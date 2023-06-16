import { createUseStyles } from 'react-jss'

const useHeaderStyles = createUseStyles((theme) => ({
  header: {
    width: '100%',
    position: 'relative',
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    justifyContent: 'center',
    marginTop: theme.spacing(-1),
    padding: '25px 0 0px'
  },
  title: {
    fontSize: theme.spacing(5),
    fontWeight: theme.fontWeights.bold,
    float: 'left'
  },
  titleText: {
    alignItems: 'bottom',
    justifyContent: 'bottom',
    padding: '20px 7px 0'
    
  },
  logo:{
    maxWidth: 500,
    maxHeight: 500
    }
}))

export default useHeaderStyles
