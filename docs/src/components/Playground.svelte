<script lang="ts">
  import { onMount } from 'svelte';
  let ready = $state(false);
  onMount(() => {
    ready = true;
  });
  import CodePanel from './CodePanel.svelte';
  import {
    exampleCases,
    completeExample,
    sampleResult,
    endpoint,
    fetchJSON,
  } from '../lib/demo';
  let selected = $state(0),
    mode = $state('sample'),
    base = $state('/gobi-api'),
    body = $state('{"name":"page_view"}');
  let data = $state<unknown>(sampleResult('array')),
    busy = $state(false),
    error = $state(''),
    status = $state('Sample fixture'),
    duration = $state<number | null>(null);
  let active = $derived(exampleCases[selected]);
  function choose(index: number) {
    selected = index;
    data = sampleResult(exampleCases[index].id);
    status = 'Sample fixture';
    error = '';
    duration = null;
  }
  async function run() {
    error = '';
    duration = null;
    if (mode === 'sample') {
      data = sampleResult(active.id);
      status = 'Sample fixture';
      return;
    }
    busy = true;
    const start = performance.now();
    try {
      if (active.method === 'POST') {
        JSON.parse(body);
      }
      const result = await fetchJSON(
        endpoint(base, active.path),
        active.method,
        active.method === 'POST' ? body : undefined,
      );
      data = result.data;
      status = `HTTP ${result.status}`;
      duration = Math.round(performance.now() - start);
    } catch (e) {
      error = `${e instanceof Error ? e.message : 'Request failed'}. Check that the Gobi demo server is running and the proxy or CORS is configured.`;
      status = 'Request failed';
      data = null;
    } finally {
      busy = false;
    }
  }
</script>

<div inert={!ready} aria-busy={!ready}>
  <div class="tabs mb-6" aria-label="Example selection">
    {#each exampleCases as item, index}<button
        aria-pressed={selected === index}
        onclick={() => choose(index)}
        disabled={busy}>{item.label}</button
      >{/each}
  </div>
  <div class="panel mb-6">
    <div class="panel-body">
      <div class="toolbar">
        <div>
          <label class="label" for="pg-mode">Response source</label><select
            id="pg-mode"
            bind:value={mode}
            disabled={busy}
            onchange={() => {
              status = 'Run to load selected source';
              data = null;
              error = '';
            }}
            ><option value="sample">Sample fixture</option><option value="live"
              >Live Gobi server</option
            ></select
          >
        </div>
        <div class="grow">
          <label class="label" for="pg-base">Gobi server / proxy URL</label
          ><input
            class="w-full"
            id="pg-base"
            bind:value={base}
            disabled={mode === 'sample' || busy}
          />
        </div>
        <button class="button" onclick={run} disabled={busy}
          >{busy
            ? 'Sending request…'
            : mode === 'sample'
              ? 'Load sample'
              : 'Send request'}</button
        >
      </div>
      <p class="hint mb-0">
        {mode === 'sample'
          ? 'Sample fixtures are illustrative browser data. Switch to Live Gobi server to execute the schema in Go.'
          : 'The bundled demo serves these routes. /gobi-api uses the development or Docker proxy; a different origin requires CORS on your server.'}
      </p>
    </div>
  </div>
  {#if error}<p class="error mb-6" role="alert">{error}</p>{/if}
  <div class="flex justify-between items-center gap-3 mb-4">
    <p class="text-sm text-muted m-0">{active.description}</p>
    <span class="mono text-xs whitespace-nowrap"
      ><span class="method">{active.method}</span> {active.path}</span
    >
  </div>
  {#if active.method === 'POST'}<label class="label" for="pg-body"
      >Request body (JSON)</label
    ><textarea
      class="w-full mono mb-5"
      id="pg-body"
      rows="3"
      bind:value={body}
      disabled={busy}></textarea>{/if}
  <div class="split">
    <CodePanel
      value={completeExample(selected)}
      title="Schema · api.gobi.json"
      download
    />
    <div>
      <div class="flex justify-between mb-3 items-center">
        <span class:sample={mode === 'sample'} class="status" role="status"
          >{status}</span
        >{#if duration !== null}<span class="mono text-xs text-muted"
            >{duration} ms</span
          >{/if}
      </div>
      <CodePanel value={data} title="JSON response" />
      <p class="hint">
        {mode === 'sample'
          ? 'These fixed fixtures demonstrate the response shape. Gobi generates new values on each live request.'
          : 'Response received from the selected server. Timing is measured in your browser.'}
      </p>
    </div>
  </div>
</div>
