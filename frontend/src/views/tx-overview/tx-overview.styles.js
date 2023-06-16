import { createUseStyles } from 'react-jss'

const useTxOverviewStyles = createUseStyles((theme) => ({
  txOverview: {
    flex: 1,
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center'
  },
  checkIcon: {
    marginLeft: theme.spacing(3),
    marginTop: theme.spacing(4)
  },
  title: {
    fontSize: theme.spacing(3),
    fontWeight: theme.fontWeights.bold,
    marginTop: theme.spacing(4)
  },
  buttonGroup: {
    marginTop: theme.spacing(2),
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center'
  },
  button: {
    fontSize: theme.spacing(2.5),
    color: theme.palette.grey.dark,
    appearance: 'none',
    border: 'none',
    background: 'transparent',
    cursor: 'pointer',
    display: 'flex',
    alignItems: 'center',
    padding: theme.spacing(1),
    marginTop: theme.spacing(2),
    '&:first-child': {
      marginTop: theme.spacing(0)
    }
  },
  buttonIcon: {
    marginLeft: theme.spacing(1)
  },
  howToUseLink: {
    color: '#000',
    fontSize: theme.spacing(2.5),
    textDecoration: 'underline',
    padding: theme.spacing(1)
  },
  note: {
    backgroundColor: '#fff3d6',
    padding: '10px 20px',
    fontSize: theme.spacing(2),
    lineHeight: '18px',
    marginTop: theme.spacing(2)
  },
  noteTitle: {
    color: '#ffb800',
    fontSize: theme.spacing(2.5),
    marginBottom: theme.spacing(0.5)
  },
  noteLink: {
    color: '#7280f7'
  }
}))

export default useTxOverviewStyles
