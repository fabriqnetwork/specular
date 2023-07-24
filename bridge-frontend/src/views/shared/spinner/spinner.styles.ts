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

const useSpinnerStyles = createUseStyles((theme: Theme) => ({
  '@keyframes spin': {
    from: { transform: 'rotate(0deg)' },
    to: { transform: 'rotate(360deg)' }
  },
  spinner:{
    width: 64,
    height: 64,
    overflow: 'hidden',
  },
  svg: {
    animation: '$spin 0.8s linear infinite'
  },
  topCircle: {
    stroke: theme.palette.purple,
    strokeLinecap: 'round',
    strokeDasharray: '30px 200px',
    strokeDashoffset: '0px'
  },
  bottomCircle: {
    stroke: theme.palette.purple,
    strokeOpacity: 0.2
  }
}))

export default useSpinnerStyles
