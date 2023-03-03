import {
  IL1Portal__factory,
  IL2Portal__factory,
} from "../../../contracts/typechain-types";
import type { IL1Portal, IL2Portal } from "../../../contracts/typechain-types";
import type { Provider } from "@ethersproject/providers";

export function getL1PortalContract(
  address: string,
  provider: Provider
): IL1Portal {
  return IL1Portal__factory.connect(address, provider);
}

export function getL2PortalContract(
  address: string,
  provider: Provider
): IL2Portal {
  return IL2Portal__factory.connect(address, provider);
}
