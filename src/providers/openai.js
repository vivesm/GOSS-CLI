import { OpenAICompatibleProvider } from './openai-compatible.js';

export class OpenAIProvider extends OpenAICompatibleProvider {
  constructor(config) {
    // Force OpenAI API base if not specified
    const openaiConfig = {
      ...config,
      apiBase: config.apiBase || 'https://api.openai.com/v1',
      model: config.model || 'gpt-3.5-turbo'
    };
    
    if (!openaiConfig.apiKey) {
      throw new Error('OpenAI provider requires OPENAI_API_KEY environment variable');
    }
    
    super(openaiConfig);
  }

  getName() {
    return 'openai';
  }
}