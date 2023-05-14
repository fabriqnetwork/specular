package bridge

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/specularl2/specular/clients/geth/specular/rollup/rpc/txmgr"
)

// Wrapper around EthTxManager for bridge contract methods.
type TxManager struct {
	EthTxManager
	cfg bridgeConfig
}

type EthTxManager interface {
	Send(ctx context.Context, candidate txmgr.TxCandidate) (*types.Receipt, error)
}

type bridgeConfig interface {
	SequencerInboxAddr() common.Address
	RollupAddr() common.Address
}

func NewTxManager(txMgr EthTxManager, cfg bridgeConfig) (*TxManager, error) {
	err := ensureUtilInit()
	return &TxManager{EthTxManager: txMgr, cfg: cfg}, err
}

// ISequencerInbox

func (m *TxManager) AppendTxBatch(
	ctx context.Context,
	contexts,
	txLengths []*big.Int,
	txs []byte,
) (*types.Receipt, error) {
	data, err := packAppendTxBatchInput(contexts, txLengths, txs)
	if err != nil {
		return nil, err
	}
	addr := m.cfg.SequencerInboxAddr()
	return m.Send(ctx, txmgr.TxCandidate{TxData: data, To: &addr})
}

// IRollup

func (m *TxManager) Stake(ctx context.Context, stakeAmount *big.Int) (*types.Receipt, error) {
	data, err := packStakeInput(stakeAmount)
	if err != nil {
		return nil, err
	}
	return m.sendRollupTx(ctx, data)
}

func (m *TxManager) AdvanceStake(ctx context.Context, assertionID *big.Int) (*types.Receipt, error) {
	data, err := packAdvanceStakeInput(assertionID)
	if err != nil {
		return nil, err
	}
	return m.sendRollupTx(ctx, data)
}

func (m *TxManager) CreateAssertion(ctx context.Context, vmHash common.Hash, inboxSize *big.Int) (*types.Receipt, error) {
	data, err := packCreateAssertionInput(vmHash, inboxSize)
	if err != nil {
		return nil, err
	}
	return m.sendRollupTx(ctx, data)
}

func (m *TxManager) ConfirmFirstUnresolvedAssertion(ctx context.Context) (*types.Receipt, error) {
	data, err := packConfirmFirstUnresolvedAssertionInput()
	if err != nil {
		return nil, err
	}
	return m.sendRollupTx(ctx, data)
}

func (m *TxManager) RejectFirstUnresolvedAssertion(ctx context.Context, stakerAddress common.Address) (*types.Receipt, error) {
	data, err := packRejectFirstUnresolvedAssertionInput(stakerAddress)
	if err != nil {
		return nil, err
	}
	return m.sendRollupTx(ctx, data)
}

func (m *TxManager) sendRollupTx(ctx context.Context, data []byte) (*types.Receipt, error) {
	addr := m.cfg.RollupAddr()
	return m.Send(ctx, txmgr.TxCandidate{TxData: data, To: &addr})
}
