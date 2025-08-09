import fs from 'fs/promises';
import path from 'path';
import { fileURLToPath } from 'url';
import os from 'os';

const __dirname = path.dirname(fileURLToPath(import.meta.url));

export async function saveConversation(messages, prompt = null) {
  const logsDir = path.join(process.cwd(), 'logs');
  
  try {
    // Ensure logs directory exists with proper permissions
    await fs.mkdir(logsDir, { recursive: true });
    
    // Test write permissions
    const testFile = path.join(logsDir, '.write-test');
    await fs.writeFile(testFile, 'test');
    await fs.unlink(testFile);
  } catch (err) {
    if (err.code === 'EACCES' || err.code === 'EPERM') {
      throw new Error(`Cannot write to logs directory: ${logsDir}. Check permissions or use a different directory.`);
    }
    throw err;
  }
  
  // Create timestamp filename (Windows-safe)
  const timestamp = new Date().toISOString().replace(/[:.]/g, '-').replace('T', '_').slice(0, -5);
  const filename = `conversation_${timestamp}.txt`;
  const filepath = path.resolve(logsDir, filename);
  
  // Format conversation
  let content = `GOSS-CLI Conversation Log\n`;
  content += `Timestamp: ${new Date().toISOString()}\n`;
  content += `${'='.repeat(50)}\n\n`;
  
  if (prompt) {
    content += `Single Prompt Mode\n`;
    content += `${'='.repeat(50)}\n\n`;
  }
  
  for (const msg of messages) {
    const role = msg.role.charAt(0).toUpperCase() + msg.role.slice(1);
    content += `${role}:\n${msg.content}\n\n`;
    content += `${'-'.repeat(30)}\n\n`;
  }
  
  await fs.writeFile(filepath, content);
  return filepath;
}

export async function loadContextFile(filepath) {
  try {
    // Resolve relative paths and normalize for cross-platform
    const resolvedPath = path.resolve(process.cwd(), filepath);
    const content = await fs.readFile(resolvedPath, 'utf-8');
    const messages = [];
    
    // Simple parser for conversation format
    const sections = content.split(/^(User|Assistant|System):/gm);
    
    for (let i = 1; i < sections.length; i += 2) {
      const role = sections[i].toLowerCase();
      const content = sections[i + 1].trim().replace(/-{30,}/g, '').trim();
      
      if (content) {
        messages.push({
          role: role === 'assistant' ? 'assistant' : role === 'system' ? 'system' : 'user',
          content
        });
      }
    }
    
    return messages;
  } catch (err) {
    if (err.code === 'ENOENT') {
      throw new Error(`Context file not found: ${filepath}`);
    }
    throw err;
  }
}