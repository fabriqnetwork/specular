
import { createUseStyles } from 'react-jss'

const useAppStyles = createUseStyles((theme) => ({
  container: {
    marginLeft:'0',
    marginRight:'0',
    verticalAlign:'middle', 
    margin:'middle'
  },
  '@font-face': [
    {
      fontFamily: 'Modern Era',
      src: "url('./fonts/modern-era/ModernEra-Regular.woff2') format('woff2')",
      fallbacks: [
        { src: "url('./fonts/modern-era/ModernEra-Regular.woff') format('woff')" },
        { src: "url('./fonts/modern-era/ModernEra-Regular.ttf') format('truetype')" }
      ],
      fontWeight: 400,
      fontStyle: 'normal'
    },
    {
      fontFamily: 'Modern Era',
      src: "url('./fonts/modern-era/ModernEra-Medium.woff2') format('woff2')",
      fallbacks: [
        { src: "url('./fonts/modern-era/ModernEra-Medium.woff') format('woff')" },
        { src: "url('./fonts/modern-era/ModernEra-Medium.ttf') format('truetype')" }
      ],
      fontWeight: 500,
      fontStyle: 'normal'
    },
    {
      fontFamily: 'Modern Era',
      src: "url('./fonts/modern-era/ModernEra-Bold.woff2') format('woff2')",
      fallbacks: [
        { src: "url('./fonts/modern-era/ModernEra-Bold.woff') format('woff')" },
        { src: "url('./fonts/modern-era/ModernEra-Bold.ttf') format('truetype')" }
      ],
      fontWeight: 700,
      fontStyle: 'normal'
    }
  ],
  '@global': {
    '*': {
      fontFamily: 'Modern Era',
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
}))

export default useAppStyles
