import { useEffect,useState, useCallback } from 'react'
import useStepperStyles from './stepper.styles'
import useWallet from '../../hooks/use-wallet'
import useStep, { Step } from '../../hooks/use-stepper-data'
import Login from '../login/login.view'
import DepositForm from '../deposit-form/deposit-form.view'
import WithdrawForm from '../withdraw-form/withdraw-form.view'
import TxConfirm from '../tx-confirm/tx-confirm.view'
import useDeposit from '../../hooks/use-deposit'
import useWithdraw from '../../hooks/use-withdraw'
import * as React from 'react';
import useFinalizeDeposit from '../../hooks/use-finalize-deposit'
import useFinalizeWithdraw from '../../hooks/use-finalize-withdraw'
import TxPending from '../tx-pending/tx-pending.view'
import TxPendingFinalizeDeposit from '../tx-pending-finalize-deposit/tx-pending-finalize-deposit.view'
import TxPendingFinalizeWithdraw from '../tx-pending-finalize-withdraw/tx-pending-finalize-withdraw.view'
import TxFinalizeDeposit from '../tx-finalize-deposit/tx-finalize-deposit.view'
import TxFinalizeWithdraw from '../tx-finalize-withdraw/tx-finalize-withdraw.view'
import TxOverview from '../tx-overview/tx-overview.view'
import NetworkError from '../network-error/network-error.view'
import { ethers } from 'ethers'
import DataLoader from '../data-loader/data-loader'
import {
  CHIADO_NETWORK_ID,
  SPECULAR_NETWORK_ID,
  CHIADO_RPC_URL,
  SPECULAR_RPC_URL,
  L1ORACLE_ADDRESS,
  L1PORTAL_ADDRESS
} from "../../constants";
import type { PendingDeposit, PendingWithdrawal } from "../../types";
import {
  L1Oracle__factory,
  IL1Portal__factory,
} from "../../typechain-types";

