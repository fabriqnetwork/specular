import { useState } from 'react';

export enum Step {
  Loading = 'loading',
  Login = 'login',
  Deposit = 'deposit',
  Withdraw = 'withdraw',
  ConfirmDeposit = 'confirmDeposit',
  ConfirmWithdraw = 'confirmWithdraw',
  PendingDeposit = 'pendingDeposit',
  PendingWithdraw = 'pendingWithdraw',
  Overview = 'overview',
  FinalizeDeposit = 'finalizeDeposit',
  FinalizeWithdrawl = 'finalizeWithdrawl',
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
