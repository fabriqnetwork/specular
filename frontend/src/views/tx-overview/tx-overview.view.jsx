import Header from '../shared/header/header.view'
import useTxOverviewStyles from './tx-overview.styles'
import { ReactComponent as CheckIcon } from '../../images/check-icon.svg'
import { ReactComponent as LinkIcon } from '../../images/link-icon.svg'
import { ReactComponent as MetaMaskLogo } from '../../images/metamask-logo.svg'
import useWatchAsset from '../../hooks/use-watch-asset'

import { NETWORKS } from '../../chains'

function TxOverview ({ wallet, depositData, onGoBack, onDisconnectWallet, isMetamask }) {
  const classes = useTxOverviewStyles()
  const watchAsset = useWatchAsset()

  return (
    <div className={classes.txOverview}>
      <Header
        address={wallet.address}
        title={`xDAI → ETH`}
        onGoBack={onGoBack}
        onDisconnectWallet={onDisconnectWallet}
      />
      <CheckIcon className={classes.checkIcon} />
      <div className={classes.buttonGroup}>
        <a
          className={classes.button}
          href={`${NETWORKS[wallet.chainId].blockExplorerUrl}/tx/${depositData.data.hash}`}
          target='_blank'
          rel='noopener noreferrer'
        >
          Check transaction details here
          <LinkIcon className={classes.buttonIcon} />
        </a>
        {isMetamask && (
          <button
            className={classes.button}
          >
            Add xDAI token to MetaMask
            <MetaMaskLogo className={classes.buttonIcon} />
          </button>
        )}
      </div>
    </div>
  )
}

export default TxOverview
