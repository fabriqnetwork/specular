import { useState, useCallback } from 'react'
import useStepperStyles from './stepper.styles'
import useWallet from '../../hooks/use-wallet'
import useStep, { Step } from '../../hooks/use-stepper-data'
import Login from '../login/login.view'
import WithdrawForm from '../withdraw-form/withdraw-form.view'
import DepositForm from '../deposit-form/deposit-form.view'
import TxConfirm from '../tx-confirm/tx-confirm.view'
import useDeposit from '../../hooks/use-deposit'
import TxPending from '../tx-pending/tx-pending.view'
import TxOverview from '../tx-overview/tx-overview.view'
import NetworkError from '../network-error/network-error.view'
import DataLoader from '../data-loader/data-loader'


function Stepper () {
  const classes = useStepperStyles()
  const { wallet, loadWallet, disconnectWallet, isMetamask, switchChainInMetaMask } = useWallet()
  const { step, switchStep } = useStep()
  const { deposit, data: depositData, resetData: resetDepositData } = useDeposit(wallet?.address, wallet)


  const tabs = [
    { name: 'Deposit', step: Step.Deposit },
  ]
  tabs.push({name: 'Withdraw', step: Step.Withdraw })
  

  const [activeTab, setActiveTab] = useState(tabs[0].name)

  const selectTab = useCallback((tab) => {
    if (activeTab === tab.name) return
    setActiveTab(tab.name)
    switchStep(tab.step)
  }, [activeTab, switchStep])

  if (wallet && wallet.chainId !== process.env.REACT_APP_NETWORK_ID) {
    return (
      <div className={classes.stepper}>
        <NetworkError {...{ isMetamask, switchChainInMetaMask }} />
      </div>
    )
  }

  return (
    <div className={classes.container}>
      {![Step.Login, Step.Loading].includes(step) && (
        <div className={classes.tabs}>
          {tabs.map(tab =>
            <button
              key={tab.name}
              className={activeTab === tab.name ? classes.tabActive : classes.tab}
              onClick={() => selectTab(tab)}
              disabled={![Step.Withdraw, Step.Deposit].includes(step)}
            >
              <span className={classes.tabName}>{tab.name}</span>
            </button>
          )}
        </div>
      )}
      <div className={classes.stepper}>
        {(() => {
          switch (step) {
            case Step.Loading: {
              console.log("Loading attempted")
              return (
                <DataLoader
                onGoToNextStep={() => switchStep(Step.Deposit)}
                />
              )
            }
            case Step.Login: {
              console.log("Login attempted")
              return (
                <Login
                  wallet={wallet}
                  onLoadWallet={loadWallet}
                  onGoToNextStep={() => switchStep(Step.Deposit)}
                />
              )
            }
            case Step.Deposit: {
              console.log("Deposit attempted")
              return (
                <DepositForm
                  wallet={wallet}
                  depositData={depositData}
                  onAmountChange={resetDepositData}
                  onSubmit={(fromAmount) => {
                    deposit(wallet, fromAmount)
                    switchStep(Step.Confirm)
                  }}
                  onDisconnectWallet={disconnectWallet}
                />
              )
            }
            case Step.Withdraw: {
              console.log("Withdraw attempted")
              return (
                <WithdrawForm
                  wallet={wallet}
                  depositData={depositData}
                  onAmountChange={resetDepositData}
                  onSubmit={(fromAmount) => {
                    deposit(wallet, fromAmount)
                    switchStep(Step.Confirm)
                  }}
                  onDisconnectWallet={disconnectWallet}
                />
              )
            }
            case Step.Confirm: {
              console.log("Tx Confirmed")
              return (
                <TxConfirm
                  wallet={wallet}
                  depositData={depositData}
                  onGoBack={() => switchStep(Step.Deposit)}
                  onGoToPendingStep={() => switchStep(Step.Pending)}
                />
              )
            }
            case Step.Pending: {
              console.log("Pending")
              return (
                <TxPending
                  wallet={wallet}
                  depositData={depositData}
                  onGoBack={() => switchStep(Step.Deposit)}
                  onGoToOverviewStep={() => switchStep(Step.Overview)}
                />
              )
            }
            case Step.Overview: {
              console.log("Overview")
              return (
                <TxOverview
                  wallet={wallet}
                  depositData={depositData}
                  onGoBack={() => window.location.reload()}
                  onDisconnectWallet={disconnectWallet}
                  isMetamask={isMetamask}
                />
              )
            }
            default: {
              return <></>
            }
          }
        })()}

      </div>
    </div>
  )
}

export default Stepper
