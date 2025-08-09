import inquirer from 'inquirer';
import chalk from 'chalk';
import { createClient } from '../api/client.js';
import { openAIToText } from '../api/converter.js';
import { explainConnectionError } from '../utils/error-handler.js';
import { saveConversation, loadContextFile } from '../utils/file-logger.js';

export async function chatCommand(cfg) {
  const client = createClient(cfg);
  let history = [];
  
  // Load context file if provided
  if (cfg.contextFile) {
    try {
      const context = await loadContextFile(cfg.contextFile);
      history.push(...context);
      console.log(chalk.dim(`Loaded ${context.length} messages from context file\n`));
    } catch (err) {
      console.error(chalk.yellow(`Warning: ${err.message}\n`));
    }
  }

  while (true) {
    const { prompt } = await inquirer.prompt([
      { type: 'input', name: 'prompt', message: chalk.cyan('You:'), validate: v => !!v || 'Enter a prompt (or /q to quit)' }
    ]);
    if (prompt.trim() === '/q') {
      // Save on quit if requested
      if (cfg.save && history.length > 0) {
        const filepath = await saveConversation(history);
        console.log(chalk.dim(`\nConversation saved to: ${filepath}`));
      }
      break;
    }

    history.push({ role: 'user', content: prompt });

    try {
      if (cfg.stream) {
        const r = await client.chatComplete({
          messages: history,
          temperature: cfg.temperature,
          maxTokens: cfg.maxTokens,
          stream: true,
        });
        process.stdout.write(chalk.green('Assistant: '));
        if (r.type === 'stream') {
          let buffer = '';
          await new Promise((resolve, reject) => {
            r.stream.on('data', (chunk) => {
              const lines = chunk.toString().split('\n').filter(Boolean);
              for (const line of lines) {
                if (!line.startsWith('data:')) continue;
                const payload = line.slice(5).trim();
                if (payload === '[DONE]') return;
                try {
                  const json = JSON.parse(payload);
                  const delta = json?.choices?.[0]?.delta?.content || '';
                  if (delta) {
                    buffer += delta;
                    process.stdout.write(delta);
                  }
                } catch { /* ignore */ }
              }
            });
            r.stream.on('end', resolve);
            r.stream.on('error', reject);
          });
          process.stdout.write('\n');
          history.push({ role: 'assistant', content: buffer });
        }
      } else {
        const r = await client.chatComplete({
          messages: history,
          temperature: cfg.temperature,
          maxTokens: cfg.maxTokens,
          stream: false,
        });
        const text = openAIToText(r.data);
        console.log(chalk.green('Assistant:'), text);
        history.push({ role: 'assistant', content: text });
      }
    } catch (err) {
      console.error(chalk.red(explainConnectionError(cfg, err)));
      break;
    }
  }
}