import { createUseStyles } from 'react-jss'

const useNetworkErrorStyles = createUseStyles((theme) => ({
  networkError: {
    flex: 1,
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    color: theme.palette.black,
    paddingTop: theme.spacing(15.5)
  },
  title: {
    marginTop: theme.spacing(3),
    fontWeight: theme.fontWeights.bold,
    fontSize: theme.spacing(2.5)
  },
  description: {
    fontWeight: theme.fontWeights.medium,
    fontSize: theme.spacing(2),
    textAlign: 'center'
  },
  descriptionContainer: {
    display: 'flex',
    alignItems: 'center',
    marginTop: theme.spacing(2),
  },
  networkName: {
    fontWeight: theme.fontWeights.bold,
    fontSize: theme.spacing(2),
    marginLeft: '4px',
  },
  switchNetworkButton: {
    fontWeight: theme.fontWeights.medium,
    fontSize: theme.spacing(2),
    display: 'inline-flex',
    alignItems: 'center',
    padding: '3px 7px',
    position: 'relative',
    cursor: 'pointer',
    backgroundColor: '#f3e2fd',
    borderRadius: '4px',
    margin: '0 7px',
    hover: {
      backgroundColor: '#f2dbff',
    }
  }
}))

export default useNetworkErrorStyles
