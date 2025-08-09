export function explainConnectionError(cfg, err) {
  const base = cfg.apiBase;
  const msg = err?.code === 'ECONNREFUSED' || err?.message?.includes('ECONNREFUSED')
    ? `LM Studio API unreachable at ${base}.
Start LM Studio, enable "Local Server", and ensure the API is at ${base}`
    : err?.message || String(err);
  return msg;
}