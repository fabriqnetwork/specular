import { formatUnits } from 'ethers/lib/utils';
import ArrowLeft from '@mui/icons-material/ArrowLeft';
import CrossIcon from '@mui/icons-material/Clear';
import useHeaderStyles from './header.styles';

interface HeaderProps {
  address?: string;
  title: string;
  isGoBackButtonDisabled?: boolean;
  onGoBack?: () => void;
  onDisconnectWallet?: () => void;
  onClose?: () => void;
  tokenInfo?: {
    decimals: number;
    symbol: string;
  };
  balance?: number;
}

function Header({
  address,
  title,
  isGoBackButtonDisabled,
  onGoBack,
  onDisconnectWallet,
  onClose,
  tokenInfo,
  balance,
}: HeaderProps) {
  const classes = useHeaderStyles();

  function getPartiallyHiddenEthereumAddress(ethereumAddress: string) {
    const firstAddressSlice = ethereumAddress.slice(0, 6);
    const secondAddressSlice = ethereumAddress.slice(
      ethereumAddress.length - 4,
      ethereumAddress.length
    );

    return `${firstAddressSlice} *** ${secondAddressSlice}`;
  }

  return (
    <div className={classes.header}>
      {onGoBack && (
        <button
          disabled={isGoBackButtonDisabled}
          className={classes.goBackButton}
          onClick={onGoBack}
        >
          <ArrowLeft />
        </button>
      )}
      {onClose && (
        <button className={classes.closeButton} onClick={onClose}>
          <CrossIcon className={classes.closeIcon} />
        </button>
      )}
      <p className={classes.title}>{title}</p>
      {address && (
        <p className={classes.address}>{getPartiallyHiddenEthereumAddress(address)}</p>
      )}
      {balance && tokenInfo && (
        <p className={classes.balance}>
          Balance: {Number(formatUnits(balance, tokenInfo.decimals))} {tokenInfo.symbol}
        </p>
      )}
      {onDisconnectWallet && (
        <button className={classes.disconnectButton} onClick={onDisconnectWallet}>
          Disconnect
        </button>
      )}
    </div>
  );
}

export default Header;
