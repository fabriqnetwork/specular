package rollup

import (
	"github.com/specularl2/specular/clients/geth/specular/rollup/rpc/eth/txmgr"
	"github.com/urfave/cli/v2"
)

// All supported flags.
func CLIFlags() []cli.Flag {
	var flags = []cli.Flag{}
	flags = append(flags, l1Flags...)
	flags = append(flags, l2Flags...)
	flags = append(flags, sequencerCLIFlags...)
	flags = append(flags, validatorCLIFlags...)
	flags = append(flags, driverCLIFlags...)
	flags = append(flags, txmgr.CLIFlags(sequencerTxMgrNamespace)...)
	flags = append(flags, txmgr.CLIFlags(validatorTxMgrNamespace)...)
	return flags
}

const (
	RequiredFlagName = "rollup.l1.endpoint"
	// txmgr flag namespaces
	sequencerTxMgrNamespace = "rollup.sequencer.txmgr"
	validatorTxMgrNamespace = "rollup.validator.txmgr"
)

// These are all the command line flags we support.
// If you add to this list, please remember to include the
// flag in the appropriate command definition.
var (
	// L1 config flags
	l1EndpointFlag = &cli.StringFlag{
		Name:  "rollup.l1.endpoint",
		Usage: "The API endpoint of L1 client",
	}
	l1ChainIDFlag = &cli.Uint64Flag{
		Name:  "rollup.l1.chainid",
		Usage: "The chain ID of L1 client",
		Value: 31337,
	}
	l1RollupGenesisBlockFlag = &cli.Uint64Flag{
		Name:  "rollup.l1.rollup-genesis-block",
		Usage: "The block number of the L1 block containing the rollup genesis",
		Value: 0,
	}
	sequencerInboxAddrFlag = &cli.StringFlag{
		Name:  "rollup.l1.sequencer-inbox-addr",
		Usage: "The contract address of L1 sequencer inbox",
	}
	rollupAddrFlag = &cli.StringFlag{
		Name:  "rollup.l1.rollup-addr",
		Usage: "The contract address of L1 rollup",
	}
	// L2 config flags
	l2EndpointFlag = &cli.StringFlag{
		Name:  "rollup.l2.endpoint",
		Usage: "The API endpoint of L2 client",
		Value: "ws://0.0.0.0:4012", // TODO: read this from http params?
	}
	l2ClefEndpointFlag = &cli.StringFlag{
		Name:  "rollup.l2.clef-endpoint",
		Usage: "The endpoint of the Clef instance that should be used as a signer",
	}
	// Sequencer config flags
	sequencerAddrFlag = &cli.StringFlag{
		Name:  "rollup.sequencer.addr",
		Usage: "The sequencer address to be unlocked (pass passphrash via --password)",
	}
	sequencerMinExecIntervalFlag = &cli.UintFlag{
		Name:  "rollup.sequencer.min-execution-interval",
		Usage: "Minimum time between block executions (seconds)",
		Value: 0,
	}
	sequencerMaxExecIntervalFlag = &cli.UintFlag{
		Name:  "rollup.sequencer.max-execution-interval",
		Usage: "Maximum time between block executions (seconds)",
		Value: 1,
	}
	sequencerSequencingIntervalFlag = &cli.UintFlag{
		Name:  "rollup.sequencer.sequencing-interval",
		Usage: "Time between batch sequencing steps (seconds)",
		Value: 5,
	}
	// Validator config flags
	validatorAddrFlag = &cli.StringFlag{
		Name:  "rollup.validator.addr",
		Usage: "The validator address to be unlocked (pass passphrash via --password)",
	}
	validatorIsActiveStakerFlag = &cli.BoolFlag{
		Name:  "rollup.validator.is-active-staker",
		Usage: "Whether the validator should be an active staker",
		Value: false,
	}
	validatorIsActiveCreatorFlag = &cli.BoolFlag{
		Name:  "rollup.validator.is-active-creator",
		Usage: "Whether the validator should be an active assertion creator",
		Value: false,
	}
	validatorIsActiveChallengerFlag = &cli.BoolFlag{
		Name:  "rollup.validator.is-active-challenger",
		Usage: "Whether the validator should be an active challenger (i.e. issue challenges)",
		Value: false,
	}
	validatorIsResolverFlag = &cli.BoolFlag{
		Name:  "rollup.validator.is-resolver",
		Usage: "Whether the validator should resolve (confirm/reject) assertions",
		Value: false,
	}
	// TODO: read this from the contract
	validatorStakeAmountFlag = &cli.Uint64Flag{
		Name:  "rollup.validator.stake-amount",
		Usage: "Required staking amount",
		Value: 1000000000000000000,
	}
	// Driver config flags
	driverStepIntervalFlag = &cli.Uint64Flag{
		Name:  "rollup.driver.step-interval",
		Usage: "Time between driver steps (seconds)",
		Value: 2,
	}
	driverRetryDelayFlag = &cli.Uint64Flag{
		Name:  "rollup.driver.retry-delay",
		Usage: "Time between driver retries (seconds)",
		Value: 8,
	}
	driverNumAttemptsFlag = &cli.Uint64Flag{
		Name:  "rollup.driver.num-attempts",
		Usage: "Number of driver step attempts before giving up",
		Value: 4,
	}
)

var (
	l1Flags = []cli.Flag{
		l1EndpointFlag,
		l1ChainIDFlag,
		l1RollupGenesisBlockFlag,
		sequencerInboxAddrFlag,
		rollupAddrFlag,
	}
	l2Flags = []cli.Flag{
		l2EndpointFlag,
		l2ClefEndpointFlag,
	}
	sequencerCLIFlags = []cli.Flag{
		sequencerAddrFlag,
		sequencerMinExecIntervalFlag,
		sequencerMaxExecIntervalFlag,
		sequencerSequencingIntervalFlag,
	}
	validatorCLIFlags = []cli.Flag{
		validatorAddrFlag,
		validatorIsActiveStakerFlag,
		validatorIsActiveCreatorFlag,
		validatorIsActiveChallengerFlag,
		validatorIsResolverFlag,
		validatorStakeAmountFlag,
	}
	driverCLIFlags = []cli.Flag{
		driverStepIntervalFlag,
		driverRetryDelayFlag,
		driverNumAttemptsFlag,
	}
)
