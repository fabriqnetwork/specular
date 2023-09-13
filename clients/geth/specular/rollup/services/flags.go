package services

import (
	"github.com/specularl2/specular/clients/geth/specular/rollup/rpc/eth/txmgr"
	"github.com/urfave/cli/v2"
)

// Returns all supported flags.
func CLIFlags() []cli.Flag {
	return mergeFlagGroups(
		l1Flags,
		l2Flags,
		sequencerCLIFlags,
		txmgr.CLIFlags(sequencerTxMgrNamespace),
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
	CmdlineFlagName = "rollup.l1.endpoint"
	// txmgr flag namespaces
	sequencerTxMgrNamespace = "rollup.sequencer.txmgr"
	validatorTxMgrNamespace = "rollup.validator.txmgr"
)

// These are all the command line flags we support.
// If you add to this list, please remember to include the
// flag in the appropriate command definition.
var (
	l1EndpointFlag = &cli.StringFlag{
		Name:  "rollup.l1.endpoint",
		Usage: "The L1 API endpoint",
	}
	l1ChainIDFlag = &cli.Uint64Flag{
		Name:  "rollup.l1.chainid",
		Usage: "The L1 chain ID",
		Value: 31337,
	}
	l1RollupGenesisBlockFlag = &cli.Uint64Flag{
		Name:  "rollup.l1.rollup-genesis-block",
		Usage: "The block number of the L1 block containing the rollup genesis",
		Value: 0,
	}
	l1SequencerInboxAddrFlag = &cli.StringFlag{
		Name:  "rollup.l1.sequencer-inbox-addr",
		Usage: "The contract address of L1 sequencer inbox",
	}
	l1RollupAddrFlag = &cli.StringFlag{
		Name:  "rollup.l1.rollup-addr",
		Usage: "The contract address of L1 rollup",
	}
	// L2 config flags
	l2EndpointFlag = &cli.StringFlag{
		Name:  "rollup.l2.endpoint",
		Usage: "The L2 API endpoint",
		Value: "ws://0.0.0.0:4012", // TODO: read this from http params?
	}
	l2ChainIDFlag = &cli.Uint64Flag{
		Name:  "rollup.l2.chainid",
		Usage: "The L2 chain ID",
	}
	l2L1FeeOverheadFlag = &cli.Int64Flag{
		Name:  "rollup.l2.l1-fee-overhead",
		Usage: "Gas cost of sequencing a Tx",
		Value: 0,
	}
	l2L1FeeMultiplierFlag = &cli.Float64Flag{
		Name:  "rollup.l2.l1-fee-multiplier",
		Usage: "Scalar value to increase the L1 Fee",
		Value: 1.5,
	}
	l2L1OracleAddressFlag = &cli.StringFlag{
		Name:  "rollup.l2.l1-oracle-address",
		Usage: "The address of the L1Oracle contract",
		Value: "0xff00000000000000000000000000000000000002",
	}
	l2L1OracleBaseFeeSlotFlag = &cli.StringFlag{
		Name:  "rollup.l2.l1-oracle-base-fee-slot",
		Usage: "The address of the L1Oracle contract",
		Value: "0x18b94da8c18f49ac05520153402a0591c3c917271b9d13711fd6fdb213ded168", // keccak256("specular.basefee")
	}
	// Sequencer config flags
	sequencerEnableSequencerFlag = &cli.BoolFlag{
		Name:  "rollup.sequencer",
		Usage: "Whether this node is a sequencer",
	}
	sequencerAddrFlag = &cli.StringFlag{
		Name:  "rollup.sequencer.addr",
		Usage: "The sequencer address",
	}
	sequencerClefEndpointFlag = &cli.StringFlag{
		Name:  "rollup.sequencer.clef-endpoint",
		Usage: "The endpoint of the Clef instance that should be used as a sequencer signer",
	}
	sequencerPassphraseFlag = &cli.StringFlag{
		Name:  "rollup.sequencer.passphrase",
		Usage: "The passphrase of the sequencer account",
	}
	sequencerSequencingIntervalFlag = &cli.UintFlag{
		Name:  "rollup.sequencer.sequencing-interval",
		Usage: "Time between batch sequencing steps (seconds)",
		Value: 8,
	}
	// Validator config flags
	validatorEnableValidatorFlag = &cli.BoolFlag{
		Name:  "rollup.validator",
		Usage: "Whether this node is a validator",
	}
	validatorAddrFlag = &cli.StringFlag{
		Name:  "rollup.validator.addr",
		Usage: "The validator address",
	}
	validatorClefEndpointFlag = &cli.StringFlag{
		Name:  "rollup.validator.clef-endpoint",
		Usage: "The endpoint of the Clef instance that should be used as a validator signer",
	}
	validatorPassphraseFlag = &cli.StringFlag{
		Name:  "rollup.validator.passphrase",
		Usage: "The passphrase of the validator account",
	}
	validatorValidationIntervalFlag = &cli.UintFlag{
		Name:  "rollup.validator.validation-interval",
		Usage: "Time between batch validation steps (seconds)",
		Value: 10,
	}
)

var (
	l1Flags = []cli.Flag{
		l1EndpointFlag,
		l1ChainIDFlag,
		l1RollupGenesisBlockFlag,
		l1SequencerInboxAddrFlag,
		l1RollupAddrFlag,
	}
	l2Flags = []cli.Flag{
		l2EndpointFlag,
		l2ChainIDFlag,
		l2L1FeeOverheadFlag,
		l2L1FeeMultiplierFlag,
		l2L1OracleAddressFlag,
		l2L1OracleBaseFeeSlotFlag,
	}
	sequencerCLIFlags = []cli.Flag{
		sequencerEnableSequencerFlag,
		sequencerAddrFlag,
		sequencerClefEndpointFlag,
		sequencerPassphraseFlag,
		sequencerSequencingIntervalFlag,
	}
	validatorCLIFlags = []cli.Flag{
		validatorEnableValidatorFlag,
		validatorAddrFlag,
		validatorClefEndpointFlag,
		validatorPassphraseFlag,
		validatorValidationIntervalFlag,
	}
)
