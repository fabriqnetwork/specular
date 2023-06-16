import { useState } from 'react'
import { ethers } from 'ethers'

import {
  L1PORTAL_ADDRESS,
} from "../constants";


const INITIAL_DATA = { status: 'pending' }

function useDeposit (address, provider) {
  const [data, setData] = useState(INITIAL_DATA)

  const deposit = async (wallet, amount) => {
    setData({ status: 'loading' })

    if (!wallet) {
      setData({ status: 'failed', error: 'Wallet doesn\'t exist' })
      return
    }

    try {
      const balance = await provider.getBalance(address);
      const formattedBalance = ethers.utils.formatEther(balance);
      const signer = await provider.getSigner();
      const tx = await signer.sendTransaction({
        to: L1PORTAL_ADDRESS,
        value: ethers.utils.parseUnits(
          Number(amount.value).toString(),
          ethers.utils.parseEther(formattedBalance)
        ),
      });

      setData({ status: 'pending', data: tx })
      await tx.wait()
      setData({ status: 'successful', data: tx })
    } catch (err) {
      let error = 'Transaction failed.'
      if (err?.code === -32603) {
        error = 'Transaction was not sent because of the low gas price. Try to increase it.'
      }
      setData({ status: 'failed', error })
      console.log(err)
    }
  }

  const resetData = () => {
    setData(INITIAL_DATA)
  }

  return { deposit, data, resetData }
}

export default useDeposit
