import { useEffect, useState } from 'react'
import { BigNumber, ethers } from 'ethers'
import { formatUnits, parseUnits } from 'ethers/lib/utils'

// ethers BigNumber doesn't support decimals, so we need to workaround it
// using 35 as the SWAP_FACTOR instead of 3.5 and later divide the result by 10
// in the multiplyAmountBySwapFactor function
const INITIAL_VALUES = { from: '', to: '' }
const INITIAL_AMOUNTS = { from: BigNumber.from(0), to: BigNumber.from(0) }

function useWithdrawFormData () {
  const [values, setValues] = useState(INITIAL_VALUES)
  const [error, setError] = useState()
  const [amounts, setAmounts] = useState(INITIAL_AMOUNTS)

  useEffect(() => {
    setValues(INITIAL_VALUES)
    setAmounts(INITIAL_AMOUNTS)
  }, [])

  const setToAmount = (value: ethers.BigNumber): ethers.BigNumber => {
    return value.mul("10").div(ethers.constants.WeiPerEther);
  };

  const changeWithdrawValue = (newFromValue: string): void => {
    console.log("Event is " + newFromValue);
    const INPUT_REGEX = new RegExp(`^\\d*(?:\\.\\d{0,${18}})?$`);
    if (INPUT_REGEX.test(newFromValue)) {
      try {
        const newFromAmount = ethers.utils.parseUnits(newFromValue.length > 0 ? newFromValue : '0', 18);
        const newToAmount = setToAmount(newFromAmount);
  
        setAmounts({ from: newFromAmount, to: newToAmount });
        setValues({ from: newFromValue, to: ethers.utils.formatUnits(newToAmount, 18) });
      } catch (err) {
        console.log(err);
      }
    }
  };



  return { values, amounts, error, setError, changeWithdrawValue}
}

export default useWithdrawFormData
