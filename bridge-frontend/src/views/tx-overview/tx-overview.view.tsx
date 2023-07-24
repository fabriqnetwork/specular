import Header from '../shared/header/header.view';
import useTxOverviewStyles from './tx-overview.styles';
import { ReactComponent as CheckIcon } from '../../images/check-icon.svg';
import LinkIcon from '@mui/icons-material/OpenInNew';
import { ReactComponent as MetaMaskLogo } from '../../images/metamask-logo.svg';
import useWatchAsset from '../../hooks/use-watch-asset';
import { NETWORKS } from '../../chains';
import {
  CHIADO_NETWORK_ID,
  SPECULAR_NETWORK_ID
} from "../../constants";

interface TxOverviewProps {
  wallet: {
    address: string;
    chainId: number;
    provider: any;
  };
  transactionData: {
    status: string;
    data?: {
      hash: string;
    };
  };
  finalizeTransactionData: {
    status: string;
    data?: string;
  };
  onDisconnectWallet: () => void;
  isMetamask: boolean;
}

function TxOverview({
  wallet,
  transactionData,
  finalizeTransactionData,
  onDisconnectWallet,
  isMetamask,
}: TxOverviewProps) {
  const classes = useTxOverviewStyles();
  const watchAsset = useWatchAsset();

  return (
    <div className={classes.txOverview}>
      <Header
        address={wallet.address}
        title={`xDAI → ETH`}
        onDisconnectWallet={onDisconnectWallet}
      />
      <CheckIcon className={classes.checkIcon} />
      <div className={classes.buttonGroup}>
        <a
          className={classes.button}
          href={`${NETWORKS[CHIADO_NETWORK_ID].blockExplorerUrl}/tx/${transactionData?.data?.hash}`}
          target='_blank'
          rel='noopener noreferrer'
          >
          Check transaction details here
          <LinkIcon className={classes.buttonIcon} />
        </a>
        {isMetamask && (
          <button className={classes.button}>
            Add xDAI token to MetaMask
            <MetaMaskLogo className={classes.buttonIcon} />
          </button>
        )}
      </div>
      <div className={classes.buttonGroup}>
        <a
          className={classes.button}
          href={`${NETWORKS[SPECULAR_NETWORK_ID].blockExplorerUrl}/tx/${finalizeTransactionData?.data}`}
          target='_blank'
          rel='noopener noreferrer'
        >
          Check transaction details here
          <LinkIcon className={classes.buttonIcon} />
        </a>
        {isMetamask && (
          <button className={classes.button}>
            Add xDAI token to MetaMask
            <MetaMaskLogo className={classes.buttonIcon} />
          </button>
        )}
      </div>
    </div>
  );
}

export default TxOverview;
