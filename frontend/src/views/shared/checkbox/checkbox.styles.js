import { createUseStyles } from 'react-jss'

const useStyles = createUseStyles((theme) => ({
  checkIcon: {
    '& > path': {
      stroke: theme.palette.black
    }
  },
  row: {
    display: 'flex',
    alignItems: 'center',
    cursor: 'pointer',
    marginBottom: 30
  },
  checkbox: {
    width: 18,
    height: 18,
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    borderRadius: 4,
    marginRight: 16,
    border: `1px solid ${theme.palette.black}`
  },
  label: {
    flex: 1,
    lineHeight: '26px'
  }
}))

export default useStyles
