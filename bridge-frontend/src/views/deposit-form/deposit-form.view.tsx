import { useEffect, useRef, useState } from 'react'
import { BigNumber } from 'ethers'
import { formatUnits } from 'ethers/lib/utils'
import {Select, MenuItem} from "@mui/material";
import useDepositFormStyles from './deposit-form.styles'
import useDepositFormData from '../../hooks/use-deposit-form-data'
import Header from '../shared/header/header.view'
import InfoIcon from '@mui/icons-material/Info';
import DownArrow from '@mui/icons-material/ArrowDownward';
import { NETWORKS } from '../../chains';
import {SPECULAR_NETWORK_ID, CHIADO_NETWORK_ID} from '../../constants';
import { TOKEN } from '../../tokens';

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
  onSubmit: (amount: BigNumber, selectedTokenKey: number) => void,
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

  const [selectedTokenKey, setSelectedTokenKey] = useState(1);
  const handleChange = (event:any) => {
    setSelectedTokenKey(event.target.value);
  };
  const selectedToken = TOKEN[selectedTokenKey];
  var { values, amounts, l1balance, l2balance, error, changeDepositValue } = useDepositFormData(wallet, selectedTokenKey, l1Provider, l2Provider);

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
        title={`Specular Bridge`}
        onDisconnectWallet={onDisconnectWallet}
      />
      <form
        className={classes.form}
        onSubmit={(event) => {
          event.preventDefault()
          onSubmit(amounts.from,selectedTokenKey)
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
              {TOKEN[key].l1TokenName}
            </MenuItem>
          ))}
        </Select>
          <input
            ref={inputEl}
            className={classes.fromInput}
            placeholder='0.0'
            value={values.from}
            onChange={event => changeDepositValue(event.target.value)}
          />
          <p className={classes.toValue}>
            Balance: {formatUnits(l1balance, NETWORKS[CHIADO_NETWORK_ID].nativeCurrency.decimals)} {selectedToken.l1TokenSymbol}
          </p>
        </div>
        <DownArrow className={classes.cardIcon} />
        <div className={classes.card}>
          <p className={classes.cardTitleText}>
            {selectedToken.l2TokenName}
          </p>
          <p>
            {formatUnits(amounts.to, NETWORKS[SPECULAR_NETWORK_ID].nativeCurrency.decimals)} {NETWORKS[SPECULAR_NETWORK_ID].nativeCurrency.symbol}
          </p>
          <p className={classes.toValue}>
            Balance: {formatUnits(l2balance, NETWORKS[SPECULAR_NETWORK_ID].nativeCurrency.decimals)} {selectedToken.l2TokenSymbol}
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
