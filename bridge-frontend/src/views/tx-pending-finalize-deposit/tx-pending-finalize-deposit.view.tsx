import { useEffect } from 'react';
import Header from '../shared/header/header.view';
import useTxPendingStyles from './tx-pending-finalize-deposit.styles';
import LinkIcon from '@mui/icons-material/OpenInNew';
import Spinner from '../shared/spinner/spinner.view';
import { NETWORKS } from '../../chains';
import { ethers } from 'ethers'
import {
  CHIADO_NETWORK_ID,
  SPECULAR_NETWORK_ID,
  L1PORTAL_ADDRESS,
  CHIADO_RPC_URL,
  L1ORACLE_ADDRESS,
  SPECULAR_RPC_URL,
} from "../../constants";
import {
  IL1Portal__factory,
  IL2Portal__factory,
  L1Oracle__factory,
  IRollup__factory,
  ISequencerInbox__factory,
} from "../../typechain-types";
import type { PendingDeposit } from "../../types";

interface TxPendingProps {
 wallet: {
    address: string;
    chainId: number;
    provider: any;
  };
  depositData: {
    status: string;
    data?: {
      hash: string;
    };
  };
  setPendingDeposit: (arg0:PendingDeposit) => void;
  onGoToFinalizeStep: () => void;
}

function TxPendingFinalizeDeposit({ wallet, depositData, setPendingDeposit, onGoToFinalizeStep }: TxPendingProps) {
  const classes = useTxPendingStyles();

  const l1Provider =  wallet.provider
  const l1Portal = IL1Portal__factory.connect(
    L1PORTAL_ADDRESS,
    l1Provider
  );

  useEffect(() => {
    l1Portal.on(
      l1Portal.filters.DepositInitiated(),
      (nonce, sender, target, value, gasLimit, data, depositHash, event) => {
        if (event.transactionHash === depositData.data?.hash) {
          const newPendingDeposit: PendingDeposit = {
            l1BlockNumber: event.blockNumber,
            proofL1BlockNumber: undefined,
            depositHash: depositHash,
            depositTx: {
              nonce,
              sender,
              target,
              value,
              gasLimit,
              data,
            },
          }
         setPendingDeposit(newPendingDeposit);
         onGoToFinalizeStep();
        }
      }
    )
  }, [l1Portal, depositData, setPendingDeposit, onGoToFinalizeStep]);

  return (
    <div className={classes.txOverview}>
      <Header address={wallet.address} title={`xDAI → ETH`} />
      <div className={classes.spinnerWrapper}>
        <Spinner/>
      </div>
      <div className={classes.buttonGroup}>
        <a
          className={classes.button}
          href={`${NETWORKS[wallet.chainId].blockExplorerUrl}/tx/${depositData?.data?.hash}`}
          target='_blank'
          rel='noopener noreferrer'
        >
          Check transaction status here
          <LinkIcon className={classes.buttonIcon} />
        </a>
      </div>
    </div>
  );
}

export default TxPendingFinalizeDeposit;
