package txmgr

import (
	"errors"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/urfave/cli/v2"
)

type Config struct {
	// ResubmissionTimeout is the interval at which, if no previously
	// published transaction has been mined, the new tx with a bumped gas
	// price will be published. Only one publication at MaxGasPrice will be
	// attempted.
	ResubmissionTimeout time.Duration

	// The multiplier applied to fee suggestions to put a hard limit on fee increases.
	FeeLimitMultiplier uint64

	// ChainID is the chain ID of the L1 chain.
	ChainID *big.Int

	// TxSendTimeout is how long to wait for sending a transaction.
	// By default it is unbounded. If set, this is recommended to be at least 20 minutes.
	TxSendTimeout time.Duration

	// TxNotInMempoolTimeout is how long to wait before aborting a transaction send if the transaction does not
	// make it to the mempool. If the tx is in the mempool, TxSendTimeout is used instead.
	TxNotInMempoolTimeout time.Duration

	// NetworkTimeout is the allowed duration for a single network request.
	// This is intended to be used for network requests that can be replayed.
	NetworkTimeout time.Duration

	// RequireQueryInterval is the interval at which the tx manager will
	// query the backend to check for confirmations after a tx at a
	// specific gas price has been published.
	ReceiptQueryInterval time.Duration

	// NumConfirmations specifies how many blocks are need to consider a
	// transaction confirmed.
	NumConfirmations uint64

	// SafeAbortNonceTooLowCount specifies how many ErrNonceTooLow observations
	// are required to give up on a tx at a particular nonce without receiving
	// confirmation.
	SafeAbortNonceTooLowCount uint64

	From common.Address
}

func (m Config) Validate() error {
	if m.NumConfirmations == 0 {
		return errors.New("NumConfirmations must not be 0")
	}
	if m.NetworkTimeout == 0 {
		return errors.New("must provide NetworkTimeout")
	}
	if m.FeeLimitMultiplier == 0 {
		return errors.New("must provide FeeLimitMultiplier")
	}
	if m.ResubmissionTimeout == 0 {
		return errors.New("must provide ResubmissionTimeout")
	}
	if m.ReceiptQueryInterval == 0 {
		return errors.New("must provide ReceiptQueryInterval")
	}
	if m.TxNotInMempoolTimeout == 0 {
		return errors.New("must provide TxNotInMempoolTimeout")
	}
	if m.SafeAbortNonceTooLowCount == 0 {
		return errors.New("SafeAbortNonceTooLowCount must not be 0")
	}
	if m.ChainID == nil {
		return errors.New("must provide the ChainID")
	}
	return nil
}

const (
	ResubmissionTimeoutFlagName       = "resubmission-timeout"
	FeeLimitMultiplierFlagName        = "fee-limit-multiplier"
	TxSendTimeoutFlagName             = "send-timeout"
	TxNotInMempoolTimeoutFlagName     = "not-in-mempool-timeout"
	NetworkTimeoutFlagName            = "network-timeout"
	ReceiptQueryIntervalFlagName      = "receipt-query-interval"
	NumConfirmationsFlagName          = "num-confirmations"
	SafeAbortNonceTooLowCountFlagName = "safe-abort-nonce-too-low-count"
)

type DefaultFlagValues struct {
	NumConfirmations          uint64
	SafeAbortNonceTooLowCount uint64
	FeeLimitMultiplier        uint64
	ResubmissionTimeout       time.Duration
	NetworkTimeout            time.Duration
	TxSendTimeout             time.Duration
	TxNotInMempoolTimeout     time.Duration
	ReceiptQueryInterval      time.Duration
}

