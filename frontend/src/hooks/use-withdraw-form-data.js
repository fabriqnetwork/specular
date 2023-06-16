import { useEffect, useState } from 'react'
import { BigNumber, ethers } from 'ethers'
import { formatUnits, parseUnits } from 'ethers/lib/utils'

// ethers BigNumber doesn't support decimals, so we need to workaround it
// using 35 as the SWAP_FACTOR instead of 3.5 and later divide the result by 10
// in the multiplyAmountBySwapFactor function
const INITIAL_VALUES = { from: '', to: '' }
const INITIAL_AMOUNTS = { from: BigNumber.from(0), to: BigNumber.from(0) }

function useWithdrawFormData (wallet) {
  const [values, setValues] = useState(INITIAL_VALUES)
  const [error, setError] = useState()
  const [amounts, setAmounts] = useState(INITIAL_AMOUNTS)

  useEffect(() => {
    setValues(INITIAL_VALUES)
    setAmounts(INITIAL_AMOUNTS)
    setError()
  }, [wallet])

 
  const changeValue = (newFromValue) => {
  }

  const convertAll = () => {
  }

  return { values, amounts, error, setError, changeValue, convertAll }
}

export default useWithdrawFormData
