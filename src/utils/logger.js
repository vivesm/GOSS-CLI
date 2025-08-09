export function logDebug(cfg, ...args) {
  if (cfg.debug) console.error('[DEBUG]', ...args);
}