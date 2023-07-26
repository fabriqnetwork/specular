import { useEffect, useState } from 'react';
import { BigNumber, ethers } from 'ethers';
import { NETWORKS } from '../chains';
import {SPECULAR_NETWORK_ID, CHIADO_NETWORK_ID} from '../constants';

const INITIAL_VALUES = { from: '', to: '' };
const INITIAL_AMOUNTS = { from: BigNumber.from(0), to: BigNumber.from(0) };

interface WithdrawFormData {
  values: { from: string; to: string };
  amounts: { from: BigNumber; to: BigNumber };
  l1balance:BigNumber;
  l2balance:BigNumber;
  error: string | undefined;
  setError: (error?: string) => void;
  changeWithdrawValue: (newFromValue: string) => void;
}

function useWithdrawFormData(wallet: any, l1Provider:any, l2Provider:any): WithdrawFormData {
  const [values, setValues] = useState<{ from: string; to: string }>(INITIAL_VALUES);
  const [error, setError] = useState<string>();
  const [amounts, setAmounts] = useState<{ from: BigNumber; to: BigNumber }>(INITIAL_AMOUNTS);
  const [l1balance, setl1balances] = useState<BigNumber>(BigNumber.from(0));
  const [l2balance, setl2balances] = useState<BigNumber>(BigNumber.from(0));

  useEffect(() => {
    setValues(INITIAL_VALUES);
    setAmounts(INITIAL_AMOUNTS);
    setError(undefined);

    const GetL1Balance = async (wallet: any) => {
      if (wallet) {
        const balance = await l1Provider.getBalance(wallet.address);

        return balance;
      }
      return BigNumber.from(0);
    };

    const GetL2Balance = async (wallet: any) => {
        if (wallet) {
        const balance  = await l2Provider.getBalance(wallet.address);
        return balance;
    }
    return BigNumber.from(0);
    };

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
  }, [wallet, l1Provider,l2Provider]);

  const changeWithdrawValue = (newFromValue: string): void => {
    const INPUT_REGEX = new RegExp(`^\\d*(?:\\.\\d{0,${NETWORKS[SPECULAR_NETWORK_ID].nativeCurrency.decimals}})?$`);
    if (INPUT_REGEX.test(newFromValue)) {
      try {
        const newFromAmount = ethers.utils.parseUnits(newFromValue.length > 0 ? newFromValue : '0', NETWORKS[SPECULAR_NETWORK_ID].nativeCurrency.decimals);
        const newToAmount = newFromAmount;

        setAmounts({ from: newFromAmount, to: newToAmount });
        setValues({ from: newFromValue, to: ethers.utils.formatUnits(newToAmount, NETWORKS[CHIADO_NETWORK_ID].nativeCurrency.decimals) });
        if (newFromAmount.gt(l2balance)) {
          setError("You don't have enough funds");
        } else {
          setError(undefined);
        }
      } catch (err) {
        console.log(err);
      }
    }
  };

  return { values, amounts, l1balance,l2balance, error, setError, changeWithdrawValue };
}

export default useWithdrawFormData;
