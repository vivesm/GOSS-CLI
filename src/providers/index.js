import { OpenAICompatibleProvider } from './openai-compatible.js';
import { OpenAIProvider } from './openai.js';
import { OllamaProvider } from './ollama.js';

export function createProvider(config) {
  const provider = config.provider || detectProvider(config);
  
  switch (provider) {
    case 'openai':
      return new OpenAIProvider(config);
    case 'ollama':
      return new OllamaProvider(config);
    case 'lmstudio':
    case 'localai':
    case 'openai-compatible':
    default:
      return new OpenAICompatibleProvider(config);
  }
}

function detectProvider(config) {
  const apiBase = config.apiBase?.toLowerCase() || '';
  
  if (apiBase.includes('openai.com')) {
    return 'openai';
  }
  if (apiBase.includes('11434')) {
    return 'ollama';
  }
  if (apiBase.includes('1234')) {
    return 'lmstudio';
  }
  if (apiBase.includes('8080')) {
    return 'localai';
  }
  
  return 'openai-compatible';
}