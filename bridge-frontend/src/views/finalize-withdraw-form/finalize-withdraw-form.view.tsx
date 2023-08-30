import { useEffect, useRef } from 'react'

import useFinalizeDepositFormStyles from './finalize-withdraw-form'
import Header from '../shared/header/header.view'

interface FinalizeDepositFormProps {
  wallet: {
    address: string;
    chainId: number;
    provider: any;
  },
  onSubmit: () => void,
  onDisconnectWallet: () => void
}

function FinalizeDepositForm ({
  wallet,
  onSubmit,
  onDisconnectWallet
}: FinalizeDepositFormProps) {
  const classes = useFinalizeDepositFormStyles()


  return (

    <div className={classes.finalizeDepositForm}>
      <Header
        address={wallet.address}
        title={`Finalize Withdraw`}
        onDisconnectWallet={onDisconnectWallet}
      />
      <form
        className={classes.form}
        onSubmit={(event) => {
          event.preventDefault()
          onSubmit()
        }}
      >
        <button
          className={classes.submitButton}
          type='submit'
        >
          Withdraw
        </button>
      </form>
    </div>
  )
}

export default FinalizeDepositForm
