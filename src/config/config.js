import 'dotenv/config';
import { validateConfig } from '../utils/config-validator.js';

const DEFAULTS = {
  apiBase: 'http://localhost:1234/v1',
  model: 'openai/gpt-oss-20b',
  temperature: 0.7,
  maxTokens: 2048,
  stream: true,
  debug: false,
};

export function loadConfig(cli) {
  let cfg = {
    apiBase: cli.apiBase || process.env.API_BASE || DEFAULTS.apiBase,
    model: cli.model || process.env.MODEL || DEFAULTS.model,
    temperature: num(cli.temperature || process.env.TEMPERATURE, DEFAULTS.temperature),
    maxTokens: int(cli.maxTokens || process.env.MAX_TOKENS, DEFAULTS.maxTokens),
    stream: cli.stream ?? DEFAULTS.stream,
    debug: !!cli.debug,
    save: !!cli.save,
    contextFile: cli.contextFile || null,
    provider: cli.provider || process.env.PROVIDER || null,
    apiKey: process.env.OPENAI_API_KEY || null,
  };
  
  // Validate and fix configuration
  cfg = validateConfig(cfg);
  
  return cfg;
}

function num(v, d) { const n = Number(v); return Number.isFinite(n) ? n : d; }
function int(v, d) { const n = parseInt(v, 10); return Number.isInteger(n) ? n : d; }