import { createUseStyles } from 'react-jss'


interface Theme {
  spacing: {
    (value: number): string;
    unit: number;
  };
  fontWeights: {
    regular: number;
    bold: number;
    normal: number;
    medium: number;
    extraBold: number;
  };

  palette: {
    grey: {
      light1: string;
      light2: string;
      dark: string;
      main: string;
    };
    red: string;
    purple: string;
    primary: string;
    white: string;
    black: string;
  };
  buttonTransition: string;
}

const useTxConfirmChainStyles = createUseStyles((theme: Theme) => ({
  txConfirmChain: {
    flex: 1,
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center'
  },
  spinnerWrapper: {
    marginTop: theme.spacing(14.5)
  },
  title: {
    fontSize: theme.spacing(3),
    fontWeight: theme.fontWeights.bold,
    marginTop: theme.spacing(8)
  }
}))

export default useTxConfirmChainStyles
