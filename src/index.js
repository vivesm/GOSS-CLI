import { Command } from 'commander';
import { loadConfig } from './config/config.js';
import { chatCommand } from './commands/chat.js';
import { validateModel } from './utils/model-detector.js';

const program = new Command();

program
  .name('goss')
  .description('GOSS-CLI: Universal local/remote model CLI')
  .option('--provider <name>', 'Provider (lmstudio, ollama, openai, localai)', process.env.PROVIDER)
  .option('--api-base <url>', 'API endpoint URL', process.env.API_BASE)
  .option('--model <name>', 'Model name', process.env.MODEL)
  .option('--temperature <num>', 'Temperature', process.env.TEMPERATURE)
  .option('--max-tokens <n>', 'Max tokens', process.env.MAX_TOKENS)
  .option('--debug', 'Verbose logging', false)
  .option('--no-stream', 'Disable streaming', false)
  .option('--save', 'Save conversation to logs/ folder', false)
  .option('--context-file <path>', 'Pre-load conversation from file');

program
  .command('chat')
  .description('Interactive chat mode')
  .action(async () => {
    const cfg = loadConfig(program.opts());
    await validateModel(cfg);
    await chatCommand(cfg);
  });

program
  .argument('[prompt...]', 'Single-prompt mode')
  .action(async (promptParts) => {
    const cfg = loadConfig(program.opts());
    if (promptParts.length === 0) {
      program.outputHelp();
      process.exit(1);
    }
    await validateModel(cfg);
    const { singlePrompt } = await import('./commands/single.js');
    await singlePrompt(cfg, promptParts.join(' '));
  });

program.parseAsync(process.argv).catch((e) => {
  console.error('Error:', e?.message || e);
  process.exit(1);
});