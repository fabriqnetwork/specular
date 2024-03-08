package services

import (
	"github.com/ethereum/go-ethereum/log"
	"github.com/urfave/cli/v2"

	"github.com/specularL2/specular/services/sidecar/rollup/rpc/eth/txmgr"
)

// Returns all supported flags.
func CLIFlags() []cli.Flag {
	return mergeFlagGroups(
		generalFlags,
		protocolFlags,
		disseminatorCLIFlags,
		txmgr.CLIFlags(disseminatorTxMgrNamespace, txmgr.DefaultDisseminatorFlagValues),
		validatorCLIFlags,
		txmgr.CLIFlags(validatorTxMgrNamespace, txmgr.DefaultValidatorFlagValues),
	)
}

// Merges flag groups into a single slice.
func mergeFlagGroups(groups ...[]cli.Flag) []cli.Flag {
	var flags []cli.Flag
	for _, group := range groups {
		flags = append(flags, group...)
	}
	return flags
}

const (
	// txmgr flag namespaces
	disseminatorTxMgrNamespace = "disseminator.txmgr"
	validatorTxMgrNamespace    = "validator.txmgr"
)

// These are all the command line flags we support.
// If you add to this list, please remember to include the
// flag in the appropriate command definition.
var (
	VerbosityFlag = &cli.IntFlag{
		Name:  "verbosity",
		Usage: "Set the log verbosity level. 0 = silent, 1 = error, 2 = warn, 3 = info, 4 = debug, 5 = trace",
		Value: int(log.LvlInfo),
	}
	// L1 config flags
	l1EndpointFlag = &cli.StringFlag{
		Name:     "l1.endpoint",
		Usage:    "The L1 API endpoint",
		Required: true,
	}
	l1SubmissionEndpointFlag = &cli.StringFlag{
		Name:     "l1.submission-endpoint",
		Usage:    "The L1 API submission endpoint",
		Required: false,
	}
	// L2 config flags
	l2EndpointFlag = &cli.StringFlag{
		Name:     "l2.endpoint",
		Usage:    "The L2 API endpoint",
		Required: true,
	}
	// Chain config protocol flags.
	protocolRollupCfgPathFlag = &cli.StringFlag{
		Name:     "protocol.rollup-cfg-path",
		Usage:    "The path to the L2 rollup config file",
		Required: true,
	}
	// Disseminator config flags
	disseminatorEnableFlag = &cli.BoolFlag{
		Name:  "disseminator",
		Usage: "Whether this node is a disseminator",
	}
	disseminatorPrivateKeyFlag = &cli.StringFlag{
		Name:  "disseminator.private-key",
		Usage: "The private key for rollup_cfg['system_config']['batcherAddr']]",
	}
	disseminatorClefEndpointFlag = &cli.StringFlag{
		Name:  "disseminator.clef-endpoint",
		Usage: "The endpoint of the Clef instance that should be used as a disseminator signer",
	}
	disseminatorIntervalFlag = &cli.UintFlag{
		Name:  "disseminator.interval",
		Usage: "Time between batch dissemination steps (seconds)",
		Value: 8,
	}
	disseminatorSubSafetyMarginFlag = &cli.Uint64Flag{
		Name:  "disseminator.sub-safety-margin",
		Usage: "The safety margin for batch tx submission (in # of L1 blocks)",
	}
	disseminatorTargetBatchSizeFlag = &cli.Uint64Flag{
		Name:  "disseminator.target-batch-size",
		Usage: "The target size of a batch tx submitted to L1 (bytes)",
	}
	disseminatorMaxBatchSizeFlag = &cli.Uint64Flag{
		Name:  "disseminator.max-batch-size",
		Usage: "The maximun size of a batch tx submitted to L1 (bytes)",
	}
	disseminatorMaxSafeLagFlag = &cli.Uint64Flag{
		Name:  "disseminator.max-safe-lag",
		Usage: "The maximum, in l2 blocks, that is safe for the disseminator to lag the sequencer",
	}
	disseminatorMaxSafeLagDeltaFlag = &cli.Uint64Flag{
		Name:  "disseminator.max-safe-lag-delta",
		Usage: "The delta gap, in l2 blocks, to use for forcing a batch when lagging",
	}
	// Validator config flags
	validatorEnableFlag = &cli.BoolFlag{
		Name:  "validator",
		Usage: "Whether this node is a validator",
	}
	validatorPrivateKeyFlag = &cli.StringFlag{
		Name:  "validator.private-key",
		Usage: "The private key for the validator",
	}
	validatorClefEndpointFlag = &cli.StringFlag{
		Name:  "validator.clef-endpoint",
		Usage: "The endpoint of the Clef instance that should be used as a validator signer",
	}
	validatorValidationIntervalFlag = &cli.UintFlag{
		Name:  "validator.validation-interval",
		Usage: "Time between validation steps (seconds)",
		Value: 10,
	}
)

var (
	generalFlags         = []cli.Flag{VerbosityFlag, l1EndpointFlag, l1SubmissionEndpointFlag, l2EndpointFlag}
	protocolFlags        = []cli.Flag{protocolRollupCfgPathFlag}
	disseminatorCLIFlags = []cli.Flag{
		disseminatorEnableFlag,
		disseminatorPrivateKeyFlag,
		disseminatorClefEndpointFlag,
		disseminatorIntervalFlag,
		disseminatorSubSafetyMarginFlag,
		disseminatorTargetBatchSizeFlag,
		disseminatorMaxBatchSizeFlag,
		disseminatorMaxSafeLagFlag,
		disseminatorMaxSafeLagDeltaFlag,
	}
	validatorCLIFlags = []cli.Flag{
		validatorEnableFlag,
		validatorPrivateKeyFlag,
		validatorClefEndpointFlag,
		validatorValidationIntervalFlag,
	}
)
