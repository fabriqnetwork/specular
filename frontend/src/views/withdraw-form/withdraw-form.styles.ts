import { createUseStyles } from 'react-jss';

interface WithdrawFormStylesProps {
  error: boolean;
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


const useWithdrawFormStyles = createUseStyles((theme:Theme) => ({
  withdrawForm: {
    flex: 1,
  },
  form: {
    display: 'flex',
    flexDirection: 'column',
  },
  card: ({ error }: WithdrawFormStylesProps) => ({
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    marginTop: theme.spacing(3),
    paddingTop: theme.spacing(5),
    border: `2px solid ${error ? theme.palette.red : theme.palette.grey.light2}`,
    borderRadius: `${theme.spacing(2.5)}px`,
  }),
  cardTitleText: {
    paddingBottom:  theme.spacing(1),
    fontSize: theme.spacing(4),
    fontWeight: theme.fontWeights.bold,
  },
  fromInput: {
    fontSize: theme.spacing(5),
    fontWeight: theme.fontWeights.bold,
    appearance: 'none',
    border: 'none',
    outline: 'none',
    width: '100%',
    textAlign: 'center',
    marginTop: theme.spacing(2),
    caretColor: theme.palette.purple,
    padding: `0 ${theme.spacing(5)}px`,
    '&:disabled': {
      background: theme.palette.white,
    },
  },
  toValue: {
    fontSize: theme.spacing(2),
    color: theme.palette.grey.main,
    textAlign: 'center',
    width: '100%',
    marginTop: theme.spacing(4),
    borderTop: `2px solid ${theme.palette.grey.light2}`,
    padding: `${theme.spacing(2)}px ${theme.spacing(5)}px`,
  },
  inputErrorContainer: {
    fontSize: theme.spacing(2),
    fontWeight: theme.fontWeights.medium,
    display: 'flex',
    alignItems: 'center',
    color: theme.palette.red,
    marginTop: theme.spacing(2),
  },
  cardErrorIcon: {
    marginRight: theme.spacing(1),
    minWidth: theme.spacing(2),
    '& path': {
      fill: theme.palette.red,
    },
  },
  cardIcon: {
    marginTop: theme.spacing(4),
    alignSelf: 'center',
    maxWidth: theme.spacing(10),
    maxHeight: theme.spacing(10),
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
  withdrawError: {
    fontSize: theme.spacing(2),
    fontWeight: theme.fontWeights.medium,
    color: theme.palette.red,
    textAlign: 'center',
    marginTop: theme.spacing(4),
  },
}));

export default useWithdrawFormStyles;
