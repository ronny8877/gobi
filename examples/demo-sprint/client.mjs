/** A small fetch client for use with Gobi or another JSON API. */
export function createClient(baseURL, fetcher = fetch) {
  const base = baseURL.replace(/\/+$/, '');
  return async function request(path, { method = 'GET', body, signal } = {}) {
    const response = await fetcher(`${base}${path}`, {
      method,
      headers: body === undefined ? undefined : { 'Content-Type': 'application/json' },
      body: body === undefined ? undefined : JSON.stringify(body),
      signal: signal ? AbortSignal.any([signal, AbortSignal.timeout(10000)]) : AbortSignal.timeout(10000),
    });
    const data = await response.json();
    if (!response.ok) throw new Error(`HTTP ${response.status}: ${data.error || 'Request failed'}`);
    return data;
  };
}
