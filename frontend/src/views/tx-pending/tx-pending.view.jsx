import { useEffect } from 'react'
import Header from '../shared/header/header.view'
import useTxPendingStyles from './tx-pending.styles'
import { ReactComponent as LinkIcon } from '../../images/link-icon.svg'
import Spinner from '../shared/spinner/spinner.view'

import { NETWORKS } from '../../chains'

function TxPending ({ wallet, depositData, onGoBack, onGoToOverviewStep }) {
  const classes = useTxPendingStyles()

  useEffect(() => {
    if (depositData.status === 'failed') {
      onGoBack()
    }
    if (depositData.status === 'successful') {
      onGoToOverviewStep()
    }
  }, [depositData, onGoBack, onGoToOverviewStep])

  return (
    <div className={classes.txOverview}>
      <Header
        address={wallet.address}
        title={`xDAI → ETH`}
      />
      <div className={classes.spinnerWrapper}>
        <Spinner className={classes.title} />
      </div>
      <div className={classes.buttonGroup}>
        <a
          className={classes.button}
          href={`${NETWORKS[wallet.chainId].blockExplorerUrl}/tx/${depositData?.data?.hash}`}
          target='_blank'
          rel='noopener noreferrer'
        >
          Check transaction status here
          <LinkIcon className={classes.buttonIcon} />
        </a>
      </div>
    </div>
  )
}

export default TxPending
