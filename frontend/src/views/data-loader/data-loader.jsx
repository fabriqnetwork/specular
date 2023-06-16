import { useEffect } from 'react'

import Spinner from '../shared/spinner/spinner.view'
import useDataLoaderStyles from './data-loader.styles'

function DataLoader ({onFinishLoading }) {
  const classes = useDataLoaderStyles()

  

  return (
    <div className={classes.dataLoader}>
      <Spinner />
    </div>
  )
}

export default DataLoader
