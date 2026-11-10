<script lang="ts">
  import { onMount } from 'svelte';
  let ready = $state(false);
  onMount(() => {
    ready = true;
  });
  import CodePanel from './CodePanel.svelte';
  import {
    generators,
    initialFields,
    validateBuilder,
    buildConfig,
  } from '../lib/builder.mjs';
  import { downloadJSON } from '../lib/demo';
  let method = $state('GET'),
    path = $state('/users'),
    port = $state(8080),
    prefix = $state('/api'),
    array = $state(true),
    count = $state(6),
    latency = $state(0),
    failRate = $state(0),
    fields = $state(initialFields()),
    nextId = 5,
    notice = $state('');
  let draft = $derived({
    method,
    path,
    port,
    prefix,
    array,
    count,
    latency,
    failRate,
    fields,
  });
  let errors = $derived(validateBuilder(draft));
  let config = $derived(errors.length ? null : buildConfig(draft));
  function add() {
    fields.push({
      id: nextId++,
      name: `field_${nextId - 1}`,
      kind: 'name',
      value: 'User()',
    });
  }
  function change(field: ReturnType<typeof initialFields>[number]) {
    field.value = generators.find((g) => g[0] === field.kind)?.[2] ?? '';
    notice = '';
  }
  function reset() {
    method = 'GET';
    path = '/users';
    port = 8080;
    prefix = '/api';
    array = true;
    count = 6;
    latency = 0;
    failRate = 0;
    fields = initialFields();
    nextId = 5;
    notice = 'Starter schema restored.';
  }
  function download() {
    if (config) {
      downloadJSON(config, 'api.gobi.json');
      notice =
        'Downloaded api.gobi.json. Run it with gobi -config api.gobi.json.';
    }
  }
</script>

