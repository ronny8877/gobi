import { test } from 'node:test';
import assert from 'node:assert/strict';
import { readFileSync } from 'node:fs';
import {
  buildConfig,
  initialFields,
  validateBuilder,
} from '../src/lib/builder.mjs';
const starter = () => ({
  method: 'GET',
  path: '/users',
  port: 8080,
  prefix: '/api',
  array: true,
  count: 6,
  latency: 0,
  failRate: 0,
  fields: initialFields(),
});
test('builder exports a complete typed route', () => {
  const config = buildConfig(starter());
  assert.equal(config.config.port, 8080);
  assert.equal(config.api[0].responseType, 'Array(6)');
  assert.equal(config.api[0].response.active, 'Bool(true)');
  assert.equal(config.api[0].response.id, 'Uuid()');
});
test('reference fields include their reusable template', () => {
  const state = starter();
  state.fields.push({
    id: 5,
    name: 'owner',
    kind: 'reference',
    value: 'Ref(ref.profile)',
  });
  assert.equal(buildConfig(state).ref.profile.name, 'User()');
});
test('empty array and zero latency remain valid', () => {
  const state = starter();
  state.count = 0;
  assert.equal(buildConfig(state).api[0].responseType, 'Array(0)');
  assert.deepEqual(validateBuilder(state), []);
});
test('rejects invalid routes, ports, duplicates and empty schemas', () => {
  for (const patch of [
    { path: 'users' },
    { prefix: '', path: '/health' },
    { path: '/users?x=1' },
    { port: 0 },
    { port: 65536 },
    { count: 10001 },
    { count: 1.5 },
    { latency: -1 },
    { failRate: 2 },
    { fields: [] },
    { fields: [...initialFields(), initialFields()[0]] },
  ]) {
    const state = { ...starter(), ...patch };
    assert.ok(validateBuilder(state).length);
    assert.throws(() => buildConfig(state));
  }
});
test('object responses omit responseType', () => {
  const state = starter();
  state.array = false;
  assert.equal(buildConfig(state).api[0].responseType, undefined);
});
test('public demo config matches the runnable Go example', () => {
  const a = JSON.parse(
    readFileSync(new URL('../public/demo.gobi.json', import.meta.url)),
  );
  const b = JSON.parse(
    readFileSync(
      new URL('../../examples/docs-demo.gobi.json', import.meta.url),
    ),
  );
  assert.deepEqual(a, b);
  assert.equal(b.api.length, 5);
});
