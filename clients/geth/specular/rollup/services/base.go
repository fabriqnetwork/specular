package services

import (
	"context"
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/log"
	"github.com/specularl2/specular/clients/geth/specular/bindings"
	"github.com/specularl2/specular/clients/geth/specular/proof"
	rollupTypes "github.com/specularl2/specular/clients/geth/specular/rollup/types"
)

type challengeCtx struct {
	opponentAssertion      *rollupTypes.Assertion
	ourAssertion           *rollupTypes.Assertion
	lastValidatedAssertion *rollupTypes.Assertion
}

type BaseService struct {
	Config *Config

	Eth                   Backend
	ProofBackend          proof.Backend
	Chain                 *core.BlockChain
	L1                    *ethclient.Client
	TransactOpts          *bind.TransactOpts
	Inbox                 *bindings.ISequencerInboxSession
	Rollup                *bindings.IRollupSession
	AssertionMap          *bindings.AssertionMapCallerSession
	challengeCh           chan *challengeCtx
	challengeResolutionCh chan *rollupTypes.Assertion

	Ctx    context.Context
	Cancel context.CancelFunc
	Wg     sync.WaitGroup
}

func NewBaseService(eth Backend, proofBackend proof.Backend, cfg *Config, auth *bind.TransactOpts) (*BaseService, error) {
	if eth == nil {
		return nil, fmt.Errorf("can not use light client with rollup")
	}
	ctx, cancel := context.WithCancel(context.Background())
	l1, err := ethclient.DialContext(ctx, cfg.L1Endpoint)
	if err != nil {
		cancel()
		return nil, err
	}
	callOpts := bind.CallOpts{
		Pending: true,
		Context: ctx,
	}
	transactOpts := bind.TransactOpts{
		From:     auth.From,
		Signer:   auth.Signer,
		GasPrice: big.NewInt(800000000),
		Context:  ctx,
	}
	inbox, err := bindings.NewISequencerInbox(common.Address(cfg.SequencerInboxAddr), l1)
	if err != nil {
		cancel()
		return nil, err
	}
	inboxSession := &bindings.ISequencerInboxSession{
		Contract:     inbox,
		CallOpts:     callOpts,
		TransactOpts: transactOpts,
	}
	rollup, err := bindings.NewIRollup(common.Address(cfg.RollupAddr), l1)
	if err != nil {
		cancel()
		return nil, err
	}
	rollupSession := &bindings.IRollupSession{
		Contract:     rollup,
		CallOpts:     callOpts,
		TransactOpts: transactOpts,
	}
	assertionMapAddr, err := rollupSession.Assertions()
	if err != nil {
		cancel()
		return nil, err
	}
	assertionMap, err := bindings.NewAssertionMapCaller(assertionMapAddr, l1)
	if err != nil {
		cancel()
		return nil, err
	}
	assertionMapSession := &bindings.AssertionMapCallerSession{
		Contract: assertionMap,
		CallOpts: callOpts,
	}
	b := &BaseService{
		Config:                cfg,
		Eth:                   eth,
		ProofBackend:          proofBackend,
		L1:                    l1,
		TransactOpts:          &transactOpts,
		Inbox:                 inboxSession,
		Rollup:                rollupSession,
		AssertionMap:          assertionMapSession,
		challengeCh:           make(chan *challengeCtx),
		challengeResolutionCh: make(chan *rollupTypes.Assertion),
		Ctx:                   ctx,
		Cancel:                cancel,
	}
	b.Chain = eth.BlockChain()
	return b, nil
}

