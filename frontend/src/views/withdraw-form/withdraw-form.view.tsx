import { useEffect, useRef } from 'react'
import { BigNumber } from 'ethers'
import { formatUnits } from 'ethers/lib/utils'

import useWithdrawFormStyles from './withdraw-form.styles'
import useWithdrawFormData from '../../hooks/use-withdraw-form-data'
import Header from '../shared/header/header.view'
import { ReactComponent as InfoIcon } from '../../images/info-icon.svg'

interface WithdrawFormProps {
  wallet: any,
  depositData: any,
  onAmountChange: () => void,
  onSubmit: (amount: BigNumber) => void,
  onDisconnectWallet: () => void
}

function WithdrawForm ({
  wallet,
  depositData,
  onAmountChange,
  onSubmit,
  onDisconnectWallet
}: WithdrawFormProps) {
  const { values, amounts, error, changeWithdrawValue } = useWithdrawFormData()
  const classes = useWithdrawFormStyles({ error: error ?? false })

  const inputEl = useRef<HTMLInputElement>(null)

  useEffect(() => {
    if (inputEl.current) {
      inputEl.current.focus()
    }
  }, [inputEl])

  useEffect(() => {
    if (!amounts.from.eq(BigNumber.from(0))) {
      onAmountChange()
    }
  }, [amounts, onAmountChange])

  return (
    <div className={classes.withdrawForm}>
      <Header
        address={wallet.address}
        title={`xDAI → ETH`}
        onDisconnectWallet={onDisconnectWallet}
      />
      
      <form
        className={classes.form}
        onSubmit={(event) => {
          event.preventDefault()
          onSubmit(amounts.from)
        }}
      >
        <div className={classes.fromInputGroup}>
          <p className={classes.fromTokenSymbol}>
            {"ETH"}
          </p>
          <input
            ref={inputEl}
            className={classes.fromInput}
            placeholder='0.0'
            value={values.from}
            onChange={event => changeWithdrawValue(event.target.value)}
          />
          <p className={classes.toValue}>
            {formatUnits(amounts.to, '1')} {'ETH'}
          </p>
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