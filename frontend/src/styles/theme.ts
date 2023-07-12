const theme = {
  palette: {
    white: '#ffffff',
    black: '#2b2b2b',
    primary: '#0000FF',
    secondary: '#f6f7fa',
    grey: {
      light1: '#f5f5f5', // Light Grey
      light2: '#e0e0e0', // Silver
      main: '#9e9e9e',   // Gray
      dark: '#424242'    // Dim Gray
    },
    orange: '#006cff',
    purple: '#00008B',
    red: '#C41E3A'
  },
  hoverTransition: 'all 100ms',
  fontWeights: {
    normal: '400',
    medium: '500',
    bold: '700',
    extraBold: '800'
  },
  breakpoints: {
    upSm: '@media (min-width: 576px)'
  },
  spacing: (value: number) => value * 8,
  buttonTransition: 'all 100ms'
};

export default theme;
