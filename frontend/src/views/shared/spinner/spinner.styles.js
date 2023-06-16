import { createUseStyles } from 'react-jss'

const useSpinnerStyles = createUseStyles(theme => ({
  '@keyframes spin': {
    from: { transform: 'rotate(0deg)' },
    to: { transform: 'rotate(360deg)' }
  },
  spinner: ({ size }) => ({
    width: size,
    height: size,
    overflow: 'hidden'
  }),
  svg: {
    animation: '$spin 0.8s linear infinite'
  },
  topCircle: {
    stroke: theme.palette.purple,
    strokeLinecap: 'round',
    strokeDasharray: '30px 200px',
    strokeDashoffset: '0px'
  },
  bottomCircle: {
    stroke: theme.palette.purple,
    strokeOpacity: 0.2
  }
}))

export default useSpinnerStyles
