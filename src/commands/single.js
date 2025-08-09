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
    let hasReceivedData = false;
    
    stream.on('data', (chunk) => {
      hasReceivedData = true;
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
        } catch (parseError) {
          // Log malformed JSON for debugging but continue
          if (process.env.DEBUG) {
            console.error(`\n[DEBUG] Malformed JSON chunk: ${payload}`);
          }
        }
      }
    });
    
    stream.on('end', () => {
      if (!hasReceivedData && !fullText) {
        reject(new Error('Stream ended without receiving any data'));
      } else {
        resolve(fullText);
      }
    });
    
    stream.on('error', (err) => {
      if (fullText) {
        // If we got partial content, show it and warn
        console.error(chalk.yellow('\nConnection interrupted. Partial response received.'));
        resolve(fullText);
      } else {
        reject(err);
      }
    });
    
    // Timeout for stuck streams
    const timeout = setTimeout(() => {
      reject(new Error('Stream timeout - no data received in 30 seconds'));
    }, 30000);
    
    stream.on('data', () => clearTimeout(timeout));
    stream.on('end', () => clearTimeout(timeout));
    stream.on('error', () => clearTimeout(timeout));
  });
}