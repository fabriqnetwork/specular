import useStyles from './checkbox.styles';

import { ReactComponent as CheckIcon } from '../../../images/check-icon-small.svg';

interface CheckboxProps {
  onClick: () => void;
  checked: boolean;
  children: React.ReactNode;
}

function Checkbox({ onClick, checked, children }: CheckboxProps) {
  const classes = useStyles();
  return (
    <div className={classes.row} onClick={onClick}>
      <div className={classes.checkbox}>
        {checked && <CheckIcon className={classes.checkIcon} />}
      </div>
      <span className={classes.label}>{children}</span>
    </div>
  );
}

export default Checkbox;
