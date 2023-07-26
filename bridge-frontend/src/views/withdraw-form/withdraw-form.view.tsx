import { useEffect, useRef } from 'react'
import { BigNumber } from 'ethers'
import { formatUnits } from 'ethers/lib/utils'
import * as React from 'react';
import useWithdrawFormStyles from './withdraw-form.styles'
import useWithdrawFormData from '../../hooks/use-withdraw-form-data'
import Header from '../shared/header/header.view'
import InfoIcon from '@mui/icons-material/Info';
import DownArrow from '@mui/icons-material/ArrowDownward';

import { NETWORKS } from '../../chains';
import {SPECULAR_NETWORK_ID, CHIADO_NETWORK_ID} from '../../constants';

interface WithdrawFormProps {
  wallet: {
    address: string;
    chainId: number;
    provider: any;
  },
  withdrawData: any,
  onAmountChange: () => void,
  l1Provider: any,
  l2Provider: any,
  onSubmit: (amount: BigNumber) => void,
  onDisconnectWallet: () => void
}

function WithdrawForm ({
  wallet,
  withdrawData,
  onAmountChange,
  l1Provider,
  l2Provider,
  onSubmit,
  onDisconnectWallet
}: WithdrawFormProps) {
  const { values, amounts, l1balance, l2balance, error, changeWithdrawValue } = useWithdrawFormData(wallet, l1Provider, l2Provider)
  const classes = useWithdrawFormStyles({ error: !!error })

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
        <div className={classes.card}>
          <p className={classes.cardTitleText}>
          {NETWORKS[SPECULAR_NETWORK_ID].name+" "+NETWORKS[SPECULAR_NETWORK_ID].nativeCurrency.symbol}
          </p>
          <input
            ref={inputEl}
            className={classes.fromInput}
            placeholder='0.0'
            value={values.from}
            onChange={event => changeWithdrawValue(event.target.value)}
          />
          <p className={classes.toValue}>
            Balance: {formatUnits(l2balance, NETWORKS[SPECULAR_NETWORK_ID].nativeCurrency.decimals)} {NETWORKS[SPECULAR_NETWORK_ID].nativeCurrency.symbol}
          </p>
        </div>
        <DownArrow className={classes.cardIcon} />
        <div className={classes.card}>
          <p className={classes.cardTitleText}>
          {NETWORKS[CHIADO_NETWORK_ID].name+" "+NETWORKS[CHIADO_NETWORK_ID].nativeCurrency.symbol}
          </p>
          <p>
            {formatUnits(amounts.to, NETWORKS[CHIADO_NETWORK_ID].nativeCurrency.decimals)} {NETWORKS[CHIADO_NETWORK_ID].nativeCurrency.symbol}
          </p>
          <p className={classes.toValue}>
            Balance: {formatUnits(l1balance, NETWORKS[CHIADO_NETWORK_ID].nativeCurrency.decimals)} {NETWORKS[CHIADO_NETWORK_ID].nativeCurrency.symbol}
          </p>
        </div>
        {(error || withdrawData.status === 'failed') && (
          <div className={classes.inputErrorContainer}>
            <InfoIcon className={classes.cardErrorIcon} />
            <p>{error || withdrawData.error}</p>
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
