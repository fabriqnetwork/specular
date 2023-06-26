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
  };
  buttonTransition: string;
}
const useFooterStyles = createUseStyles((theme: Theme) => ({
  footer: {
    width: '100%',
    position: 'relative',
    display: 'flex',
    flexDirection: 'column',
    flexWrap: 'wrap',
    alignItems: 'center',
    justifyContent: 'center',
    marginTop: theme.spacing(4),
    padding: '15px 0 15px',
    backgroundColor: "black"
  },
  logo:{
    maxWidth: 100,
    maxHeight: 100
    },
  bottom: {
    padding: '20px 0 15px',
  },
}))

export default useFooterStyles
