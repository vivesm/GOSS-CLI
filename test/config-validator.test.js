import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { validateConfig } from '../src/utils/config-validator.js';

describe('Config Validator', () => {
  let consoleErrorSpy;
  let processExitSpy;

  beforeEach(() => {
    consoleErrorSpy = vi.spyOn(console, 'error').mockImplementation(() => {});
    processExitSpy = vi.spyOn(process, 'exit').mockImplementation(() => {});
  });

  afterEach(() => {
    consoleErrorSpy.mockRestore();
    processExitSpy.mockRestore();
  });

  it('should pass valid config', () => {
    const config = {
      apiBase: 'http://localhost:1234/v1',
      temperature: 0.7,
      maxTokens: 1000,
      stream: true
    };
    
    const result = validateConfig(config);
    expect(result).toEqual(config);
    expect(processExitSpy).not.toHaveBeenCalled();
  });

  it('should warn about extreme temperature', () => {
    const config = {
      apiBase: 'http://localhost:1234/v1',
      temperature: 3.0,
      maxTokens: 1000
    };
    
    validateConfig(config);
    expect(consoleErrorSpy).toHaveBeenCalledWith(
      expect.stringContaining('Temperature should be between 0 and 2')
    );
  });

  it('should error on invalid API base', () => {
    const config = {
      apiBase: 'invalid-url',
      temperature: 0.7,
      maxTokens: 1000
    };
    
    validateConfig(config);
    expect(processExitSpy).toHaveBeenCalledWith(1);
  });

  it('should error when OpenAI provider has no API key', () => {
    const config = {
      apiBase: 'https://api.openai.com/v1',
      provider: 'openai',
      temperature: 0.7,
      maxTokens: 1000
    };
    
    validateConfig(config);
    expect(processExitSpy).toHaveBeenCalledWith(1);
  });
});