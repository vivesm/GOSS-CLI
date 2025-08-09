import { createProvider } from '../providers/index.js';

export function createClient(cfg) {
  const provider = createProvider(cfg);
  
  return {
    chatComplete: (params) => provider.chatComplete(params),
    listModels: () => provider.listModels(),
    getProviderName: () => provider.getName()
  };
}