import { createUseStyles } from 'react-jss'

const useLoginStyles = createUseStyles((theme) => ({
  login: {
    flex: 1,
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    paddingTop: theme.spacing(7),
    paddingBottom: theme.spacing(4),
    justifyContent: 'space-between'
  },
  tokenLogos: {
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center'
  },
  swapArrow: {
    margin: `0 ${theme.spacing(2)}px`
  },
  title: {
    fontSize: theme.spacing(3),
    fontWeight: theme.fontWeights.bold
  },
  column: {
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center'
  },
  metaMaskButton: {
    padding: theme.spacing(2.5),
    borderRadius: theme.spacing(3.5),
    background: theme.palette.white,
    boxShadow: '0px 3.75px 17px #b3b3b3',
    marginBottom: theme.spacing(3),
    border: 'none',
    transition: theme.buttonTransition,
    cursor: 'pointer'
  },
  connectText: {
    fontSize: theme.spacing(2),
  },
  metaMaskIcon: {
    width: theme.spacing(7),
    height: theme.spacing(7)
  },
  metaMaskNameText: {
    fontWeight: theme.fontWeights.bold,
    marginTop: theme.spacing(2)
  },
  footer: {
    color: '#000',
    fontSize: '16px',
    textDecoration: 'underline'
  }
}))

export default useLoginStyles
