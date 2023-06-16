import { createUseStyles } from 'react-jss'

const useFooterStyles = createUseStyles((theme) => ({
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
