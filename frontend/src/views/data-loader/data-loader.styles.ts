import { createUseStyles } from 'react-jss'

const useDataLoaderStyles = createUseStyles((theme) => ({
  dataLoader: {
    flex: 1,
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    justifyContent: 'center'
  }
}))

export default useDataLoaderStyles