var (
	DefaultDisseminatorFlagValues = DefaultFlagValues{
		NumConfirmations:          uint64(6),
		SafeAbortNonceTooLowCount: uint64(3),
		FeeLimitMultiplier:        uint64(5),
		ResubmissionTimeout:       48 * time.Second,
		NetworkTimeout:            10 * time.Second,
		TxSendTimeout:             0 * time.Second,
		TxNotInMempoolTimeout:     2 * time.Minute,
		ReceiptQueryInterval:      12 * time.Second,
	}
	// TODO: tweak these values.
	DefaultValidatorFlagValues = DefaultFlagValues{
		NumConfirmations:          uint64(10),
		SafeAbortNonceTooLowCount: uint64(3),
		FeeLimitMultiplier:        uint64(5),
		ResubmissionTimeout:       48 * time.Second,
		NetworkTimeout:            10 * time.Second,
		TxSendTimeout:             0 * time.Second,
		TxNotInMempoolTimeout:     2 * time.Minute,
		ReceiptQueryInterval:      12 * time.Second,
	}
)

func CLIFlags(namespace string, defaults DefaultFlagValues) []cli.Flag {
	return []cli.Flag{
		&cli.Uint64Flag{
			Name:  namespace + "." + NumConfirmationsFlagName,
			Usage: "Number of confirmations which we will wait after sending a transaction",
			Value: defaults.NumConfirmations,
		},
		&cli.Uint64Flag{
			Name:  namespace + "." + SafeAbortNonceTooLowCountFlagName,
			Usage: "Number of ErrNonceTooLow observations required to give up on a tx at a particular nonce without receiving confirmation",
			Value: defaults.SafeAbortNonceTooLowCount,
		},
		&cli.Uint64Flag{
			Name:  namespace + "." + FeeLimitMultiplierFlagName,
			Usage: "The multiplier applied to fee suggestions to put a hard limit on fee increases",
			Value: defaults.FeeLimitMultiplier,
		},
		&cli.DurationFlag{
			Name:  namespace + "." + ResubmissionTimeoutFlagName,
			Usage: "Duration we will wait before resubmitting a transaction to L1",
			Value: defaults.ResubmissionTimeout,
		},
		&cli.DurationFlag{
			Name:  namespace + "." + NetworkTimeoutFlagName,
			Usage: "Timeout for all network operations",
			Value: defaults.NetworkTimeout,
		},
		&cli.DurationFlag{
			Name:  namespace + "." + TxSendTimeoutFlagName,
			Usage: "Timeout for sending transactions. If 0 it is disabled.",
			Value: defaults.TxSendTimeout,
		},
		&cli.DurationFlag{
			Name:  namespace + "." + TxNotInMempoolTimeoutFlagName,
			Usage: "Timeout for aborting a tx send if the tx does not make it to the mempool.",
			Value: defaults.TxNotInMempoolTimeout,
		},
		&cli.DurationFlag{
			Name:  namespace + "." + ReceiptQueryIntervalFlagName,
			Usage: "Frequency to poll for receipts",
			Value: defaults.ReceiptQueryInterval,
		},
	}
}

func NewConfigFromCLI(
	cliCtx *cli.Context,
	namespace string,
	chainID *big.Int,
	from common.Address,
) Config {
	return Config{
		ResubmissionTimeout:       cliCtx.Duration(namespace + "." + ResubmissionTimeoutFlagName),
		FeeLimitMultiplier:        cliCtx.Uint64(namespace + "." + FeeLimitMultiplierFlagName),
		ChainID:                   chainID,
		TxSendTimeout:             cliCtx.Duration(namespace + "." + TxSendTimeoutFlagName),
		TxNotInMempoolTimeout:     cliCtx.Duration(namespace + "." + TxNotInMempoolTimeoutFlagName),
		NetworkTimeout:            cliCtx.Duration(namespace + "." + NetworkTimeoutFlagName),
		ReceiptQueryInterval:      cliCtx.Duration(namespace + "." + ReceiptQueryIntervalFlagName),
		NumConfirmations:          cliCtx.Uint64(namespace + "." + NumConfirmationsFlagName),
		SafeAbortNonceTooLowCount: cliCtx.Uint64(namespace + "." + SafeAbortNonceTooLowCountFlagName),
		From:                      from,
	}
}
