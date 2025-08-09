import axios from 'axios';
import { BaseProvider } from './base.js';
import { logDebug } from '../utils/logger.js';
import { Readable } from 'stream';

export class OllamaProvider extends BaseProvider {
  constructor(config) {
    super(config);
    this.apiBase = config.apiBase || 'http://localhost:11434';
    this.http = axios.create({
      baseURL: this.apiBase,
      timeout: 60_000,
      headers: { 'Content-Type': 'application/json' },
    });
  }

  async chatComplete({ messages, temperature, maxTokens, stream }) {
    // Convert OpenAI format to Ollama format
    const body = {
      model: this.config.model || 'llama2',
      messages,
      stream,
      options: {
        temperature,
        num_predict: maxTokens
      }
    };

    logDebug(this.config, 'OLLAMA REQUEST', JSON.stringify(body, null, 2));

    if (!stream) {
      const { data } = await this.http.post('/api/chat', body);
      // Convert Ollama response to OpenAI format
      const openaiResponse = {
        choices: [{
          message: { content: data.message?.content || '' }
        }]
      };
      logDebug(this.config, 'OLLAMA RESPONSE', JSON.stringify(openaiResponse, null, 2));
      return { type: 'final', data: openaiResponse };
    }

    // Streaming response
    const res = await this.http.post('/api/chat', body, {
      responseType: 'stream',
    });

    // Transform Ollama stream to OpenAI format
    const transformStream = new Readable({
      read() {}
    });

    res.data.on('data', (chunk) => {
      try {
        const lines = chunk.toString().split('\n').filter(Boolean);
        for (const line of lines) {
          const json = JSON.parse(line);
          if (json.done) {
            transformStream.push('data: [DONE]\n\n');
          } else if (json.message?.content) {
            const openaiChunk = {
              choices: [{
                delta: { content: json.message.content }
              }]
            };
            transformStream.push(`data: ${JSON.stringify(openaiChunk)}\n\n`);
          }
        }
      } catch (err) {
        logDebug(this.config, 'OLLAMA STREAM ERROR', err);
      }
    });

    res.data.on('end', () => transformStream.push(null));
    res.data.on('error', (err) => transformStream.destroy(err));

    return { type: 'stream', stream: transformStream };
  }

  async listModels() {
    try {
      const response = await this.http.get('/api/tags', { timeout: 5000 });
      const models = response.data?.models || [];
      return models.map(m => m.name).filter(Boolean);
    } catch {
      return [];
    }
  }

  getName() {
    return 'ollama';
  }
}