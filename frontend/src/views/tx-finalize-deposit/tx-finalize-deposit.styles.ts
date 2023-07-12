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

const useTxOverviewStyles = createUseStyles((theme: Theme) => ({
  txOverview: {
    flex: 1,
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center'
  },
  checkIcon: {
    marginLeft: theme.spacing(3),
    marginTop: theme.spacing(3)
  },
  title: {
    fontSize: theme.spacing(3),
    fontWeight: theme.fontWeights.bold,
    marginTop: theme.spacing(8)
  },
  buttonGroup: {
    marginTop: theme.spacing(2),
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center'
  },
  button: {
    fontSize: theme.spacing(2.5),
    color: theme.palette.grey.dark,
    appearance: 'none',
    border: 'none',
    background: 'transparent',
    cursor: 'pointer',
    display: 'flex',
    alignItems: 'center',
    padding: theme.spacing(1),
    '&:not(:first-child)': {
      marginTop: theme.spacing(2)
    }
  },
  buttonIcon: {
    marginLeft: theme.spacing(1)
  },
  howToUseLink: {
    color: '#000',
    fontSize: theme.spacing(2.5),
    textDecoration: 'underline',
    padding: theme.spacing(1)
  },
  spinnerWrapper: {
    marginTop: theme.spacing(14.5)
  },
}))

export default useTxOverviewStyles
