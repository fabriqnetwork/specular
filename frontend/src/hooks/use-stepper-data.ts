import { useState } from 'react';

export enum Step {
  Loading = 'loading',
  Login = 'login',
  Deposit = 'deposit',
  Withdraw = 'withdraw',
  Confirm = 'confirm',
  Pending = 'pending',
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
