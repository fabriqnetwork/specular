import React from 'react'
import PropTypes from 'prop-types'

import useSpinnerStyles from './spinner.styles'

const SIZE = 64
const THICKNESS = 8

function Spinner () {
  const classes = useSpinnerStyles({ size: SIZE })

  return (
    <div className={classes.spinner}>
      <svg
        className={classes.svg}
        viewBox={`${SIZE / 2} ${SIZE / 2} ${SIZE} ${SIZE}`}
      >
        <circle
          className={classes.bottomCircle}
          cx={SIZE}
          cy={SIZE}
          r={(SIZE - THICKNESS) / 2}
          fill='none'
          strokeWidth={THICKNESS}
        />
        <circle
          className={classes.topCircle}
          cx={SIZE}
          cy={SIZE}
          r={(SIZE - THICKNESS) / 2}
          fill='none'
          strokeWidth={THICKNESS}
        />
      </svg>
    </div>
  )
}

Spinner.propTypes = {
  size: PropTypes.number
}

export default Spinner
