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

const useHeaderStyles = createUseStyles((theme: Theme) => ({
  
  header: {
    width: '100%',
    position: 'relative',
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    justifyContent: 'center',
    marginTop: theme.spacing(4)
  },
  title: {
    fontSize: theme.spacing(3),
    fontWeight: theme.fontWeights.bold
  },
  address: {
    marginTop: theme.spacing(2),
    fontSize: theme.spacing(2.5),
    color: theme.palette.grey.dark
  },
  balance: {
    marginTop: theme.spacing(1),
    fontSize: theme.spacing(2),
    color: theme.palette.black
  },
  goBackButton: {
    position: 'absolute',
    background: 'transparent',
    appearance: 'none',
    border: 'none',
    cursor: 'pointer',
    padding: theme.spacing(1),
    left: -(theme.spacing(1)),
    top: -(theme.spacing(0.5)),
    '&:disabled': {
      cursor: 'default'
    },
  },
  closeButton: {
    position: 'absolute',
    background: 'transparent',
    appearance: 'none',
    border: 'none',
    cursor: 'pointer',
    padding: theme.spacing(1),
    right: -(theme.spacing(1)),
    top: -(theme.spacing(0.5)),
  },
  disconnectButton: {
    background: 'transparent',
    border: 'none',
    cursor: 'pointer',
    textDecoration: 'underline',
    fontSize: theme.spacing(2),
    marginTop: theme.spacing(1.5)
  },
  closeIcon: {
    width: 16
  }
}))

export default useHeaderStyles
