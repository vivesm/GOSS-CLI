import axios from 'axios';
import { BaseProvider } from './base.js';
import { logDebug } from '../utils/logger.js';

export class OpenAICompatibleProvider extends BaseProvider {
  constructor(config) {
    super(config);
    this.http = axios.create({
      baseURL: config.apiBase,
      timeout: 60_000,
      headers: { 
        'Content-Type': 'application/json',
        ...(config.apiKey ? { 'Authorization': `Bearer ${config.apiKey}` } : {})
      },
    });
  }

  async chatComplete({ messages, temperature, maxTokens, stream }) {
    const body = {
      model: this.config.model,
      messages,
      temperature,
      max_tokens: maxTokens,
      stream,
    };

    logDebug(this.config, 'REQUEST', JSON.stringify(body, null, 2));

    if (!stream) {
      const { data } = await this.http.post('/chat/completions', body);
      logDebug(this.config, 'RESPONSE', JSON.stringify(data, null, 2));
      return { type: 'final', data };
    }

    const res = await this.http.post('/chat/completions', body, {
      responseType: 'stream',
    });

    return { type: 'stream', stream: res.data };
  }

  async listModels() {
    try {
      const response = await this.http.get('/models', { timeout: 5000 });
      const models = response.data?.data || [];
      return models.map(m => m.id || m.name).filter(Boolean);
    } catch {
      return [];
    }
  }

  getName() {
    return 'openai-compatible';
  }
}