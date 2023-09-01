import { useState, useCallback } from 'react'
import useStepperStyles from './stepper.styles'
import useWallet from '../../hooks/use-wallet'
import useStep, { Step } from '../../hooks/use-stepper-data'
import Login from '../login/login.view'
import DepositForm from '../deposit-form/deposit-form.view'
import FinalizeDepositForm from '../finalize-deposit-form/finalize-deposit-form.view'
import FinalizeWithdrawForm from '../finalize-withdraw-form/finalize-withdraw-form.view'
import WithdrawForm from '../withdraw-form/withdraw-form.view'
import TxConfirm from '../tx-confirm/tx-confirm.view'
import TxConfirmChain from '../tx-confirm-chain/tx-confirm-chain.view'
import useDeposit from '../../hooks/use-deposit'
import useWithdraw from '../../hooks/use-withdraw'
import * as React from 'react';
import useFinalizeDeposit from '../../hooks/use-finalize-deposit'
import useFinalizeWithdraw from '../../hooks/use-finalize-withdraw'
import TxPendingOracleConfirmation from '../tx-pending-oracle-confirmation/tx-pending-oracle-confirmation.view'
import TxPendingDeposit from '../tx-pending-deposit/tx-pending-deposit.view'
import TxPendingWithdraw from '../tx-pending-withdraw/tx-pending-withdraw.view'
import TxConfirmAssertion from '../tx-confirm-assertion/tx-confirm-assertion.view'
import TxCreateAssertion from '../tx-create-assertion/tx-create-assertion.view'
import TxBatchAppend from '../tx-batch-append/tx-batch-append.view'
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
  SPECULAR_RPC_URL
} from "../../constants";
import type { PendingDeposit, PendingWithdrawal } from "../../types";


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
      version: 0,
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

  interface PendingWithdrawlData {
    status: string;
    data: PendingWithdrawal;
  }
  const PENDINGWITHDRAW: PendingWithdrawal = {
    l2BlockNumber: 0,
    proofL2BlockNumber: undefined,
    inboxSize: undefined,
    assertionID: undefined,
    withdrawalHash: "",
    withdrawalTx: {
      version: 0,
      nonce:ethers.BigNumber.from("0"),
      sender:"",
      target: "",
      value:ethers.BigNumber.from("0"),
      gasLimit:ethers.BigNumber.from("0"),
      data: ""
    }
  };

  const INITIALPENDINGWITHDRAW = {status: 'pending', data: PENDINGWITHDRAW}
  const [pendingWithdraw, setPendingWithdraw] = useState<PendingWithdrawlData>(INITIALPENDINGWITHDRAW);
  const [isDeposit, setIsDeposit] = useState<boolean>(true);



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
                  onSubmit={(fromAmount, selectedTokenKey) => {
                    deposit(wallet, fromAmount,selectedTokenKey)
                    setAmount(fromAmount)
                    switchStep(Step.ConfirmDeposit)
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
                  onGoToNextStep={() => switchStep(Step.PendingDeposit)}
                />
              )
            }
            case Step.PendingDeposit: {
              console.log("PendingDeposit")
              return (
                <TxPendingDeposit
                  wallet={wallet}
                  depositData={depositData}
                  l1Provider={l1Provider}
                  pendingDeposit={pendingDeposit}
                  setPendingDeposit={setPendingDeposit}
                  onGoBack={() => switchStep(Step.Deposit)}
                  onGoToFinalizeStep={() => {
                    switchStep(Step.ConfirmOracle)
                  }}
                />
              )
            }
            case Step.ConfirmOracle: {
              console.log("ConfirmDeposit tab")
              return (
                <TxPendingOracleConfirmation
                  wallet={wallet}
                  depositData={depositData}
                  pendingDeposit={pendingDeposit}
                  setPendingDeposit={setPendingDeposit}
                  switchChain={switchChain}
                  onGoToNextStep={() => switchStep(Step.FinalizeDepositForm)}
                />
              )
            }

            case Step.FinalizeDepositForm: {
              console.log("FinalizeDeposit tab")
              return (
                <FinalizeDepositForm
                  wallet={wallet}
                  onSubmit={() => {
                    switchChain(SPECULAR_NETWORK_ID.toString())
                    switchStep(Step.ConfirmDepositChain)
                  }}
                  onDisconnectWallet={disconnectWallet}
                />
              )
            }
            case Step.ConfirmDepositChain: {
              console.log("ConfirmChain tab")
              return (
                <TxConfirmChain
                  wallet={wallet}
                  networkId={SPECULAR_NETWORK_ID.toString()}
                  onGoBack={() => switchStep(Step.Deposit)}
                  onGoToNextStep={() => {
                    finalizeDeposit(wallet,amount,pendingDeposit,setPendingDeposit)
                    switchStep(Step.FinalizeDeposit)
                  }}
                />
              )
            }
            case Step.FinalizeDeposit: {
              console.log("FinalizeDeposit")
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
                  onSubmit={(fromAmount, selectedTokenKey) => {
                    withdraw(wallet, fromAmount, selectedTokenKey)
                    switchStep(Step.ConfirmWithdraw)
                  }}
                  onDisconnectWallet={disconnectWallet}
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
                  onGoToNextStep={() => switchStep(Step.PendingWithdraw)}
                />
              )
            }
            case Step.PendingWithdraw: {
              console.log("PendingWithdraw")
              return (
                <TxPendingWithdraw
                  wallet={wallet}
                  withdrawData={withdrawData}
                  l2Provider={l2Provider}
                  pendingWithdraw={pendingWithdraw}
                  setPendingWithdraw={setPendingWithdraw}
                  onGoBack={() => switchStep(Step.Deposit)}
                  onGoToFinalizeStep={() => {
                    switchStep(Step.BatchAppend)
                  }}
                />
              )
            }
            case Step.BatchAppend: {
              console.log("BatchAppend")
              return (
                <TxBatchAppend
                  wallet={wallet}
                  withdrawData={withdrawData}
                  pendingWithdraw={pendingWithdraw}
                  setPendingWithdraw={setPendingWithdraw}
                  switchChain={switchChain}
                  onGoBack={() => switchStep(Step.Withdraw)}
                  onGoToNextStep={() => {
                    switchStep(Step.CreateAssertion)
                  }}
                />
              )
            }
            case Step.CreateAssertion: {
              console.log("CreateAssertion")
              return (
                <TxCreateAssertion
                  wallet={wallet}
                  withdrawData={withdrawData}
                  pendingWithdraw={pendingWithdraw}
                  setPendingWithdraw={setPendingWithdraw}
                  switchChain={switchChain}
                  onGoBack={() => switchStep(Step.Withdraw)}
                  onGoToNextStep={() => {
                    switchStep(Step.ConfirmAssertion)
                  }}
                />
              )
            }
            case Step.ConfirmAssertion: {
              console.log("ConfirmAssertion")
              return (
                <TxConfirmAssertion
                  wallet={wallet}
                  withdrawData={withdrawData}
                  pendingWithdraw={pendingWithdraw}
                  setPendingWithdraw={setPendingWithdraw}
                  switchChain={switchChain}
                  onGoBack={() => switchStep(Step.Withdraw)}
                  onGoToNextStep={() => {
                    switchStep(Step.FinalizeWithdrawForm)
                  }}
                />
              )
            }
            case Step.FinalizeWithdrawForm: {
              console.log("FinalizeDeposit tab")
              return (
                <FinalizeWithdrawForm
                  wallet={wallet}
                  onSubmit={() => {
                    switchChain(CHIADO_NETWORK_ID.toString())
                    switchStep(Step.ConfirmWithdrawChain)
                  }}
                  onDisconnectWallet={disconnectWallet}
                />
              )
            }
            case Step.ConfirmWithdrawChain: {
              console.log("ConfirmChain tab")
              return (
                <TxConfirmChain
                  wallet={wallet}
                  networkId={CHIADO_NETWORK_ID.toString()}
                  onGoBack={() => switchStep(Step.Deposit)}
                  onGoToNextStep={() => {
                    setIsDeposit(false)
                    finalizeWithdraw(wallet,amount,pendingWithdraw.data)
                    switchStep(Step.FinalizeWithdrawl)
                  }}
                />
              )
            }
            case Step.FinalizeWithdrawl: {
              console.log("FinalizeWithdrawl")
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
                  isDeposit={isDeposit}
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
