import { useState, useEffect } from 'react'
import { Contract } from 'ethers'

import ERC20ABI from '../abis/erc20'

function useTokenContract (address, wallet) {
  const [contract, setContract] = useState()

  useEffect(() => {
    if (address && wallet?.provider) {
      const contract = new Contract(address, ERC20ABI, wallet.provider.getSigner(0))
      setContract(contract)
    }
  }, [address, wallet])

  return contract
}

export default useTokenContract
