import { useEffect, useRef } from 'react'
import { BigNumber } from 'ethers'
import { formatUnits } from 'ethers/lib/utils'

import useWithdrawFormStyles from './withdraw-form.styles'
import useWithdrawFormData from '../../hooks/use-withdraw-form-data'
import Header from '../shared/header/header.view'
import { ReactComponent as InfoIcon } from '../../images/info-icon.svg'

function WithdrawForm ({
  wallet,
  depositData,
  onAmountChange,
  onSubmit,
  onDisconnectWallet,
}) {
  const { values, amounts, error, depositAll, changeValue } = useWithdrawFormData(wallet)
  const classes = useWithdrawFormStyles({ error })
  const inputEl = useRef()

  useEffect(() => {
    if (inputEl) {
      inputEl.current.focus()
    }
  }, [inputEl])

  useEffect(() => {
    if (!amounts.from.eq(BigNumber.from(0))) {
      onAmountChange()
    }
  }, [amounts, onAmountChange])

  return (
    <div className={classes.swapForm}>
      <Header
        address={wallet.address}
        title={`xDai → ETH`}
        onDisconnectWallet={onDisconnectWallet}
      />
      <div className={classes.balanceCard}>
        
        <button
          className={classes.depositAllButton}
          type='button'
          onClick={depositAll}
        >
          Withdraw All
        </button>
      </div>
      <form
        className={classes.form}
        onSubmit={(event) => {
          event.preventDefault()
          onSubmit(amounts.from)
        }}
      >
        <div className={classes.fromInputGroup}>
          <input
            ref={inputEl}
            className={classes.fromInput}
            placeholder='0.0'
            value={values.from}
            onChange={event => changeValue(event.target.value)}
          />
          
        </div>
        {(error || depositData.status === 'failed') && (
          <div className={classes.inputErrorContainer}>
            <InfoIcon className={classes.inputErrorIcon} />
            <p>{error || depositData.error}</p>
          </div>
        )}
        <button
          className={classes.submitButton}
          disabled={amounts.from.eq(BigNumber.from(0)) || !!error}
          type='submit'
        >
          Withdraw
        </button>
      </form>
    </div>
  )
}

export default WithdrawForm