export const profile = {
  name: 'User()',
  email: 'User(email)',
  id: 'Uuid()',
  avatar: 'DiceBearImage(initials)',
};
export const exampleCases = [
  {
    id: 'array',
    label: 'Arrays',
    method: 'GET',
    path: '/api/users',
    description: 'One template. Six independently generated users.',
    schema: {
      method: 'GET',
      path: '/users',
      responseType: 'Array(6)',
      response: {
        name: 'User()',
        email: 'User(email)',
        role: 'Company(jobTitle)',
        active: 'Bool(true)',
      },
    },
  },
  {
    id: 'ref',
    label: 'References',
    method: 'GET',
    path: '/api/profile',
    description: 'Reuse a profile without repeating the field definitions.',
    schema: {
      method: 'GET',
      path: '/profile',
      response: {
        owner: 'Ref(ref.profile)',
        team: 'Array(len=3,type=Ref(ref.profile))',
      },
    },
  },
  {
    id: 'types',
    label: 'Generators',
    method: 'GET',
    path: '/api/metrics',
    description: 'Numbers, booleans and identifiers keep their JSON types.',
    schema: {
      method: 'GET',
      path: '/metrics',
      response: {
        id: 'Uuid()',
        visitors: 'Int(min=1200,max=8000)',
        conversion: 'Float(min=1,max=8,precision=2)',
        healthy: 'Bool(true)',
        updated: 'Time(iso)',
      },
    },
  },
  {
    id: 'methods',
    label: 'HTTP methods',
    method: 'POST',
    path: '/api/events',
    description:
      'A POST route with required JSON fields and a generated response.',
    schema: {
      method: 'POST',
      path: '/events',
      validate: { body: ['name'] },
      response: {
        id: 'Uuid()',
        message: 'Message(Event created)',
        created: 'Time(iso)',
      },
    },
  },
] as const;
export const samplePeople = [
  {
    name: 'Alex Morgan',
    email: 'alex@example.com',
    role: 'Designer',
    active: true,
  },
  {
    name: 'Jordan Lee',
    email: 'jordan@example.com',
    role: 'Engineer',
    active: true,
  },
  {
    name: 'Sam Rivera',
    email: 'sam@example.com',
    role: 'Product manager',
    active: true,
  },
  {
    name: 'Casey Chen',
    email: 'casey@example.com',
    role: 'Researcher',
    active: true,
  },
  {
    name: 'Taylor Reed',
    email: 'taylor@example.com',
    role: 'Developer',
    active: true,
  },
  {
    name: 'Robin Patel',
    email: 'robin@example.com',
    role: 'Designer',
    active: true,
  },
];
const sampleProfile = {
  ...samplePeople[0],
  id: 'sample-user-001',
  avatar: '/favicon.svg',
};
export const sampleDashboard = {
  visitors: 4826,
  revenue: 18740,
  orders: 316,
  conversion: 3.8,
  series: [
    { day: 'Mon', visits: 480 },
    { day: 'Tue', visits: 620 },
    { day: 'Wed', visits: 540 },
    { day: 'Thu', visits: 780 },
    { day: 'Fri', visits: 890 },
    { day: 'Sat', visits: 650 },
    { day: 'Sun', visits: 866 },
  ],
  customers: samplePeople.slice(0, 5).map((p, i) => ({
    ...p,
    id: `sample-${i}`,
    orders: 12 + i * 3,
    spend: 240 + i * 83,
  })),
  updated: 'Sample fixture',
};
export function sampleResult(id: string): unknown {
  switch (id) {
    case 'array':
      return samplePeople;
    case 'ref':
      return {
        owner: sampleProfile,
        team: [
          sampleProfile,
          { ...sampleProfile, name: 'Jordan Lee', id: 'sample-user-002' },
          { ...sampleProfile, name: 'Sam Rivera', id: 'sample-user-003' },
        ],
      };
    case 'types':
      return {
        id: 'sample-metric-001',
        visitors: 4826,
        conversion: 3.8,
        healthy: true,
        updated: '2026-01-01T12:00:00Z',
      };
    default:
      return {
        id: 'sample-event-001',
        message: 'Event created',
        created: '2026-01-01T12:00:00Z',
      };
  }
}
export function completeExample(index: number) {
  return {
    config: { port: 8080, prefix: '/api' },
    ref: { profile },
    api: [exampleCases[index].schema],
  };
}
export function endpoint(base: string, path: string) {
  const clean = base.trim().replace(/\/+$/, '');
  if (clean.startsWith('/') && !clean.startsWith('//')) return clean + path;
  const url = new URL(clean);
  if (
    !['http:', 'https:'].includes(url.protocol) ||
    url.username ||
    url.password ||
    url.search ||
    url.hash
  )
    throw new Error(
      'Use an HTTP(S) server URL without credentials, query parameters or a fragment.',
    );
  return clean + path;
}
export async function fetchJSON(
  url: string,
  method = 'GET',
  body?: string,
  signal?: AbortSignal,
) {
  const timeout = AbortSignal.timeout(12000);
  const response = await fetch(url, {
    method,
    headers: body ? { 'Content-Type': 'application/json' } : undefined,
    body: method === 'GET' ? undefined : body,
    signal: signal ? AbortSignal.any([timeout, signal]) : timeout,
    credentials: 'omit',
  });
  let json: unknown;
  try {
    json = await response.json();
  } catch {
    throw new Error(
      `The server returned ${response.status} without valid JSON.`,
    );
  }
  if (!response.ok) {
    const message =
      typeof json === 'object' && json && 'error' in json
        ? String(json.error)
        : response.statusText;
    throw new Error(`HTTP ${response.status}: ${message}`);
  }
  return { data: json, status: response.status };
}
export function downloadJSON(value: unknown, filename: string) {
  const url = URL.createObjectURL(
    new Blob([JSON.stringify(value, null, 2) + '\n'], {
      type: 'application/json',
    }),
  );
  const link = document.createElement('a');
  link.href = url;
  link.download = filename;
  link.click();
  setTimeout(() => URL.revokeObjectURL(url), 1000);
}
