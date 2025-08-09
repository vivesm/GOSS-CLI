import axios from 'axios';
import { logDebug } from '../utils/logger.js';

export function createClient(cfg) {
  const http = axios.create({
    baseURL: cfg.apiBase,
    timeout: 60_000,
    headers: { 'Content-Type': 'application/json' },
  });

  async function chatComplete({ messages, temperature, maxTokens, stream }) {
    const body = {
      model: cfg.model,
      messages,
      temperature,
      max_tokens: maxTokens,
      stream,
    };

    logDebug(cfg, 'REQUEST', JSON.stringify(body, null, 2));

    if (!stream) {
      const { data } = await http.post('/chat/completions', body);
      logDebug(cfg, 'RESPONSE', JSON.stringify(data, null, 2));
      return { type: 'final', data };
    }

    // Streamed: Server-Sent Events-like chunks: "data: {...}\n\n", ending with "data: [DONE]"
    const res = await http.post('/chat/completions', body, {
      responseType: 'stream',
    });

    return { type: 'stream', stream: res.data };
  }

  return { chatComplete };
}