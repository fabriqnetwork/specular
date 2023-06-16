import { useState, useEffect } from 'react'

function useTokenBalance (address, contract) {
  const [balance, setBalance] = useState()

  useEffect(() => {
    if (address && contract) {
      contract.balanceOf(address).then(setBalance).catch(() => setBalance())
    } else {
      setBalance()
    }
  }, [address, contract])

  return balance
}

export default useTokenBalance
