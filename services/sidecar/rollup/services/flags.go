package services

import (
	"github.com/specularL2/specular/services/sidecar/rollup/rpc/eth/txmgr"
	"github.com/urfave/cli/v2"
)

// Returns all supported flags.
func CLIFlags() []cli.Flag {
	return mergeFlagGroups(
		protocolFlags,
		l1Flags,
		l2Flags,
		disseminatorCLIFlags,
		txmgr.CLIFlags(disseminatorTxMgrNamespace),
		validatorCLIFlags,
		txmgr.CLIFlags(validatorTxMgrNamespace),
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
	// L1 config flags
	l1EndpointFlag = &cli.StringFlag{
		Name:  "l1.endpoint",
		Usage: "The L1 API endpoint",
	}
	// L2 config flags
	l2EndpointFlag = &cli.StringFlag{
		Name:  "l2.endpoint",
		Usage: "The L2 API endpoint",
	}
	// Chain config protocol flags.
	protocolRollupCfgPathFlag = &cli.StringFlag{
		Name:  "protocol.rollup-cfg-path",
		Usage: "The path to the L2 rollup config file",
	}
	protocolRollupAddrFlag = &cli.StringFlag{
		Name:  "protocol.rollup-addr",
		Usage: "The contract address of L1 rollup",
	}
	protocolL1OracleAddrFlag = &cli.StringFlag{
		Name:  "protocol.l1-oracle-addr",
		Usage: "The address of the L1Oracle contract",
		Value: "0xff00000000000000000000000000000000000002",
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
	// Validator config flags
	validatorEnableFlag = &cli.BoolFlag{
		Name:  "validator",
		Usage: "Whether this node is a validator",
	}
	validatorAddrFlag = &cli.StringFlag{
		Name:  "validator.addr",
		Usage: "The validator address",
	}
	validatorPrivateKeyFlag = &cli.StringFlag{
		Name:  "validator.private-key",
		Usage: "The private key for validator.addr",
	}
	validatorClefEndpointFlag = &cli.StringFlag{
		Name:  "validator.clef-endpoint",
		Usage: "The endpoint of the Clef instance that should be used as a validator signer",
	}
	validatorValidationIntervalFlag = &cli.UintFlag{
		Name:  "validator.validation-interval",
		Usage: "Time between batch validation steps (seconds)",
		Value: 10,
	}
)

var (
	l1Flags       = []cli.Flag{l1EndpointFlag}
	l2Flags       = []cli.Flag{l2EndpointFlag}
	protocolFlags = []cli.Flag{
		protocolRollupCfgPathFlag,
		protocolRollupAddrFlag,
		protocolL1OracleAddrFlag,
	}
	disseminatorCLIFlags = []cli.Flag{
		disseminatorEnableFlag,
		disseminatorPrivateKeyFlag,
		disseminatorClefEndpointFlag,
		disseminatorIntervalFlag,
	}
	validatorCLIFlags = []cli.Flag{
		validatorEnableFlag,
		validatorAddrFlag,
		validatorPrivateKeyFlag,
		validatorClefEndpointFlag,
		validatorValidationIntervalFlag,
	}
)
