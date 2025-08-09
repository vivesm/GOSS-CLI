export class BaseProvider {
  constructor(config) {
    this.config = config;
  }

  async chatComplete({ messages, temperature, maxTokens, stream }) {
    throw new Error('chatComplete must be implemented by provider');
  }

  async listModels() {
    throw new Error('listModels must be implemented by provider');
  }

  getName() {
    return 'base';
  }
}