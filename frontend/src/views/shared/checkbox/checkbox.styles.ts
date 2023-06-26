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

const useStyles = createUseStyles((theme : Theme) => ({
  checkIcon: {
    '& > path': {
      stroke: theme.palette.black
    }
  },
  row: {
    display: 'flex',
    alignItems: 'center',
    cursor: 'pointer',
    marginBottom: 30
  },
  checkbox: {
    width: 18,
    height: 18,
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    borderRadius: 4,
    marginRight: 16,
    border: `1px solid ${theme.palette.black}`
  },
  label: {
    flex: 1,
    lineHeight: '26px'
  }
}))

export default useStyles
