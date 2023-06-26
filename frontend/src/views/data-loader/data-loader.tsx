import React, { useEffect } from 'react';

import Spinner from '../shared/spinner/spinner.view';
import useDataLoaderStyles from './data-loader.styles';

interface DataLoaderProps {
  onGoToNextStep: () => void;
}

const DataLoader: React.FC<DataLoaderProps> = ({ onGoToNextStep }) => {
  const classes = useDataLoaderStyles();

  useEffect(() => {
    onGoToNextStep();
  }, [onGoToNextStep]);

  return (
    <div className={classes.dataLoader}>
      <Spinner />
    </div>
  );
};

export default DataLoader;
