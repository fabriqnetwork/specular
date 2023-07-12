import { useEffect, useState } from 'react';
import { BigNumber, ethers } from 'ethers';
import {SPECULAR_RPC_URL} from '../constants';

const INITIAL_VALUES = { from: '', to: '' };
const INITIAL_AMOUNTS = { from: BigNumber.from(0), to: BigNumber.from(0) };

interface DepositFormData {
  values: { from: string; to: string };
  amounts: { from: BigNumber; to: BigNumber };
  l1balance:BigNumber;
  l2balance:BigNumber;
  error: string | undefined;
  setError: (error?: string) => void;
  changeDepositValue: (newFromValue: string) => void;
}

function useDepositFormData(wallet: any): DepositFormData {
  const [values, setValues] = useState<{ from: string; to: string }>(INITIAL_VALUES);
  const [error, setError] = useState<string>();
  const [amounts, setAmounts] = useState<{ from: BigNumber; to: BigNumber }>(INITIAL_AMOUNTS);
  const [l1balance, setl1balances] = useState<BigNumber>(BigNumber.from(0));
  const [l2balance, setl2balances] = useState<BigNumber>(BigNumber.from(0));


  const GetL1Balance = async (wallet: any) => {
    if (wallet) {
      const balance = await wallet.provider.getBalance(wallet.address);

      return balance;
    }
    return BigNumber.from(0);
  };

  const GetL2Balance = async (wallet: any) => {
      const l2Provider = new ethers.providers.StaticJsonRpcProvider(SPECULAR_RPC_URL)
      if (wallet) {
      const balance  = await l2Provider.getBalance(wallet.address);
      return balance;
  }
  return BigNumber.from(0);
  };

  useEffect(() => {
    setValues(INITIAL_VALUES);
    setAmounts(INITIAL_AMOUNTS);
    setError(undefined);

    const fetchL1Balance = async () => {
      const balance = await GetL1Balance(wallet);
      setl1balances(balance);
    };

    const fetchL2Balance = async () => {
      const balance = await GetL2Balance(wallet);
      setl2balances(balance);
    };

    fetchL1Balance();
    fetchL2Balance();
  }, [wallet]);

  const changeDepositValue = (newFromValue: string): void => {
    const INPUT_REGEX = new RegExp(`^\\d*(?:\\.\\d{0,${18}})?$`);
    if (INPUT_REGEX.test(newFromValue)) {
      try {
        const newFromAmount = ethers.utils.parseUnits(newFromValue.length > 0 ? newFromValue : '0', 18);
        const newToAmount = newFromAmount;

        setAmounts({ from: newFromAmount, to: newToAmount });
        setValues({ from: newFromValue, to: ethers.utils.formatUnits(newToAmount, 18) });
        if (newFromAmount.gt(l1balance)) {
          setError("You don't have enough funds");
        } else {
          setError(undefined);
        }
      } catch (err) {
        console.log(err);
      }
    }
  };

  return { values, amounts, l1balance,l2balance, error, setError, changeDepositValue };
}

export default useDepositFormData;
