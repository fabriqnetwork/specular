import { useState, useCallback, useEffect } from 'react';
import { SafeAppWeb3Modal } from '@gnosis.pm/safe-apps-web3modal';
import WalletConnectProvider from '@walletconnect/web3-provider';
import WalletLink from 'walletlink';
import { providers, utils } from 'ethers';

import coinbaseLogo from '../images/coinbase.png';
import walletConnectLogo from '../images/walletconnect.svg';

import { NETWORKS } from '../chains';

const web3Modal = new SafeAppWeb3Modal({
  cacheProvider: true,
  providerOptions: {
    'custom-walletlink': {
      display: {
        logo: coinbaseLogo,
        name: 'Coinbase Wallet',
        description: 'Scan with Coinbase Wallet to connect',
      },
      package: WalletLink,
      connector: async (ProviderPackage) => {
        const provider = new ProviderPackage({ appName: 'xDAI to ETH' }).makeWeb3Provider({}, 0);
        await provider.enable();
        return provider;
      },
    },
    'custom-walletconnect': {
      display: {
        logo: walletConnectLogo,
        name: 'WalletConnect',
        description: 'Scan with WalletConnect to connect',
      },
      package: WalletConnectProvider,
      options: {
        chainId: Number(process.env.REACT_APP_NETWORK_ID),
        rpc: {
          [process.env.REACT_APP_NETWORK_ID]: process.env.REACT_APP_RPC_URL
        },
      },
      connector: async (ProviderPackage, options) => {
        const provider = new ProviderPackage(options);
        await provider.enable();
        return provider;
      }
    }
  },
});

async function switchChainInMetaMask(chainId) {
  const { name, symbol, chainName, rpcUrl, blockExplorerUrl } = NETWORKS[chainId];
  try {
    await window.ethereum.request({
      method: 'wallet_switchEthereumChain',
      params: [
        {
          chainId: utils.hexValue(Number(chainId)),
        },
      ],
    });
    return true;
  } catch (switchError) {
    // This error code indicates that the chain has not been added to MetaMask.
    if (switchError.code === 4902) {
      try {
        if (chainId !== '100') throw Error();
        await window.ethereum.request({
          method: 'wallet_addEthereumChain',
          params: [
            {
              chainId: utils.hexValue(Number(chainId)),
              chainName,
              nativeCurrency: {
                name,
                symbol,
                decimals: 18,
              },
              rpcUrls: [rpcUrl],
              blockExplorerUrls: [blockExplorerUrl],
            },
          ],
        });
        return true;
      } catch (addError) {
        console.log(addError);
      }
    } else {
      console.log(switchError);
    }
    return false;
  }
};

function useWallet() {
  const [wallet, setWallet] = useState();
  const [isMetamask, setIsMetamask] = useState(false);

  const closeConnection = useCallback(async () => {
    const provider = wallet?.provider;
    if (provider && provider.currentProvider && provider.currentProvider.close) {
      await provider.currentProvider.close();
    }
    await web3Modal.clearCachedProvider();
    await web3Modal.off();
    window.localStorage.clear();
    window.location.reload();
  }, [wallet]);

  const loadWallet = useCallback(async () => {
    const provider = await web3Modal.requestProvider();
    async function connect() {
      const library = new providers.Web3Provider(provider);
      const network = await library.getNetwork();
      const address = await library.getSigner().getAddress();
      const chainId = String(network.chainId);
      setIsMetamask(library?.connection?.url === 'metamask');
      setWallet({ provider: library, address, chainId });
    }
    if (provider.on) {
      provider.on('close', closeConnection);
      provider.on('disconnect', closeConnection);
      provider.on('accountsChanged', accounts => accounts.length ? connect() : window.location.reload());
      // provider.on('networkChanged', connect);
      provider.on('chainChanged', closeConnection);
    }
    provider.autoRefreshOnNetworkChange = false;
    await connect();
  }, [closeConnection]);

  const disconnectWallet = useCallback(async () => {
    await web3Modal.clearCachedProvider();
    await web3Modal.off();
    window.localStorage.clear();
    window.location.reload();
  }, []);

  useEffect(() => {
    async function connect() {
      if (await web3Modal.isSafeApp()) {
        try {
          await loadWallet();
        } catch (error) {
          console.log(error);
          await disconnectWallet();
        }
      }
    }
    connect();
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  return ({
    wallet,
    isMetamask,
    loadWallet,
    disconnectWallet,
    switchChainInMetaMask,
  });
};

export default useWallet;