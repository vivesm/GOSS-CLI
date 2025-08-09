import axios from 'axios';
import chalk from 'chalk';

export async function listAvailableModels(apiBase) {
  try {
    const response = await axios.get(`${apiBase}/models`, {
      timeout: 5000
    });
    
    const models = response.data?.data || [];
    return models.map(m => m.id || m.name).filter(Boolean);
  } catch (err) {
    // Silently fail - model listing is optional
    return [];
  }
}

export async function validateModel(cfg) {
  const models = await listAvailableModels(cfg.apiBase);
  
  if (models.length === 0) {
    // Can't validate, assume it's correct
    return true;
  }
  
  const modelExists = models.some(m => 
    m.toLowerCase() === cfg.model.toLowerCase()
  );
  
  if (!modelExists) {
    console.error(chalk.yellow(`Warning: Model '${cfg.model}' not found in available models.`));
    console.error(chalk.cyan('Available models:'));
    models.forEach(m => console.error(chalk.gray(`  - ${m}`)));
    console.error(chalk.dim('\nYou can continue anyway, or use --model to select a different one.\n'));
    return false;
  }
  
  return true;
}