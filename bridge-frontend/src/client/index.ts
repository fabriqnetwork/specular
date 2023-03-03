import { configureChains, createClient, VagmiPlugin } from "vagmi";

import { chiadoChain, specularChain } from "./chains";
import { provider as defaultProvider } from "./providers";

const { chains, provider, webSocketProvider } = configureChains(
  [chiadoChain, specularChain],
  [defaultProvider]
);

console.log(chains);

const client = createClient({
  autoConnect: true,
  provider,
  webSocketProvider,
});

const plugin = VagmiPlugin(client);

export default plugin;
