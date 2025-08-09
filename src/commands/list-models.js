import chalk from 'chalk';
import { createClient } from '../api/client.js';
import { explainConnectionError } from '../utils/error-handler.js';

export async function listModelsCommand(cfg) {
  const client = createClient(cfg);
  
  try {
    console.log(chalk.cyan(`Fetching models from ${cfg.apiBase}...\n`));
    
    const models = await client.listModels();
    
    if (models.length === 0) {
      console.log(chalk.yellow('No models found or API does not support model listing.'));
      console.log(chalk.dim('Make sure your server is running and has models loaded.'));
      return;
    }
    
    console.log(chalk.green(`Found ${models.length} model(s):\n`));
    models.forEach(model => {
      console.log(chalk.white(`  • ${model}`));
    });
    
    console.log(chalk.dim(`\nProvider: ${client.getProviderName()}`));
    console.log(chalk.dim(`Current default: ${cfg.model}`));
    
  } catch (err) {
    console.error(chalk.red(explainConnectionError(cfg, err)));
    process.exit(1);
  }
}