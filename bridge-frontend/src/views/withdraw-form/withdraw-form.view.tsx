import { useEffect, useRef, useState } from 'react'
import { BigNumber } from 'ethers'
import { formatUnits } from 'ethers/lib/utils'
import {Select, MenuItem} from "@mui/material";
import useWithdrawFormStyles from './withdraw-form.styles'
import useWithdrawFormData from '../../hooks/use-withdraw-form-data'
import Header from '../shared/header/header.view'
import InfoIcon from '@mui/icons-material/Info';
import DownArrow from '@mui/icons-material/ArrowDownward';
import { NETWORKS } from '../../chains';
import {SPECULAR_NETWORK_ID, CHIADO_NETWORK_ID} from '../../constants';
import { TOKEN } from '../../tokens';

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
  onSubmit: (amount: BigNumber, selectedTokenKey: number) => void,
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
  const [selectedTokenKey, setSelectedTokenKey] = useState(1);
  const handleChange = (event:any) => {
    setSelectedTokenKey(event.target.value);
  };
  const selectedToken = TOKEN[selectedTokenKey];
  const { values, amounts, l1balance, l2balance, error, changeWithdrawValue } = useWithdrawFormData(wallet, selectedTokenKey, l1Provider, l2Provider)
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
        title={`Specular Bridge`}
        onDisconnectWallet={onDisconnectWallet}
      />
      <form
        className={classes.form}
        onSubmit={(event) => {
          event.preventDefault()
          onSubmit(amounts.from, selectedTokenKey)
        }}
      >
        <div className={classes.card}>
          <Select
          className={classes.cardTitleText}
          value={selectedTokenKey || ''}
          onChange={handleChange}
          sx={{ boxShadow: 'none', '.MuiOutlinedInput-notchedOutline': { border: 0 } }}
        >
          {Object.keys(TOKEN).map((key) => (
            <MenuItem key={key} value={key}>
              {TOKEN[key].l2TokenName}
            </MenuItem>
          ))}
        </Select>
          <input
            ref={inputEl}
            className={classes.fromInput}
            placeholder='0.0'
            value={values.from}
            onChange={event => changeWithdrawValue(event.target.value)}
          />
          <p className={classes.toValue}>
            Balance: {formatUnits(l2balance, NETWORKS[SPECULAR_NETWORK_ID].nativeCurrency.decimals)} {selectedToken.l2TokenSymbol}
          </p>
        </div>
        <DownArrow className={classes.cardIcon} />
        <div className={classes.card}>
          <p className={classes.cardTitleText}>
          {selectedToken.l1TokenName}
          </p>
          <p>
            {formatUnits(amounts.to, NETWORKS[CHIADO_NETWORK_ID].nativeCurrency.decimals)} {NETWORKS[CHIADO_NETWORK_ID].nativeCurrency.symbol}
          </p>
          <p className={classes.toValue}>
            Balance: {formatUnits(l1balance, NETWORKS[CHIADO_NETWORK_ID].nativeCurrency.decimals)} {selectedToken.l1TokenSymbol}
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