<div inert={!ready} aria-busy={!ready}>
  <div class="builder-grid">
    <div>
      <div class="panel">
        <div class="panel-head">
          <span>Endpoint</span><span class="mono text-xs text-muted"
            >{method} {prefix}{path}</span
          >
        </div>
        <div class="panel-body">
          <div class="grid grid-cols-[110px_1fr] gap-4 mb-5">
            <div>
              <label class="label" for="method">Method</label><select
                class="w-full"
                id="method"
                bind:value={method}
                >{#each ['GET', 'POST', 'PUT', 'PATCH', 'DELETE'] as m}<option
                    >{m}</option
                  >{/each}</select
              >
            </div>
            <div>
              <label class="label" for="route-path">Route path</label><input
                class="w-full mono"
                id="route-path"
                bind:value={path}
                placeholder="/users"
              />
            </div>
          </div>
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="label" for="prefix">API prefix</label><input
                id="prefix"
                class="w-full mono"
                bind:value={prefix}
              />
            </div>
            <div>
              <label class="label" for="port">Port</label><input
                id="port"
                class="w-full"
                type="number"
                min="1"
                max="65535"
                bind:value={port}
              />
            </div>
          </div>
        </div>
      </div>
      <div class="panel mt-5">
        <div class="panel-head">
          <span>Response fields</span><button
            class="button secondary small"
            onclick={add}>Add field</button
          >
        </div>
        <div class="fields">
          {#each fields as field, index (field.id)}<div class="field-row">
              <div class="grid grid-cols-2 gap-3">
                <div>
                  <label class="label" for={`name-${field.id}`}
                    >Field name</label
                  ><input
                    class="w-full mono"
                    id={`name-${field.id}`}
                    bind:value={field.name}
                  />
                </div>
                <div>
                  <label class="label" for={`kind-${field.id}`}
                    >Value type</label
                  ><select
                    class="w-full"
                    id={`kind-${field.id}`}
                    bind:value={field.kind}
                    onchange={() => change(field)}
                    >{#each generators as [key, label]}<option value={key}
                        >{label}</option
                      >{/each}</select
                  >
                </div>
              </div>
              <div class="flex gap-3 items-end mt-3">
                <div class="grow min-w-0">
                  <label class="label" for={`value-${field.id}`}
                    >{field.kind === 'literal'
                      ? 'Literal value'
                      : 'Generator expression'}</label
                  ><input
                    class="w-full mono"
                    id={`value-${field.id}`}
                    bind:value={field.value}
                  />
                </div>
                <button
                  class="remove"
                  onclick={() =>
                    (fields = fields.filter((f) => f.id !== field.id))}
                  aria-label={`Remove field ${index + 1}`}
                  title="Remove field"
                  ><svg
                    width="17"
                    height="17"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="1.7"
                    aria-hidden="true"
                    ><path
                      d="M4 7h16M9 7V4h6v3M7 7l1 14h8l1-14M10 10v7m4-7v7"
                    /></svg
                  ></button
                >
              </div>
            </div>{/each}{#if !fields.length}<p class="empty">
              Add a field to start your response template.
            </p>{/if}
        </div>
      </div>
      <div class="panel mt-5">
        <div class="panel-head"><span>Response behavior</span></div>
        <div class="panel-body">
          <div class="flex items-center justify-between gap-4 mb-5">
            <label class="flex items-center gap-3 text-sm" for="as-array"
              ><input
                type="checkbox"
                id="as-array"
                bind:checked={array}
              />Return an array</label
            >{#if array}<div class="w-28">
                <label class="sr-only" for="array-count"
                  >Response array length</label
                ><input
                  class="w-full"
                  id="array-count"
                  type="number"
                  min="0"
                  max="10000"
                  bind:value={count}
                />
              </div>{/if}
          </div>
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="label" for="latency">Latency (ms)</label><input
                id="latency"
                class="w-full"
                type="number"
                min="0"
                max="30000"
                bind:value={latency}
              />
            </div>
            <div>
              <label class="label" for="fail-rate">Failure rate (0–1)</label
              ><input
                id="fail-rate"
                class="w-full"
                type="number"
                min="0"
                max="1"
                step="0.05"
                bind:value={failRate}
              />
            </div>
          </div>
        </div>
      </div>
    </div>
    <div class="builder-preview">
      <div class="flex justify-between items-center mb-4 gap-3">
        <h2 class="text-base font-semibold m-0">Your configuration</h2>
        <button class="button secondary small" onclick={reset}
          >Reset starter</button
        >
      </div>
      {#if errors.length}<div class="error mb-4" role="alert">
          <p class="m-0 font-semibold">Complete these fields to export</p>
          <ul class="list-disc pl-5 mt-2 text-xs">
            {#each errors as error}<li>{error}</li>{/each}
          </ul>
        </div>{/if}<CodePanel
        value={config}
        title={config ? 'api.gobi.json' : 'Configuration incomplete'}
      /><button class="button w-full mt-4" disabled={!config} onclick={download}
        >Download api.gobi.json <svg
          width="17"
          height="17"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="1.7"
          aria-hidden="true"
          ><path d="M12 3v12m-5-5 5 5 5-5M5 16v5h14v-5" /></svg
        ></button
      >{#if notice}<p class="notice mt-4" role="status">{notice}</p>{/if}
      <p class="hint mt-4">
        A complete Gobi configuration, generated as you edit. Run the downloaded
        file to execute generators in Go. Reference fields automatically include
        a reusable <code>profile</code> template.
      </p>
      <div class="notice mt-5">
        <strong class="block mb-1">Take it for a spin</strong><code
          class="text-xs break-all">gobi -config api.gobi.json</code
        >
        <p class="m-0 mt-2">
          Then request <code class="break-all"
            >http://localhost:{port}{prefix}{path}</code
          >
          using {method}.
        </p>
      </div>
      <p class="hint mt-4">
        The starter builder supports one endpoint and flat fields. Nested
        objects, multiple routes, authentication, and PocketBase can be added in
        JSON. Generator expressions are editable; Gobi handles their
        interpretation.
      </p>
    </div>
  </div>
</div>

<style>
  .builder-grid {
    display: grid;
    grid-template-columns: minmax(0, 1fr) minmax(0, 1fr);
    gap: 32px;
    align-items: start;
  }
  .builder-preview {
    position: sticky;
    top: 24px;
    min-width: 0;
  }
  .field-row {
    padding: 20px 24px;
    border-bottom: 1px solid var(--color-line);
  }
  .field-row:last-child {
    border-bottom: 0;
  }
  .remove {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 38px;
    height: 40px;
    border: 1px solid #d7dfd2;
    border-radius: 7px;
    background: #fff;
    color: #6a7665;
    flex-shrink: 0;
  }
  .remove:hover {
    background: #fff2ed;
    color: #8a432b;
    border-color: #d6a892;
  }
  @media (max-width: 850px) {
    .builder-grid {
      grid-template-columns: 1fr;
    }
    .builder-preview {
      position: static;
    }
    .field-row {
      padding: 18px;
    }
  }
</style>
