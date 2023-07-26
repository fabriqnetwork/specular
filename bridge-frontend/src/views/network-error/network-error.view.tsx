import useNetworkErrorStyles from './network-error.styles';
import { ReactComponent as SwitchNetwork } from '../../images/switch-network.svg';
import { ReactComponent as MetamaskIcon } from '../../images/metamask-fox.svg';
import { NETWORKS } from '../../chains';
import {
  CHIADO_NETWORK_ID
} from "../../constants";

interface NetworkButtonProps {
  isMetamask: boolean;
  switchChain: (chainId: string) => void;
}

const NetworkButton = ({ isMetamask, switchChain }: NetworkButtonProps) => {
  const classes = useNetworkErrorStyles();
  const chainId = CHIADO_NETWORK_ID.toString();
  const name = NETWORKS[chainId].chainName;
  if (!isMetamask) {
    return <p className={classes.networkName}>{name}</p>;
  }
  return (
    <div className={classes.switchNetworkButton} onClick={() => switchChain(chainId)}>
      <MetamaskIcon width="20" height="20" style={{ marginRight: 5 }} />
      <b>{name}</b>
    </div>
  );
};

interface NetworkErrorProps {
  isMetamask: boolean;
  switchChain: (chainId: string) => void;
}

function NetworkError({ isMetamask, switchChain }: NetworkErrorProps) {
  const classes = useNetworkErrorStyles();

  return (
    <div className={classes.networkError}>
      <SwitchNetwork />
      <p className={classes.title}>Switch Network</p>
      <div className={classes.descriptionContainer}>
        <p className={classes.description}>
          Please, connect to
        </p>
        <NetworkButton isMetamask={isMetamask} switchChain={switchChain} />
      </div>
    </div>
  );
}

export default NetworkError;
