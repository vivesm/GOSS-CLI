import chalk from 'chalk';
import { createClient } from '../api/client.js';
import { openAIToText } from '../api/converter.js';
import { explainConnectionError } from '../utils/error-handler.js';
import { saveConversation, loadContextFile } from '../utils/file-logger.js';

export async function singlePrompt(cfg, prompt) {
  const client = createClient(cfg);
  let messages = [];
  
  // Load context file if provided
  if (cfg.contextFile) {
    try {
      const context = await loadContextFile(cfg.contextFile);
      messages.push(...context);
      if (cfg.debug) console.error(chalk.dim(`Loaded ${context.length} messages from context file`));
    } catch (err) {
      console.error(chalk.yellow(`Warning: ${err.message}`));
    }
  }
  
  messages.push({ role: 'user', content: prompt });

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
        const responseText = await streamToStdout(r.stream);
        process.stdout.write('\n');
        messages.push({ role: 'assistant', content: responseText });
        
        // Save conversation if requested
        if (cfg.save) {
          const filepath = await saveConversation(messages, prompt);
          console.error(chalk.dim(`\nConversation saved to: ${filepath}`));
        }
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
    messages.push({ role: 'assistant', content: text });
    
    // Save conversation if requested
    if (cfg.save) {
      const filepath = await saveConversation(messages, prompt);
      console.error(chalk.dim(`\nConversation saved to: ${filepath}`));
    }
  } catch (err) {
    console.error(chalk.red(explainConnectionError(cfg, err)));
    process.exit(1);
  }
}

async function streamToStdout(stream) {
  return new Promise((resolve, reject) => {
    let fullText = '';
    stream.on('data', (chunk) => {
      const lines = chunk.toString().split('\n').filter(Boolean);
      for (const line of lines) {
        if (!line.startsWith('data:')) continue;
        const payload = line.slice(5).trim();
        if (payload === '[DONE]') {
          resolve(fullText);
          return;
        }
        try {
          const json = JSON.parse(payload);
          const delta = json?.choices?.[0]?.delta?.content || '';
          if (delta) {
            process.stdout.write(delta);
            fullText += delta;
          }
        } catch { /* ignore parse errors in mixed chunks */ }
      }
    });
    stream.on('end', () => resolve(fullText));
    stream.on('error', reject);
  });
}