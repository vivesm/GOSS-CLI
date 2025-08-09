import chalk from 'chalk';

export function validateConfig(cfg) {
  const errors = [];
  const warnings = [];

  // Validate temperature
  if (cfg.temperature < 0 || cfg.temperature > 2) {
    warnings.push('Temperature should be between 0 and 2');
  }

  // Validate max tokens
  if (cfg.maxTokens < 1 || cfg.maxTokens > 32000) {
    warnings.push('Max tokens should be between 1 and 32000');
  }

  // Check conflicting stream flags
  if (cfg.stream === false && cfg.noStream === false) {
    warnings.push('Both --stream and --no-stream specified. Using --no-stream.');
    cfg.stream = false;
  }

  // Validate API base URL
  if (!cfg.apiBase.startsWith('http')) {
    errors.push('API base must start with http:// or https://');
  }

  // Provider-specific validations
  if (cfg.provider === 'openai' && !cfg.apiKey) {
    errors.push('OpenAI provider requires OPENAI_API_KEY environment variable');
  }

  // Display warnings
  if (warnings.length > 0) {
    warnings.forEach(warning => {
      console.error(chalk.yellow(`Warning: ${warning}`));
    });
  }

  // Display errors and exit if any
  if (errors.length > 0) {
    console.error(chalk.red('Configuration errors:'));
    errors.forEach(error => {
      console.error(chalk.red(`  • ${error}`));
    });
    process.exit(1);
  }

  return cfg;
}