function Stepper () {
  const classes = useStepperStyles()
  const { wallet, loadWallet, disconnectWallet, isMetamask, switchChain } = useWallet()
  const { step, switchStep } = useStep()
  const { deposit, data: depositData, resetData: resetDepositData } = useDeposit()
  const { finalizeDeposit, data: finalizeDepositData } = useFinalizeDeposit(switchChain)
  const { withdraw, data: withdrawData, resetData: resetWithdrawData } = useWithdraw()
  const { finalizeWithdraw, data: finalizeWithdrawData } = useFinalizeWithdraw(switchChain)
  const [amount, setAmount] = useState(ethers.BigNumber.from("0"));

  const PENDINGDEPOSIT: PendingDeposit = {
    l1BlockNumber: 0,
    proofL1BlockNumber: undefined,
    depositHash: "",
    depositTx: {
      nonce:ethers.BigNumber.from("0"),
      sender:"",
      target: "",
      value:ethers.BigNumber.from("0"),
      gasLimit:ethers.BigNumber.from("0"),
      data: ""
    }
  };
  interface PendingData {
    status: string;
    data: PendingDeposit;
  }
  const INITIALPENDINGDEPOSIT = {status: 'pending', data: PENDINGDEPOSIT}
  const [pendingDeposit, setPendingDeposit] = useState<PendingData>(INITIALPENDINGDEPOSIT);

  const initialPendingWithdrawal: PendingWithdrawal = {
    l2BlockNumber: 0,
    proofL2BlockNumber: undefined,
    inboxSize: undefined,
    assertionID: undefined,
    withdrawalHash: "",
    withdrawalTx: {
      nonce:ethers.BigNumber.from("0"),
      sender:"",
      target: "",
      value:ethers.BigNumber.from("0"),
      gasLimit:ethers.BigNumber.from("0"),
      data: ""
    }
  };
  const [pendingWithdraw, setPendingWithdraw] = useState(initialPendingWithdrawal)


  const tabs = [
    { name: 'Deposit', step: Step.Deposit },
  ]
  tabs.push({name: 'Withdraw', step: Step.Withdraw })

  const [activeTab, setActiveTab] = useState(tabs[0].name)

  const selectTab = useCallback((tab: { name: string; step: Step }) => {
    if (activeTab === tab.name) return;
    setActiveTab(tab.name);
    switchStep(tab.step);
  }, [activeTab, switchStep]);


  const l1Provider = new ethers.providers.StaticJsonRpcProvider(CHIADO_RPC_URL);
  const l2Provider = new ethers.providers.StaticJsonRpcProvider(SPECULAR_RPC_URL);

useEffect(() => {

const l1Portal = IL1Portal__factory.connect(
  L1PORTAL_ADDRESS,
  l1Provider
);

  l1Portal.on(
    l1Portal.filters.DepositInitiated(),
    (nonce:ethers.BigNumber, sender:string, target:string, value:ethers.BigNumber, gasLimit:ethers.BigNumber, data:string, depositHash:string, event:any) => {
      console.log("Main L1 Portal transactionHash is "+event.transactionHash+"Deposit data hash is "+ depositData.status+" Deposit Data is "+depositData.data+" Deposit Data hash is "+depositData.data?.hash)
      if (event.transactionHash === depositData.data?.hash) {
        const newPendingDeposit: PendingDeposit = {
          l1BlockNumber: event.blockNumber,
          proofL1BlockNumber: undefined,
          depositHash: depositHash,
          depositTx: {
            nonce,
            sender,
            target,
            value,
            gasLimit,
            data,
          },
        }
        console.log("Main Correct L1 Portal transactionHash is "+event.transactionHash)
        setPendingDeposit({ status: 'initiated', data: newPendingDeposit});
      }
    }
  );
  },[depositData]
)



if (wallet && !(wallet.chainId == CHIADO_NETWORK_ID || wallet.chainId == SPECULAR_NETWORK_ID) ){
  return (
    <div className={classes.stepper}>
      <NetworkError {...{ isMetamask, switchChain }} />
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
              return (
                <DataLoader
                onGoToNextStep={() => switchStep(Step.Deposit)}
                />
              )
            }
            case Step.Login: {
              console.log("Login tab")
              return (
                <Login
                  wallet={wallet}
                  onLoadWallet={loadWallet}
                  onGoToNextStep={() => switchStep(Step.Deposit)}
                />
              )
            }
            case Step.Deposit: {
              console.log("Deposit tab")
              switchChain(CHIADO_NETWORK_ID.toString())
              console.log("Chain Id is: "+wallet.chainId)
              return (
                <DepositForm
                  wallet={wallet}
                  depositData={depositData}
                  onAmountChange={resetDepositData}
                  l1Provider={l1Provider}
                  l2Provider={l2Provider}
                  onSubmit={(fromAmount) => {
                    deposit(wallet, fromAmount)
                    setAmount(fromAmount)
                    switchStep(Step.ConfirmDeposit)
                  }}
                  onDisconnectWallet={disconnectWallet}
                />
              )
            }
            case Step.Withdraw: {
              console.log("Withdraw tab")
              switchChain(SPECULAR_NETWORK_ID.toString())
              console.log("Chain Id is: "+wallet.chainId)
              return (
                <WithdrawForm
                  wallet={wallet}
                  withdrawData={withdrawData}
                  onAmountChange={resetWithdrawData}
                  l1Provider={l1Provider}
                  l2Provider={l2Provider}
                  onSubmit={(fromAmount) => {
                    withdraw(wallet, fromAmount)
                    switchStep(Step.ConfirmWithdraw)
                  }}
                  onDisconnectWallet={disconnectWallet}
                />
              )
            }
            case Step.ConfirmDeposit: {
              console.log("ConfirmDeposit tab")
              return (
                <TxConfirm
                  wallet={wallet}
                  transactionData={depositData}
                  onGoBack={() => switchStep(Step.Deposit)}
                  onGoToPendingStep={() => switchStep(Step.PendingDeposit)}
                />
              )
            }
            case Step.ConfirmWithdraw: {
              console.log("ConfirmWithdraw tab")
              return (
                <TxConfirm
                  wallet={wallet}
                  transactionData={withdrawData}
                  onGoBack={() => switchStep(Step.Withdraw)}
                  onGoToPendingStep={() => switchStep(Step.PendingWithdraw)}
                />
              )
            }
            case Step.PendingDeposit: {
              console.log("PendingDeposit")
              return (
                <TxPending
                  wallet={wallet}
                  transactionData={depositData}
                  onGoBack={() => switchStep(Step.Deposit)}
                  onGoToFinalizeStep={() => {
                    switchStep(Step.PendingFinalizeDeposit)
                  }}
                />
              )
            }
            case Step.PendingWithdraw: {
              console.log("PendingWithdraw")
              return (
                <TxPending
                  wallet={wallet}
                  transactionData={withdrawData}
                  onGoBack={() => switchStep(Step.Withdraw)}
                  onGoToFinalizeStep={() => switchStep(Step.PendingFinalizeWithdraw)}
                />
              )
            }
            case Step.PendingFinalizeDeposit: {
              switchChain(SPECULAR_NETWORK_ID.toString())
              console.log("PendingFinalizeDeposit")
              return (
                <TxPendingFinalizeDeposit
                  wallet={wallet}
                  depositData={depositData}
                  pendingDeposit={pendingDeposit}
                  switchChain={switchChain}
                  onGoToFinalizeStep={() => {
                    switchStep(Step.PendingFinalizeDeposit1)
                  }}
                />
              )
            }
            case Step.PendingFinalizeDeposit1: {
              console.log("PendingDeposit")
              return (
                <TxPending
                  wallet={wallet}
                  transactionData={depositData}
                  onGoBack={() => switchStep(Step.Deposit)}
                  onGoToFinalizeStep={() => {
                    finalizeDeposit(wallet,amount,pendingDeposit,setPendingDeposit)
                    switchStep(Step.FinalizeDeposit)
                  }}
                />
              )
            }
            case Step.PendingFinalizeWithdraw: {
              console.log("PendingWithdraw")
              return (
                <TxPendingFinalizeWithdraw
                  wallet={wallet}
                  depositData={withdrawData}
                  setPendingWithdraw={setPendingWithdraw}
                  onGoToFinalizeStep={() => {
                    switchStep(Step.FinalizeWithdrawl)
                  }}
                />
              )
            }

            case Step.FinalizeDeposit: {
              // if (pendingDeposit.status !== 'finalized') {
              //   setPendingDeposit({ status: 'finalized', data: pendingDeposit.data})
              //   console.log("Chain Id is: "+wallet.chainId)
              //   console.log("pendingDeposit"+pendingDeposit)
              //   finalizeDeposit(wallet,amount,pendingDeposit,setPendingDeposit)
              //   console.log("pendingDeposit done")
              // }
              return (
                <TxFinalizeDeposit
                  wallet={wallet}
                  depositData={depositData}
                  finalizeDepositData={finalizeDepositData}
                  onGoBack={() => switchStep(Step.Deposit)}
                  onGoToOverviewStep={() => switchStep(Step.Overview)}
                />
              )
            }
            case Step.FinalizeWithdrawl: {
              console.log("FinalizeWithdrawl")
              finalizeWithdraw(wallet,amount,pendingWithdraw)
              return (
                <TxFinalizeWithdraw
                  wallet={wallet}
                  withdrawData={withdrawData}
                  finalizeWithdrawData={finalizeWithdrawData}
                  onGoBack={() => switchStep(Step.Withdraw)}
                  onGoToOverviewStep={() => switchStep(Step.Overview)}
                />
              )
            }
            case Step.Overview: {
              console.log("Overview")
              return (
                <TxOverview
                  wallet={wallet}
                  finalizeTransactionData={finalizeDepositData}
                  transactionData={depositData}
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
