import { OpenAICompatibleProvider } from './openai-compatible.js';
import { OpenAIProvider } from './openai.js';
import { OllamaProvider } from './ollama.js';

export function createProvider(config) {
  const provider = config.provider || detectProvider(config);
  
  // Validate conflicting configurations
  if (provider === 'openai' && !config.apiKey && !process.env.OPENAI_API_KEY) {
    console.warn('Warning: OpenAI provider specified but no API key found. Set OPENAI_API_KEY environment variable.');
  }
  
  if (config.debug) {
    console.error(`[DEBUG] Using provider: ${provider} (${config.provider ? 'explicit' : 'auto-detected'})`);
  }
  
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
  
  // More specific domain matching
  if (apiBase.includes('api.openai.com')) {
    return 'openai';
  }
  
  // Port-based detection (more specific)
  const url = new URL(apiBase);
  const port = url.port;
  
  if (port === '11434' || apiBase.includes('ollama')) {
    return 'ollama';
  }
  if (port === '1234' || apiBase.includes('lmstudio')) {
    return 'lmstudio';
  }
  if (port === '8080' && (apiBase.includes('localai') || apiBase.includes('local'))) {
    return 'localai';
  }
  
  // Hostname patterns
  if (apiBase.includes('localhost:11434') || apiBase.includes('127.0.0.1:11434')) {
    return 'ollama';
  }
  if (apiBase.includes('localhost:1234') || apiBase.includes('127.0.0.1:1234')) {
    return 'lmstudio';
  }
  
  // Default to OpenAI-compatible for unknown endpoints
  return 'openai-compatible';
}