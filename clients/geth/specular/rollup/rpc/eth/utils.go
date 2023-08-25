package eth

import (
	"context"

	"github.com/avast/retry-go/v4"
	"github.com/specularl2/specular/clients/geth/specular/utils/fmt"
)

type LazyEthClient struct {
	*EthClient
	endpoint  string
	retryOpts []retry.Option
}

func NewLazilyDialedEthClient(endpoint string, retryOpts ...retry.Option) *LazyEthClient {
	return &LazyEthClient{endpoint: endpoint, retryOpts: retryOpts}
}

func (c *LazyEthClient) EnsureDialed(ctx context.Context) error {
	if c.EthClient != nil {
		return nil
	}
	client, err := DialWithRetry(ctx, c.endpoint, c.retryOpts...)
	if err != nil {
		return fmt.Errorf("failed to connect to node: %w", err)
	}
	c.EthClient = client
	return nil
}
