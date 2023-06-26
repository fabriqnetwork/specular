import useNetworkErrorStyles from './network-error.styles';
import { ReactComponent as SwitchNetwork } from '../../images/switch-network.svg';
import { ReactComponent as MetamaskIcon } from '../../images/metamask-fox.svg';

import { NETWORKS } from '../../chains';

interface NetworkButtonProps {
  isMetamask: boolean;
  switchChainInMetaMask: (chainId: string) => void;
}

const NetworkButton = ({ isMetamask, switchChainInMetaMask }: NetworkButtonProps) => {
  const classes = useNetworkErrorStyles();
  const chainId = process.env.REACT_APP_NETWORK_ID as string;
  const name = NETWORKS[chainId].chainName;
  if (!isMetamask) {
    return <p className={classes.networkName}>{name}</p>;
  }
  return (
    <div className={classes.switchNetworkButton} onClick={() => switchChainInMetaMask(chainId)}>
      <MetamaskIcon width="20" height="20" style={{ marginRight: 5 }} />
      <b>{name}</b>
    </div>
  );
};

interface NetworkErrorProps {
  isMetamask: boolean;
  switchChainInMetaMask: (chainId: string) => void;
}

function NetworkError({ isMetamask, switchChainInMetaMask }: NetworkErrorProps) {
  const classes = useNetworkErrorStyles();

  return (
    <div className={classes.networkError}>
      <SwitchNetwork />
      <p className={classes.title}>Switch Network</p>
      <div className={classes.descriptionContainer}>
        <p className={classes.description}>
          Please, connect to
        </p>
        <NetworkButton isMetamask={isMetamask} switchChainInMetaMask={switchChainInMetaMask} />
      </div>
    </div>
  );
}

export default NetworkError;
