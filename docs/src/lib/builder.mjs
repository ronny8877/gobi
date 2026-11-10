export const generators = [
  ['literal', 'Literal text', 'Hello, Gobi'],
  ['name', 'Full name', 'User()'],
  ['email', 'Email address', 'User(email)'],
  ['uuid', 'UUID', 'Uuid()'],
  ['integer', 'Integer', 'Int(min=1,max=100)'],
  ['float', 'Decimal', 'Float(min=1,max=100,precision=2)'],
  ['boolean', 'Boolean', 'Bool(true)'],
  ['time', 'Timestamp', 'Time(iso)'],
  ['array', 'Array', 'Array(len=3,type=User(email))'],
  ['reference', 'Reusable profile', 'Ref(ref.profile)'],
];
export function initialFields() {
  return [
    { id: 1, name: 'id', kind: 'uuid', value: 'Uuid()' },
    { id: 2, name: 'name', kind: 'name', value: 'User()' },
    { id: 3, name: 'email', kind: 'email', value: 'User(email)' },
    { id: 4, name: 'active', kind: 'boolean', value: 'Bool(true)' },
  ];
}
export function validateBuilder(state) {
  const errors = [];
  if (!state.path.startsWith('/') || /[\s?#]/.test(state.path))
    errors.push(
      'Route path must start with / and contain no spaces, query or fragment.',
    );
  if (!['GET', 'POST', 'PUT', 'PATCH', 'DELETE'].includes(state.method))
    errors.push('Choose a supported HTTP method.');
  if (
    !Number.isInteger(Number(state.port)) ||
    Number(state.port) < 1 ||
    Number(state.port) > 65535
  )
    errors.push('Port must be a whole number between 1 and 65535.');
  if (
    state.prefix &&
    (!state.prefix.startsWith('/') || /[\s{}?#]/.test(state.prefix))
  )
    errors.push(
      'Prefix must start with / and contain no spaces, braces, query or fragment.',
    );
  if (
    state.array &&
    (!Number.isInteger(Number(state.count)) ||
      Number(state.count) < 0 ||
      Number(state.count) > 10000)
  )
    errors.push('Array length must be between 0 and 10000.');
  if (
    !Number.isInteger(Number(state.latency)) ||
    Number(state.latency) < 0 ||
    Number(state.latency) > 30000
  )
    errors.push('Latency must be between 0 and 30000 milliseconds.');
  if (
    Number(state.failRate) < 0 ||
    Number(state.failRate) > 1 ||
    !Number.isFinite(Number(state.failRate))
  )
    errors.push('Failure rate must be between 0 and 1.');
  if (!state.fields.length) errors.push('Add at least one response field.');
  if (
    ['/health', '/auth'].includes(state.prefix.replace(/\/+$/, '') + state.path)
  )
    errors.push('The /health and /auth endpoints are reserved by Gobi.');
  const names = new Set();
  for (const field of state.fields) {
    if (!field.name.trim()) errors.push('Every field needs a name.');
    if (names.has(field.name))
      errors.push(`Field names must be unique: ${field.name}`);
    names.add(field.name);
    if (!generators.some((g) => g[0] === field.kind))
      errors.push('Choose a supported field type.');
    if (field.kind !== 'literal' && !field.value.trim())
      errors.push(
        `Provide a generator expression for ${field.name || 'the field'}.`,
      );
  }
  return [...new Set(errors)];
}
export function buildConfig(state) {
  const errors = validateBuilder(state);
  if (errors.length) throw new Error(errors.join(' '));
  const route = {
    method: state.method,
    path: state.path,
    ...(state.array ? { responseType: `Array(${Number(state.count)})` } : {}),
    response: Object.fromEntries(state.fields.map((f) => [f.name, f.value])),
  };
  return {
    config: {
      port: Number(state.port),
      prefix: state.prefix.replace(/\/+$/, ''),
      latency: Number(state.latency),
      failRate: Number(state.failRate),
    },
    ...(state.fields.some(
      (f) => f.kind === 'reference' || f.value.includes('Ref('),
    )
      ? {
          ref: {
            profile: { id: 'Uuid()', name: 'User()', email: 'User(email)' },
          },
        }
      : {}),
    api: [route],
  };
}
