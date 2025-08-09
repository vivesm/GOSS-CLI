// Gemini -> OpenAI messages
export function geminiToOpenAI(contents) {
  // contents: [{ role: 'user'|'model'|'system', parts:[{text:string}|...] }, ...]
  const messages = [];
  for (const c of contents || []) {
    const role = mapRole(c.role);
    const text = (c.parts || [])
      .map(p => p.text)
      .filter(Boolean)
      .join('\n');
    if (text) messages.push({ role, content: text });
  }
  return messages;
}

function mapRole(role) {
  if (role === 'model') return 'assistant';
  if (role === 'user' || role === 'assistant' || role === 'system') return role;
  return 'user';
}

// OpenAI -> plain text
export function openAIToText(resp) {
  // Non-streamed: resp.choices[0].message.content
  const ch = resp?.choices?.[0];
  return ch?.message?.content ?? '';
}