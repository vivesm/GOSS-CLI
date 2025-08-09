import { describe, it, expect } from 'vitest';
import { createProvider } from '../src/providers/index.js';

describe('Provider Detection', () => {
  it('should detect LM Studio by port', () => {
    const config = { apiBase: 'http://localhost:1234/v1', debug: false };
    const provider = createProvider(config);
    expect(provider.getName()).toBe('openai-compatible'); // LM Studio uses OpenAI-compatible provider
  });

  it('should detect Ollama by port', () => {
    const config = { apiBase: 'http://localhost:11434/v1', debug: false };
    const provider = createProvider(config);
    expect(provider.getName()).toBe('ollama');
  });

  it('should detect OpenAI by domain', () => {
    const config = { 
      apiBase: 'https://api.openai.com/v1',
      apiKey: 'test-key',
      debug: false
    };
    const provider = createProvider(config);
    expect(provider.getName()).toBe('openai');
  });

  it('should use explicit provider over detection', () => {
    const config = { 
      apiBase: 'http://localhost:11434/v1',
      provider: 'openai-compatible',
      debug: false
    };
    const provider = createProvider(config);
    expect(provider.getName()).toBe('openai-compatible');
  });

  it('should default to openai-compatible for unknown URLs', () => {
    const config = { 
      apiBase: 'https://custom-api.example.com/v1',
      debug: false
    };
    const provider = createProvider(config);
    expect(provider.getName()).toBe('openai-compatible');
  });
});