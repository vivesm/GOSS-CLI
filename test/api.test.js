import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';
import nock from 'nock';
import { createClient } from '../src/api/client.js';
import { Readable } from 'stream';

describe('API Client', () => {
  const cfg = {
    apiBase: 'http://localhost:1234/v1',
    model: 'test-model',
    debug: false
  };

  beforeEach(() => {
    nock.cleanAll();
  });

  afterEach(() => {
    nock.cleanAll();
  });

  describe('Non-streaming', () => {
    it('should make successful completion request', async () => {
      const mockResponse = {
        choices: [{
          message: { content: 'Hello from AI' }
        }]
      };

      nock('http://localhost:1234')
        .post('/v1/chat/completions')
        .reply(200, mockResponse);

      const client = createClient(cfg);
      const result = await client.chatComplete({
        messages: [{ role: 'user', content: 'Hello' }],
        temperature: 0.7,
        maxTokens: 100,
        stream: false
      });

      expect(result.type).toBe('final');
      expect(result.data).toEqual(mockResponse);
    });

    it('should handle API errors gracefully', async () => {
      nock('http://localhost:1234')
        .post('/v1/chat/completions')
        .reply(500, { error: 'Internal Server Error' });

      const client = createClient(cfg);
      
      await expect(client.chatComplete({
        messages: [{ role: 'user', content: 'Hello' }],
        temperature: 0.7,
        maxTokens: 100,
        stream: false
      })).rejects.toThrow();
    });

    it.skip('should handle connection refused', async () => {
      nock('http://localhost:1234')
        .post('/v1/chat/completions')
        .replyWithError({ code: 'ECONNREFUSED', message: 'Connection refused' });

      const client = createClient(cfg);
      
      try {
        await client.chatComplete({
          messages: [{ role: 'user', content: 'Hello' }],
          temperature: 0.7,
          maxTokens: 100,
          stream: false
        });
        expect.fail('Should have thrown an error');
      } catch (err) {
        expect(err.code).toBe('ECONNREFUSED');
      }
    });
  });

  describe('Streaming', () => {
    it('should handle streaming responses', async () => {
      const streamData = [
        'data: {"choices":[{"delta":{"content":"Hello"}}]}\n',
        'data: {"choices":[{"delta":{"content":" world"}}]}\n',
        'data: [DONE]\n'
      ];

      const stream = new Readable({
        read() {
          streamData.forEach(chunk => this.push(chunk));
          this.push(null);
        }
      });

      nock('http://localhost:1234')
        .post('/v1/chat/completions')
        .reply(200, () => stream);

      const client = createClient(cfg);
      const result = await client.chatComplete({
        messages: [{ role: 'user', content: 'Hello' }],
        temperature: 0.7,
        maxTokens: 100,
        stream: true
      });

      expect(result.type).toBe('stream');
      expect(result.stream).toBeDefined();
    });
  });
});