// Start starts the rollup service
// If cleanL1 is true, the service will only start from a clean L1 history
// If stake is true, the service will try to stake on start
// Returns the genesis block
func (b *BaseService) Start(cleanL1, stake bool) *types.Block {
	// Check if we are at genesis
	// TODO: if not, sync from L1
	genesis := b.Eth.BlockChain().CurrentBlock()
	if genesis.NumberU64() != 0 {
		log.Crit("Rollup service can only start from clean history")
	}
	log.Info("Genesis root", "root", genesis.Root())

	if cleanL1 {
		inboxSize, err := b.Inbox.GetInboxSize()
		if err != nil {
			log.Crit("Failed to get initial inbox size", "err", err)
		}
		if inboxSize.Cmp(common.Big0) != 0 {
			log.Crit("Rollup service can only start from genesis")
		}
	}

	if stake {
		// Initial staking
		// TODO: sync L1 staking status
		stakeOpts := b.Rollup.TransactOpts
		stakeOpts.Value = big.NewInt(int64(b.Config.RollupStakeAmount))
		_, err := b.Rollup.Contract.Stake(&stakeOpts)
		if err != nil {
			log.Crit("Failed to stake", "err", err)
		}
	}
	return genesis
}

func (b *BaseService) ChallengeLoop() {
	defer b.Wg.Done()

	abi, err := bindings.IChallengeMetaData.GetAbi()
	if err != nil {
		log.Crit("Failed to get IChallenge ABI", "err", err)
	}

	// Watch AssertionCreated event
	createdCh := make(chan *bindings.IRollupAssertionCreated, 4096)
	createdSub, err := b.Rollup.Contract.WatchAssertionCreated(&bind.WatchOpts{Context: b.Ctx}, createdCh)
	if err != nil {
		log.Crit("Failed to watch rollup event", "err", err)
	}
	defer createdSub.Unsubscribe()

	challengedCh := make(chan *bindings.IRollupAssertionChallenged, 4096)
	challengedSub, err := b.Rollup.Contract.WatchAssertionChallenged(&bind.WatchOpts{Context: b.Ctx}, challengedCh)
	if err != nil {
		log.Crit("Failed to watch rollup event", "err", err)
	}
	defer challengedSub.Unsubscribe()

	// Watch L1 blockchain for challenge timeout
	headCh := make(chan *types.Header, 4096)
	headSub, err := b.L1.SubscribeNewHead(b.Ctx, headCh)
	if err != nil {
		log.Crit("Failed to watch l1 chain head", "err", err)
	}
	defer headSub.Unsubscribe()

	var challengeSession *bindings.IChallengeSession
	var states []*proof.ExecutionState

	var bisectedCh chan *bindings.IChallengeBisected
	var bisectedSub event.Subscription
	var challengeCompletedCh chan *bindings.IChallengeChallengeCompleted
	var challengeCompletedSub event.Subscription

	inChallenge := false
	var ctx *challengeCtx
	var opponentTimeoutBlock uint64

	for {
		if inChallenge {
			select {
			case ev := <-bisectedCh:
				// case get bisection, if is our turn
				//   if in single step, submit proof
				//   if multiple step, track current segment, update
				responder, err := challengeSession.CurrentResponder()
				if err != nil {
					// TODO: error handling
					log.Error("Can not get current responder", "error", err)
					continue
				}
				// If it's our turn
				if responder == common.Address(b.Config.Coinbase) {
					err := RespondBisection(b, abi, challengeSession, ev, states, ctx.opponentAssertion.VmHash, false)
					if err != nil {
						// TODO: error handling
						log.Error("Can not respond to bisection", "error", err)
						continue
					}
				} else {
					opponentTimeLeft, err := challengeSession.CurrentResponderTimeLeft()
					if err != nil {
						// TODO: error handling
						log.Error("Can not get current responder left time", "error", err)
						continue
					}
					log.Info("[challenge] Opponent time left", "time", opponentTimeLeft)
					opponentTimeoutBlock = ev.Raw.BlockNumber + opponentTimeLeft.Uint64()
				}
			case header := <-headCh:
				if opponentTimeoutBlock == 0 {
					continue
				}
				// TODO: can we use >= here?
				if header.Number.Uint64() > opponentTimeoutBlock {
					_, err := challengeSession.Timeout()
					if err != nil {
						log.Error("Can not timeout opponent", "error", err)
						continue
						// TODO: wait some time before retry
						// TODO: fix race condition
					}
				}
			case ev := <-challengeCompletedCh:
				// TODO: handle if we are not winner --> state corrupted
				log.Info("[challenge] Challenge completed", "winner", ev.Winner)
				bisectedSub.Unsubscribe()
				challengeCompletedSub.Unsubscribe()
				states = []*proof.ExecutionState{}
				inChallenge = false
				b.challengeResolutionCh <- ctx.ourAssertion
			case <-b.Ctx.Done():
				bisectedSub.Unsubscribe()
				challengeCompletedSub.Unsubscribe()
				return
			}
		} else {
			select {
			case ctx = <-b.challengeCh:
				_, err = b.Rollup.CreateAssertion(
					ctx.ourAssertion.VmHash,
					ctx.ourAssertion.InboxSize,
					ctx.ourAssertion.CumulativeGasUsed,
					ctx.lastValidatedAssertion.VmHash,
					ctx.lastValidatedAssertion.CumulativeGasUsed,
				)
				if err != nil {
					log.Crit("UNHANDELED: Can't create assertion for challenge, validator state corrupted", "err", err)
				}
			case ev := <-createdCh:
				if common.Address(ev.AsserterAddr) == b.Config.Coinbase {
					if ev.VmHash == ctx.ourAssertion.VmHash {
						_, err := b.Rollup.ChallengeAssertion(
							[2]common.Address{
								common.Address(b.Config.SequencerAddr),
								common.Address(b.Config.Coinbase),
							},
							[2]*big.Int{
								ctx.opponentAssertion.ID,
								ev.AssertionID,
							},
						)
						if err != nil {
							log.Crit("UNHANDELED: Can't start challenge, validator state corrupted", "err", err)
						}
					}
				}
			case ev := <-challengedCh:
				if ctx == nil {
					continue
				}
				log.Info("validator saw challenge", "assertion id", ev.AssertionID, "expected id", ctx.opponentAssertion.ID, "block", ev.Raw.BlockNumber)
				if ev.AssertionID.Cmp(ctx.opponentAssertion.ID) == 0 {
					// start := ev.Raw.BlockNumber - 2
					challenge, err := bindings.NewIChallenge(ev.ChallengeAddr, b.L1)
					if err != nil {
						log.Crit("Failed to access ongoing challenge", "address", ev.ChallengeAddr, "err", err)
					}
					challengeSession = &bindings.IChallengeSession{
						Contract:     challenge,
						CallOpts:     bind.CallOpts{Pending: true, Context: b.Ctx},
						TransactOpts: *b.TransactOpts,
					}
					bisectedCh = make(chan *bindings.IChallengeBisected, 4096)
					bisectedSub, err = challenge.WatchBisected(&bind.WatchOpts{Context: b.Ctx}, bisectedCh)
					if err != nil {
						log.Crit("Failed to watch challenge event", "err", err)
					}
					challengeCompletedCh = make(chan *bindings.IChallengeChallengeCompleted, 4096)
					challengeCompletedSub, err = challenge.WatchChallengeCompleted(&bind.WatchOpts{Context: b.Ctx}, challengeCompletedCh)
					if err != nil {
						log.Crit("Failed to watch challenge event", "err", err)
					}
					states, err = proof.GenerateStates(
						b.ProofBackend,
						b.Ctx,
						ctx.opponentAssertion.PrevCumulativeGasUsed,
						ctx.opponentAssertion.StartBlock,
						ctx.opponentAssertion.EndBlock+1,
						nil,
					)
					if err != nil {
						log.Crit("Failed to generate states", "err", err)
					}
					inChallenge = true
				}
			case <-headCh:
				continue // consume channel values
			case <-b.Ctx.Done():
				return
			}
		}
	}
}
