package client

import (
	"errors"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"golang.org/x/net/context"
)

const confirmationTimeout = time.Minute

func waitTransaction(ctx context.Context, client bind.DeployBackend, tx *types.Transaction) (*types.Receipt, error) {
	// TODO: pass in timeout config
	ctx, cancel := context.WithTimeout(ctx, confirmationTimeout)
	defer cancel()
	return bind.WaitMined(ctx, client, tx)
}

// Workaround for the nonce issue
func retryTransactingFunction(ctx context.Context, client bind.DeployBackend, f func() (*types.Transaction, error), retryOpts []retry.Option) (*types.Transaction, error) {
	var result *types.Transaction
	var err error
	err = retry.Do(func() error {
		result, err = f()
		return err
	}, retryOpts...)
	if err != nil {
		return nil, err
	}
	receipt, err := waitTransaction(ctx, client, result)
	if err != nil {
		return nil, err
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		return nil, errors.New("transaction failed")
	}
	return result, err
}
