import { createUseStyles } from 'react-jss';

interface FinalizeDepositFormStylesProps {
}
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
  };
  buttonTransition: string;
}


const useFinalizeDepositFormStyles = createUseStyles((theme:Theme) => ({
  finalizeDepositForm: {
    flex: 1,
  },
  form: {
    display: 'flex',
    flexDirection: 'column',
  },
  submitButton: {
    fontSize: theme.spacing(2),
    fontWeight: theme.fontWeights.bold,
    margin: `${theme.spacing(12)}px auto ${theme.spacing(4)}px`,
    padding: `${theme.spacing(3)}px 0`,
    background: theme.palette.primary,
    color: theme.palette.white,
    width: '40%',
    borderRadius: theme.spacing(12.5),
    appearance: 'none',
    border: 'none',
    transition: theme.buttonTransition,
    cursor: 'pointer',
    '&:disabled': {
      background: theme.palette.grey.dark,
      cursor: 'default',
    },
  },
}));

export default useFinalizeDepositFormStyles;
