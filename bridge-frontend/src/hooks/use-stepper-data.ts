import { useState } from 'react';

export enum Step {
  Loading = 'loading',
  Login = 'login',
  Deposit = 'deposit',
  ConfirmDeposit = 'confirmDeposit',
  FinalizeDepositForm = 'finalizeDepositForm',
  FinalizingDeposit = 'finalizingDeposit',
  ConfirmOracle = 'confirmOracle',
  ConfirmDepositChain = 'confirmDepositChain',
  ConfirmWithdrawChain = 'ConfirmWithdrawChain',


  Withdraw = 'withdraw',
  ConfirmWithdraw = 'confirmWithdraw',
  PendingDeposit = 'pendingDeposit',
  PendingWithdraw = 'pendingWithdraw',
  ConfirmAssertion = 'confirmAssertion',
  FinalizeWithdrawForm= 'finalizeWithdrawForm',
  PendingFinalizeDeposit = 'pendingFinalizeDeposit',
  FinalizingWithdraw = 'finalizingWithdraw',
  PendingFinalizeWithdraw = 'pendingFinalizeWithdraw',
  FinalizeDeposit = 'finalizeDeposit',
  FinalizeWithdrawl = 'finalizeWithdrawl',
  Overview = 'overview',
}

interface StepperData {
  step: Step;
  switchStep: (step: Step) => void;
}

function useStepperData(): StepperData {
  const [step, setStep] = useState<Step>(Step.Login);

  const switchStep = (step: Step): void => {
    setStep(step);
  };

  return { step, switchStep };
}

export default useStepperData;
