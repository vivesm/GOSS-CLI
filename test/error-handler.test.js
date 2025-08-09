import { describe, it, expect } from 'vitest';
import { explainConnectionError } from '../src/utils/error-handler.js';

describe('Error Handler', () => {
  const cfg = { apiBase: 'http://localhost:1234/v1' };

  it('should explain ECONNREFUSED errors', () => {
    const error = { code: 'ECONNREFUSED' };
    const message = explainConnectionError(cfg, error);
    
    expect(message).toContain('LM Studio API unreachable');
    expect(message).toContain('http://localhost:1234/v1');
    expect(message).toContain('enable "Local Server"');
  });

  it('should handle message-based ECONNREFUSED', () => {
    const error = { message: 'connect ECONNREFUSED 127.0.0.1:1234' };
    const message = explainConnectionError(cfg, error);
    
    expect(message).toContain('LM Studio API unreachable');
  });

  it('should pass through other errors', () => {
    const error = { message: 'Invalid API key' };
    const message = explainConnectionError(cfg, error);
    
    expect(message).toBe('Invalid API key');
  });

  it('should handle errors without message', () => {
    const error = { toString: () => 'CustomError' };
    const message = explainConnectionError(cfg, error);
    
    expect(message).toBe('CustomError');
  });
});