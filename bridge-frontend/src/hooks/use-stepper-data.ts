import { useState } from 'react';

export enum Step {
  Loading = 'loading',
  Login = 'login',
  Deposit = 'deposit',
  ConfirmDeposit = 'confirmDeposit',
  FinalizeDepositForm = 'finalizeDepositForm',
  ConfirmOracle = 'confirmOracle',
  ConfirmDepositChain = 'confirmDepositChain',
  ConfirmWithdrawChain = 'ConfirmWithdrawChain',
  PendingDeposit = 'pendingDeposit',
  Withdraw = 'withdraw',
  ConfirmWithdraw = 'confirmWithdraw',
  PendingWithdraw = 'pendingWithdraw',
  ConfirmAssertion = 'confirmAssertion',
  BatchAppend = 'batchAppend',
  CreateAssertion = 'createAssertion',
  FinalizeWithdrawForm= 'finalizeWithdrawForm',
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
