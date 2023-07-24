import { createUseStyles } from 'react-jss';

const useStepperStyles = createUseStyles((theme: any) => ({
  stepper: {
    minWidth: 525,
    minHeight: 562.5,
    background: theme.palette.white,
    borderRadius: theme.spacing(3.75),
    boxShadow: '0px 0px 4px rgba(200, 200, 221, 0.53)',
    padding: `0 ${theme.spacing(3.5)}px`,
    paddingBottom: theme.spacing(3),
    display: 'flex',
    flexDirection: 'column',
    zIndex: 2,
    position: 'relative'
  },
  container: {
    position: 'relative',
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    padding: '80px 0 20px'
  },
  tabs: {
    display: 'flex',
    justifyContent: 'center',
    position: 'absolute',
    top: 40,
    left: '50%',
    transform: 'translateX(-50%)',
    zIndex: 1
  },
  tab: {
    minWidth: '235px',
    height: '40px',
    padding: '0 20px',
    background: theme.palette.secondary,
    borderTopLeftRadius: '10px',
    borderTopRightRadius: '10px',
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    boxShadow: '0px 0px 4px rgba(200, 200, 221, 0.53)',
    marginRight: '1px',
    cursor: 'pointer',
    border: 0,
    '&:disabled': {
      cursor: 'default'
    }
  },
  tabActive: {
    composes: '$tab',
    background: theme.palette.white
  },
  tabName: {
    whiteSpace: 'nowrap'
  }
}));

export default useStepperStyles;
