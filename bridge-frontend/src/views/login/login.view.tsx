import { useEffect } from 'react';

import useLoginStyles from './login.styles';

interface LoginProps {
  wallet: any;
  onLoadWallet: () => void;
  onGoToNextStep: () => void;
}

function Login({ wallet, onLoadWallet, onGoToNextStep }: LoginProps) {
  const classes = useLoginStyles();

  useEffect(() => {
    if (wallet) {
      onGoToNextStep();
    }
  }, [wallet, onGoToNextStep]);

  return (
    <div className={classes.login}>
      <h1 className={classes.title}>Welcome to Specular Bridge</h1>
      <div className={classes.column}>
        <button
          className={classes.metaMaskButton}
          onClick={onLoadWallet}
        >
          Connect wallet
        </button>
        <span className={classes.connectText}>Connect your wallet to get started</span>
      </div>
      <div></div>
    </div>
  );
}

export default Login;
