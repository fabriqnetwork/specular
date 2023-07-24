function useWatchAsset(): (wallet: any, toTokenInfo: TokenInfo) => Promise<any> {
  const watchAsset = (wallet: any, toTokenInfo: TokenInfo): Promise<any> => {
    return wallet.provider.send('wallet_watchAsset', {
      type: 'ERC20',
      options: {
        address: toTokenInfo.address,
        symbol: toTokenInfo.symbol,
        decimals: toTokenInfo.decimals
      }
    });
  };

  return watchAsset;
}

export default useWatchAsset;

interface TokenInfo {
  address: string;
  symbol: string;
  decimals: number;
}
