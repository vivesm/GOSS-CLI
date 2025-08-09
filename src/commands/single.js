import chalk from 'chalk';
import { createClient } from '../api/client.js';
import { openAIToText } from '../api/converter.js';
import { explainConnectionError } from '../utils/error-handler.js';

export async function singlePrompt(cfg, prompt) {
  const client = createClient(cfg);
  const messages = [{ role: 'user', content: prompt }];

  try {
    if (cfg.stream) {
      // Streamed
      const r = await client.chatComplete({
        messages,
        temperature: cfg.temperature,
        maxTokens: cfg.maxTokens,
        stream: true,
      });

      if (r.type === 'stream') {
        await streamToStdout(r.stream);
        process.stdout.write('\n');
      }
      return;
    }

    // Non-streamed
    const r = await client.chatComplete({
      messages,
      temperature: cfg.temperature,
      maxTokens: cfg.maxTokens,
      stream: false,
    });

    const text = openAIToText(r.data);
    process.stdout.write(chalk.reset(text) + '\n');
  } catch (err) {
    console.error(chalk.red(explainConnectionError(cfg, err)));
    process.exit(1);
  }
}

async function streamToStdout(stream) {
  return new Promise((resolve, reject) => {
    stream.on('data', (chunk) => {
      const lines = chunk.toString().split('\n').filter(Boolean);
      for (const line of lines) {
        if (!line.startsWith('data:')) continue;
        const payload = line.slice(5).trim();
        if (payload === '[DONE]') return; // end signal
        try {
          const json = JSON.parse(payload);
          const delta = json?.choices?.[0]?.delta?.content || '';
          if (delta) process.stdout.write(delta);
        } catch { /* ignore parse errors in mixed chunks */ }
      }
    });
    stream.on('end', resolve);
    stream.on('error', reject);
  });
}