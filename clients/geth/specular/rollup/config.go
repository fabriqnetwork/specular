package rollup

import (
	"time"

	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/specularl2/specular/clients/geth/specular/rollup/services"
	"github.com/urfave/cli/v2"
)

// These are all the command line flags we support.
// If you add to this list, please remember to include the
// flag in the appropriate command definition.
var (
	// L1 config flags
	L1EndpointFlag = &cli.StringFlag{
		Name:  "rollup.l1-endpoint",
		Usage: "The api endpoint of L1 client",
		Value: "",
	}
	L1ChainIDFlag = &cli.Uint64Flag{
		Name:  "rollup.l1-chainid",
		Usage: "The chain ID of L1 client",
		Value: 31337,
	}
	L1RollupGenesisBlockFlag = &cli.Uint64Flag{
		Name:  "rollup.l1-rollup-genesis-block",
		Usage: "The block number of L1 rollup genesis block to sync from",
		Value: 0,
	}
	SequencerInboxAddrFlag = &cli.StringFlag{
		Name:  "rollup.l1-sequencer-inbox-addr",
		Usage: "The contract address of L1 sequencer inbox",
		Value: "",
	}
	RollupAddrFlag = &cli.StringFlag{
		Name:  "rollup.l1-rollup-addr",
		Usage: "The contract address of L1 rollup",
		Value: "",
	}
	// L2 config flags
	L2ClefEndpointFlag = &cli.StringFlag{
		Name:  "rollup.l2-clef-endpoint",
		Usage: "The Endpoint of the Clef instance that should be used as a signer",
		Value: "",
	}
	// Sequencer config flags
	SequencerAddrFlag = &cli.StringFlag{
		Name:  "rollup.sequencer-addr",
		Usage: "The sequencer address to be unlocked (pass passphrash via --password)",
		Value: "",
	}
	SequencerMinExecutionIntervalFlag = &cli.Uint64Flag{
		Name:  "rollup.sequencer-min-execution-interval",
		Usage: "Minimum time between block executions (seconds)",
		Value: 0,
	}
	SequencerMaxExecutionIntervalFlag = &cli.Uint64Flag{
		Name:  "rollup.sequencer-max-execution-interval",
		Usage: "Maximum time between block executions (seconds)",
		Value: 1,
	}
	SequencerSequencingIntervalFlag = &cli.Uint64Flag{
		Name:  "rollup.sequencer-sequencing-interval",
		Usage: "Time between batch sequencing attempts (seconds)",
		Value: 5,
	}
	// Validator config flags
	ValidatorAddrFlag = &cli.StringFlag{
		Name:  "rollup.validator-addr",
		Usage: "The validator address to be unlocked (pass passphrash via --password)",
		Value: "",
	}
	ValidatorIsActiveStakerFlag = &cli.BoolFlag{
		Name:  "rollup.validator-is-active-staker",
		Usage: "Whether the validator should be an active staker",
		Value: false,
	}
	ValidatorIsActiveCreatorFlag = &cli.BoolFlag{
		Name:  "rollup.validator-is-active-creator",
		Usage: "Whether the validator should be an active assertion creator",
		Value: false,
	}
	ValidatorIsActiveChallengerFlag = &cli.BoolFlag{
		Name:  "rollup.validator-is-active-challenger",
		Usage: "Whether the validator should be an active challenger (i.e. issue challenges)",
		Value: false,
	}
	ValidatorIsResolverFlag = &cli.BoolFlag{
		Name:  "rollup.validator-is-resolver",
		Usage: "Whether the validator should resolve (confirm/reject) assertions",
		Value: false,
	}
	// TODO: read this from the contract
	RollupStakeAmountFlag = &cli.Uint64Flag{
		Name:  "rollup.stake-amount",
		Usage: "Required staking amount",
		Value: 1000000000000000000,
	}
)

// All supported flags.
var Flags = []cli.Flag{
	// L1 config flags
	L1EndpointFlag,
	L1ChainIDFlag,
	L1RollupGenesisBlockFlag,
	SequencerInboxAddrFlag,
	RollupAddrFlag,
	// L2 config flags
	L2ClefEndpointFlag,
	// Sequencer config flags
	SequencerAddrFlag,
	SequencerMinExecutionIntervalFlag,
	SequencerMaxExecutionIntervalFlag,
	SequencerSequencingIntervalFlag,
	// Validator config flags
	ValidatorAddrFlag,
	ValidatorIsActiveStakerFlag,
	ValidatorIsActiveCreatorFlag,
	ValidatorIsActiveChallengerFlag,
	ValidatorIsResolverFlag,
	RollupStakeAmountFlag,
}

func ParseSystemConfig(ctx *cli.Context) *services.SystemConfig {
	utils.CheckExclusive(ctx, L1EndpointFlag, utils.MiningEnabledFlag)
	utils.CheckExclusive(ctx, L1EndpointFlag, utils.DeveloperFlag)
	var (
		sequencerAddr       common.Address
		validatorAddr       common.Address
		sequencerPassphrase string
		validatorPassphrase string
	)
	pwList := utils.MakePasswordList(ctx)
	clefEndpoint := ctx.String(L2ClefEndpointFlag.Name)
	if clefEndpoint == "" {
		if len(pwList) == 0 {
			utils.Fatalf("Failed to register rollup services: no clef endpoint or password provided")
		}
		if ctx.String(SequencerAddrFlag.Name) != "" {
			sequencerAddr = common.HexToAddress(ctx.String(SequencerAddrFlag.Name))
			sequencerPassphrase = pwList[0]
			pwList = pwList[1:]
		}
		if ctx.String(ValidatorAddrFlag.Name) != "" {
			validatorAddr = common.HexToAddress(ctx.String(ValidatorAddrFlag.Name))
			validatorPassphrase = pwList[0]
			pwList = pwList[1:]
		}
	}
	return services.NewSystemConfig(
		// L1 params
		ctx.String(L1EndpointFlag.Name),
		ctx.Uint64(L1ChainIDFlag.Name),
		ctx.Uint64(L1RollupGenesisBlockFlag.Name),
		common.HexToAddress(ctx.String(SequencerInboxAddrFlag.Name)),
		common.HexToAddress(ctx.String(RollupAddrFlag.Name)),
		// L2 params
		"ws://0.0.0.0:4012", // TODO: read this from http params? or from a separate flag
		clefEndpoint,
		// Sequencer params
		sequencerAddr,
		sequencerPassphrase,
		time.Duration(ctx.Uint64(SequencerMinExecutionIntervalFlag.Name))*time.Second,
		time.Duration(ctx.Uint64(SequencerMaxExecutionIntervalFlag.Name))*time.Second,
		time.Duration(ctx.Uint64(SequencerSequencingIntervalFlag.Name))*time.Second,
		// Validator params
		validatorAddr,
		validatorPassphrase,
		ctx.Bool(ValidatorIsActiveStakerFlag.Name),
		ctx.Bool(ValidatorIsActiveCreatorFlag.Name),
		ctx.Bool(ValidatorIsActiveChallengerFlag.Name),
		ctx.Bool(ValidatorIsResolverFlag.Name),
		ctx.Uint64(RollupStakeAmountFlag.Name),
		// Driver params
		time.Duration(6)*time.Second,
		time.Duration(1)*time.Second,
		3,
	)
}
