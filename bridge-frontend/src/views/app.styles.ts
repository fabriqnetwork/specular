import { createUseStyles } from 'react-jss';

const useAppStyles = createUseStyles((theme: any) => ({
  container: {
    marginLeft: '0',
    marginRight: '0',
    verticalAlign: 'middle',
    margin: 'middle'
  },
  '@font-face': [
    {
      fontWeight: 400,
      fontStyle: 'normal'
    },
    {
      fontWeight: 500,
      fontStyle: 'normal'
    },
    {
      fontWeight: 700,
      fontStyle: 'normal'
    }
  ],
  '@global': {
    '*': {
      boxSizing: 'border-box'
    },
    'body, input, button': {
      fontSize: theme.spacing(2.5)
    },
    body: {
      margin: 0,
      display: 'flex',
      minHeight: '100vh',
      color: theme.palette.black,
      background: theme.palette.secondary,
    },
    '#root': {
      flex: 1,
      alignItems: 'center',
      justifyContent: 'center'
    },
    a: {
      textDecoration: 'none',
      color: 'inherit'
    },
    'p, h1, h2, h3, h4, h5, h6': {
      margin: 0
    }
  }
}));

export default useAppStyles;
