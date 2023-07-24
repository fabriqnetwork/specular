import { useEffect, useRef } from 'react'
import { BigNumber } from 'ethers'
import { formatUnits } from 'ethers/lib/utils'

import useDepositFormStyles from './deposit-form.styles'
import useDepositFormData from '../../hooks/use-deposit-form-data'
import Header from '../shared/header/header.view'
import InfoIcon from '@mui/icons-material/Info';
import DownArrow from '@mui/icons-material/ArrowDownward';

interface DepositFormProps {
  wallet: {
    address: string;
    chainId: number;
    provider: any;
  },
  depositData: any,
  onAmountChange: () => void,
  l1Provider: any,
  l2Provider: any,
  onSubmit: (amount: BigNumber) => void,
  onDisconnectWallet: () => void
}

function DepositForm ({
  wallet,
  depositData,
  onAmountChange,
  l1Provider,
  l2Provider,
  onSubmit,
  onDisconnectWallet
}: DepositFormProps) {
  const { values, amounts, l1balance, l2balance, error, changeDepositValue } = useDepositFormData(wallet, l1Provider, l2Provider)
  const classes = useDepositFormStyles({ error: !!error })

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

    <div className={classes.depositForm}>
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
        <div className={classes.card}>
          <p className={classes.cardTitleText}>
            {"Chiado xDai"}
          </p>
          <input
            ref={inputEl}
            className={classes.fromInput}
            placeholder='0.0'
            value={values.from}
            onChange={event => changeDepositValue(event.target.value)}
          />
          <p className={classes.toValue}>
            Balance: {formatUnits(l1balance, 18)} {'xDai'}
          </p>
        </div>
        <DownArrow className={classes.cardIcon} />
        <div className={classes.card}>
          <p className={classes.cardTitleText}>
            {"Specular ETH"}
          </p>
          <p>
            {formatUnits(amounts.to, 18)} {'ETH'}
          </p>
          <p className={classes.toValue}>
            Balance: {formatUnits(l2balance, 18)} {'ETH'}
          </p>
        </div>
        {(error || depositData.status === 'failed') && (
          <div className={classes.inputErrorContainer}>
            <InfoIcon className={classes.cardErrorIcon} />
            <p>{error || depositData.error}</p>
          </div>
        )}
        <button
          className={classes.submitButton}
          disabled={amounts.from.eq(BigNumber.from(0)) || !!error}
          type='submit'
        >
          Deposit
        </button>
      </form>
    </div>
  )
}

export default DepositForm
