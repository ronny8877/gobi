<script lang="ts">
  let {
    value,
    title = 'api.gobi.json',
    download = false,
  }: { value: unknown; title?: string; download?: boolean } = $props();
  import { downloadJSON } from '../lib/demo';
  let notice = $state('');
  async function copy() {
    try {
      await navigator.clipboard.writeText(JSON.stringify(value, null, 2));
      notice = 'Copied';
    } catch {
      notice = 'Copy unavailable. Select the JSON to copy it.';
    }
  }
</script>

<div class="panel code-panel">
  <div class="panel-head">
    <span>{title}</span>
    <div class="flex gap-2">
      <button class="button secondary small" onclick={copy}>Copy</button
      >{#if download}<button
          class="button secondary small"
          onclick={() => downloadJSON(value, 'api.gobi.json')}>Download</button
        >{/if}
    </div>
  </div>
  <pre><code>{JSON.stringify(value, null, 2)}</code></pre>
  {#if notice}<p class="m-0 px-6 pb-4 text-xs text-[#d4e5d4]" role="status">
      {notice}
    </p>{/if}
</div>
