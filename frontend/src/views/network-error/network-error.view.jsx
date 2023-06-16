import useNetworkErrorStyles from './network-error.styles'
import { ReactComponent as SwitchNetwork } from '../../images/switch-network.svg'
import { ReactComponent as MetamaskIcon } from '../../images/metamask-fox.svg'

import { NETWORKS } from '../../chains'

const NetworkButton = ({ isMetamask, switchChainInMetaMask }) => {
  const classes = useNetworkErrorStyles()
  const chainId = process.env.REACT_APP_NETWORK_ID
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
}

function NetworkError ({ isMetamask, switchChainInMetaMask }) {
  const classes = useNetworkErrorStyles()

  return (
    <div className={classes.networkError}>
      <SwitchNetwork />
      <p className={classes.title}>Switch Network</p>
      <div className={classes.descriptionContainer}>
        <p className={classes.description}>
          Please, connect to
        </p>
        <NetworkButton {...{ isMetamask, switchChainInMetaMask }} />
      </div>
    </div>
  )
}

export default NetworkError
