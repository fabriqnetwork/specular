package crypto

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/specularl2/specular/clients/geth/specular/utils/fmt"
)

// SignerFn is a generic transaction signing function.
// It takes the address that should be used to sign the transaction with.
type SignerFn func(addr common.Address, tx *types.Transaction) (*types.Transaction, error)

func NewSignerFn(privateKey string, chainID *big.Int) (SignerFn, error) {
	privKey, err := crypto.HexToECDSA(strings.TrimPrefix(privateKey, "0x"))
	if err != nil {
		return nil, fmt.Errorf("Failed to parse private key: %w", err)
	}
	from := crypto.PubkeyToAddress(privKey.PublicKey)
	signer := types.LatestSignerForChainID(chainID)
	return func(addr common.Address, tx *types.Transaction) (*types.Transaction, error) {
		if addr != from {
			return nil, bind.ErrNotAuthorized
		}
		signature, err := crypto.Sign(signer.Hash(tx).Bytes(), privKey)
		if err != nil {
			return nil, err
		}
		return tx.WithSignature(signer, signature)
	}, nil
}
