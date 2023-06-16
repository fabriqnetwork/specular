import { useState, useEffect } from 'react';
import {
  IL1Portal__factory,
  IL2Portal__factory,
} from '../../../contracts/typechain-types';

function usePortalContracts(address, provider) {
  const [l1PortalContract, setL1PortalContract] = useState(null);
  const [l2PortalContract, setL2PortalContract]= useState(null);

  useEffect(() => {
    if (!address || !provider) {
      setL1PortalContract(null);
      setL2PortalContract(null);
      return;
    }

    const l1Portal = IL1Portal__factory.connect(address, provider);
    const l2Portal = IL2Portal__factory.connect(address, provider);

    setL1PortalContract(l1Portal);
    setL2PortalContract(l2Portal);
  }, [address, provider]);

  return [l1PortalContract, l2PortalContract];
}

export default usePortalContracts;