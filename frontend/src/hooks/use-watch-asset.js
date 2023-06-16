function useWatchAsset () {
  const watchAsset = (wallet, toTokenInfo) => {
    return wallet.provider.send('wallet_watchAsset', {
      type: 'ERC20',
      options: {
        address: toTokenInfo.address,
        symbol: toTokenInfo.symbol,
        decimals: toTokenInfo.decimals
      }
    })
  }

  return watchAsset
}

export default useWatchAsset
