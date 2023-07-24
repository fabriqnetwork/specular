import { formatUnits } from 'ethers/lib/utils';
import useHeaderStyles from './header.styles';
import ArrowLeft from '@mui/icons-material/ArrowLeft';
import CrossIcon from '@mui/icons-material/Clear';

interface HeaderProps {
  address?: string;
  title: string;
  onDisconnectWallet?: () => void;
}

function Header({ address, title, onDisconnectWallet }: HeaderProps): JSX.Element {
  const classes = useHeaderStyles();

  function getPartiallyHiddenEthereumAddress(ethereumAddress: string): string {
    const firstAddressSlice = ethereumAddress.slice(0, 6);
    const secondAddressSlice = ethereumAddress.slice(
      ethereumAddress.length - 4,
      ethereumAddress.length
    );

    return `${firstAddressSlice} *** ${secondAddressSlice}`;
  }

  return (
    <div className={classes.header}>
      <p className={classes.title}>{title}</p>
      {address && (
        <p className={classes.address}>{getPartiallyHiddenEthereumAddress(address)}</p>
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
