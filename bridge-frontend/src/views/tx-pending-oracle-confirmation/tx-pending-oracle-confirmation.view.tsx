import { useEffect } from 'react';
import Header from '../shared/header/header.view';
import useTxPendingStyles from './tx-pending-oracle-confirmation.styles';
import LinkIcon from '@mui/icons-material/OpenInNew';
import Spinner from '../shared/spinner/spinner.view';
import { NETWORKS } from '../../chains';
import * as React from 'react';
import { ethers } from 'ethers'
import {
  CHIADO_NETWORK_ID,
  L1ORACLE_ADDRESS,
  SPECULAR_RPC_URL,
} from "../../constants";
import {
  L1Oracle__factory
} from "../../typechain-types";
import type { PendingDeposit } from "../../types";

interface PendingData {
  status: string;
  data: PendingDeposit;
}
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
  pendingDeposit: PendingData;
  setPendingDeposit: (args1:any) => void;
  onGoToNextStep: () => void;
  switchChain: (args1:any) => void;

}

function TxPendingOracleConfirmation({ wallet, depositData,pendingDeposit, setPendingDeposit,switchChain, onGoToNextStep }: TxPendingProps) {
  const classes = useTxPendingStyles();
  const l2Provider = new ethers.providers.StaticJsonRpcProvider(SPECULAR_RPC_URL);
  // switchChain(SPECULAR_NETWORK_ID.toString())

  useEffect(() => {
    const l1Oracle = L1Oracle__factory.connect(
      L1ORACLE_ADDRESS,
      l2Provider
    );

    l1Oracle.on(
      l1Oracle.filters.L1OracleValuesUpdated(),
      (blockNumber, stateRoot, event) => {
        console.log("Main Oracle Blocknumber is "+blockNumber)
        if (pendingDeposit.data === undefined) {
          return;
        }
        if (blockNumber.gte(pendingDeposit.data.l1BlockNumber) && pendingDeposit.status === 'initiated') {
          pendingDeposit.data.proofL1BlockNumber = blockNumber.toNumber();
          console.log("Main Oracle Correct Blocknumber is "+blockNumber);
          setPendingDeposit(pendingDeposit);
          onGoToNextStep();
        }
      }
    );

  },[pendingDeposit]
  )
  return (
    <div className={classes.txOverview}>
      <Header address={wallet.address} title={`Specular Bridge`} />
      <div className={classes.spinnerWrapper}>
        <Spinner/>
      </div>
      <div className={classes.buttonGroup}>
        <a
          className={classes.button}
          href={`${NETWORKS[CHIADO_NETWORK_ID].blockExplorerUrl}/tx/${depositData?.data?.hash}`}
          target='_blank'
          rel='noopener noreferrer'
        >
          Transaction on {NETWORKS[CHIADO_NETWORK_ID].name} is successful. Check transaction details here
          <LinkIcon className={classes.buttonIcon} />
        </a>
      </div>
      Waiting for confirmation from Oracle.
    </div>
  );
}

export default TxPendingOracleConfirmation;
