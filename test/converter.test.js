import { describe, it, expect } from 'vitest';
import { geminiToOpenAI, openAIToText } from '../src/api/converter.js';

describe('converter', () => {
  it('maps gemini contents to openai messages', () => {
    const contents = [{ role: 'user', parts: [{ text: 'Hello' }, { text: 'World' }] }];
    const msgs = geminiToOpenAI(contents);
    expect(msgs).toEqual([{ role: 'user', content: 'Hello\nWorld' }]);
  });

  it('extracts text from openai response', () => {
    const resp = { choices: [{ message: { content: 'Hi!' } }] };
    expect(openAIToText(resp)).toBe('Hi!');
  });
});