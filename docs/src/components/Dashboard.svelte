<script lang="ts">
  import { onMount } from 'svelte';
  let ready = $state(false);
  onMount(() => {
    ready = true;
  });
  import { sampleDashboard, endpoint, fetchJSON } from '../lib/demo';
  import CodePanel from './CodePanel.svelte';
  type Dashboard = typeof sampleDashboard;
  let data = $state<Dashboard | null>(structuredClone(sampleDashboard)),
    mode = $state('sample'),
    source = $state('Sample fixture'),
    base = $state('/gobi-api'),
    busy = $state(false),
    error = $state(''),
    search = $state(''),
    showSchema = $state(false);
  const money = new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'USD',
    maximumFractionDigits: 0,
  });
  let customers = $derived(
    data?.customers.filter((c) =>
      `${c.name} ${c.email}`.toLowerCase().includes(search.toLowerCase()),
    ) ?? [],
  );
  let maximum = $derived(
    Math.max(100, ...(data?.series.map((p) => p.visits) ?? [])),
  );
  let points = $derived(
    data?.series
      .map(
        (p, i) =>
          `${60 + i * (620 / Math.max(1, (data?.series.length ?? 1) - 1))},${220 - (p.visits / maximum) * 175}`,
      )
      .join(' ') ?? '',
  );
  const schema = {
    method: 'GET',
    path: '/dashboard',
    response: {
      visitors: 'Int(min=3000,max=9000)',
      revenue: 'Int(min=12000,max=30000)',
      orders: 'Int(min=100,max=500)',
      conversion: 'Float(min=1,max=8,precision=2)',
      series: ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'].map((day) => ({
        day,
        visits: 'Int(min=200,max=1200)',
      })),
      customers: 'Array(len=5,type=Ref(ref.customer))',
      updated: 'Time(iso)',
    },
  };
  function valid(value: unknown): value is Dashboard {
    if (!value || typeof value !== 'object') return false;
    const v = value as Record<string, unknown>;
    return (
      ['visitors', 'revenue', 'orders', 'conversion'].every(
        (k) =>
          typeof v[k] === 'number' &&
          Number.isFinite(v[k]) &&
          Number(v[k]) >= 0,
      ) &&
      Array.isArray(v.series) &&
      v.series.length > 0 &&
      v.series.length <= 31 &&
      v.series.every(
        (p) =>
          p &&
          typeof p.day === 'string' &&
          typeof p.visits === 'number' &&
          Number.isFinite(p.visits) &&
          p.visits >= 0,
      ) &&
      Array.isArray(v.customers) &&
      v.customers.length <= 50 &&
      v.customers.every(
        (c) =>
          c &&
          typeof c.id === 'string' &&
          typeof c.name === 'string' &&
          typeof c.email === 'string' &&
          typeof c.orders === 'number' &&
          Number.isFinite(c.orders) &&
          typeof c.spend === 'number' &&
          Number.isFinite(c.spend),
      )
    );
  }
  async function refresh() {
    busy = true;
    error = '';
    try {
      if (mode === 'sample') {
        data = structuredClone(sampleDashboard);
        source = 'Sample fixture';
      } else {
        const result = await fetchJSON(endpoint(base, '/api/dashboard'));
        if (!valid(result.data))
          throw new Error('The response does not match the dashboard schema');
        data = result.data;
        source = 'Live Gobi response';
      }
    } catch (e) {
      data = null;
      source = 'Request failed';
      error = `${e instanceof Error ? e.message : 'Unable to load dashboard'}. Start the bundled demo API, check the URL, then retry.`;
    } finally {
      busy = false;
    }
  }
</script>

<div inert={!ready} aria-busy={!ready}>
  <div class="panel mb-7">
    <div class="panel-body">
      <div class="toolbar">
        <div>
          <label for="dash-mode" class="label">Data source</label><select
            id="dash-mode"
            bind:value={mode}
            disabled={busy}
            ><option value="sample">Sample fixture</option><option value="live"
              >Live Gobi server</option
            ></select
          >
        </div>
        <div class="grow">
          <label class="label" for="dash-base">Gobi server / proxy URL</label
          ><input
            id="dash-base"
            class="w-full"
            bind:value={base}
            disabled={mode === 'sample' || busy}
          />
        </div>
        <button class="button" onclick={refresh} disabled={busy}
          >{busy ? 'Fetching data…' : 'Refresh data'}</button
        >
      </div>
      <p class="hint mb-0">
        {source === 'Live Gobi response'
          ? 'Every value below comes from GET /api/dashboard. These are mock business figures, not production analytics.'
          : 'The initial view uses a fixed sample. Choose Live Gobi server and refresh to fetch freshly generated mock data.'}
      </p>
    </div>
  </div>
  {#if error}<p class="error mb-6" role="alert">{error}</p>{/if}
  <div class="flex flex-wrap items-center justify-between gap-3 mb-5">
    <div class="flex items-center gap-3">
      <h2 class="text-xl font-semibold m-0">Store overview</h2>
      <span
        class:sample={source !== 'Live Gobi response'}
        class="status"
        role="status">{source}</span
      >
    </div>
    <button
      class="button secondary small"
      onclick={() => (showSchema = !showSchema)}
      aria-expanded={showSchema}
      >{showSchema ? 'Hide schema' : 'View the schema'}</button
    >
  </div>
  {#if showSchema}<div class="mb-7">
      <CodePanel
        title="Dashboard route · Array + Ref"
        value={{
          ref: {
            customer: {
              id: 'Uuid()',
              name: 'User()',
              email: 'User(email)',
              orders: 'Int(min=1,max=30)',
              spend: 'Int(min=80,max=900)',
            },
          },
          api: [schema],
        }}
      />
    </div>{/if}
  {#if busy}<div class="panel p-10 mb-6" role="status">
      Fetching a fresh dashboard from Gobi…
    </div>{:else if data}
    <div class="stats">
      <div>
        <span>Visitors</span><strong>{data.visitors.toLocaleString()}</strong
        ><small>Int(min=3000,max=9000)</small>
      </div>
      <div>
        <span>Revenue</span><strong>{money.format(data.revenue)}</strong><small
          >Generated amount · USD</small
        >
      </div>
      <div>
        <span>Orders</span><strong>{data.orders.toLocaleString()}</strong><small
          >Integer response</small
        >
      </div>
      <div>
        <span>Conversion</span><strong>{data.conversion}%</strong><small
          >Decimal response</small
        >
      </div>
    </div>
    <div class="panel mt-6">
      <div class="panel-head">
        <span>Visits by day</span><span class="text-muted font-normal"
          >{data.series.length} generated days</span
        >
      </div>
      <div class="chart">
        <svg
          viewBox="0 0 740 270"
          role="img"
          aria-labelledby="chart-title chart-desc"
          ><title id="chart-title">Generated daily visits</title><desc
            id="chart-desc"
            >{data.series
              .map((p) => `${p.day}: ${p.visits} visits`)
              .join('. ')}. This chart shows mock data.</desc
          >{#each [0, 0.5, 1] as ratio}<line
              x1="60"
              y1={220 - ratio * 175}
              x2="680"
              y2={220 - ratio * 175}
              class="chart-grid"
            /><text
              x="50"
              y={225 - ratio * 175}
              text-anchor="end"
              fill="#657061"
              font-size="11">{Math.round(maximum * ratio)}</text
            >{/each}<polygon
            points={`60,220 ${points} 680,220`}
            fill="#edf4e8"
          /><polyline
            {points}
            fill="none"
            stroke="#35764d"
            stroke-width="2.5"
            stroke-linecap="round"
            stroke-linejoin="round"
          />{#each data.series as point, i}<circle
              cx={60 + i * (620 / Math.max(1, (data?.series.length ?? 1) - 1))}
              cy={220 - (point.visits / maximum) * 175}
              r="3.5"
              fill="#35764d"
            /><text
              x={60 + i * (620 / Math.max(1, (data?.series.length ?? 1) - 1))}
              y="248"
              text-anchor="middle"
              fill="#657061"
              font-size="11">{point.day}</text
            >{/each}</svg
        >
      </div>
      <details class="px-6 pb-4 text-xs text-muted">
        <summary class="cursor-pointer">View chart data as a table</summary>
        <table class="data-table mt-3">
          <thead><tr><th>Day</th><th>Visits</th></tr></thead><tbody
            >{#each data.series as point}<tr
                ><td>{point.day}</td><td>{point.visits}</td></tr
              >{/each}</tbody
          >
        </table>
      </details>
    </div>
    <div class="panel mt-6">
      <div class="panel-head">
        <span
          >Customers <span class="text-muted font-normal"
            >Array(len=5,type=Ref(ref.customer))</span
          ></span
        ><label
          ><span class="sr-only">Search customers</span><input
            class="text-xs max-w-44"
            placeholder="Search customers…"
            bind:value={search}
          /></label
        >
      </div>
      <div class="table-scroll">
        <table class="data-table">
          <thead
            ><tr
              ><th>Customer</th><th>Email</th><th>Orders</th><th>Spend</th></tr
            ></thead
          ><tbody
            >{#each customers as customer}<tr
                ><td
                  ><span class="customer-name"
                    ><span class="avatar" aria-hidden="true"
                      >{customer.name
                        .split(' ')
                        .map((p) => p[0])
                        .slice(0, 2)
                        .join('')}</span
                    >{customer.name}</span
                  ></td
                ><td>{customer.email}</td><td class="tabular"
                  >{customer.orders}</td
                ><td class="tabular">{money.format(customer.spend)}</td></tr
              >{/each}</tbody
          >
        </table>
        {#if !customers.length}<p class="empty">
            No customers match “{search}”. Clear the search to see all records.
          </p>{/if}
      </div>
    </div>
    <p class="hint mt-4">
      Refresh generates an independent dataset; the summary totals are not
      calculated from the customer rows.
    </p>
  {:else}<div class="panel empty">
      No dashboard data. Check the connection above and choose Refresh data.
    </div>{/if}
</div>

<style>
  .stats {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    border: 1px solid var(--color-line);
    border-radius: 12px;
    background: white;
  }
  .stats > div {
    padding: 22px 26px;
    border-right: 1px solid var(--color-line);
  }
  .stats > div:last-child {
    border: 0;
  }
  .stats span {
    font-size: 12px;
    color: #59635c;
  }
  .stats strong {
    display: block;
    font-size: 30px;
    font-weight: 600;
    letter-spacing: -0.025em;
    margin: 7px 0;
    font-variant-numeric: tabular-nums;
  }
  .stats small {
    font-size: 10px;
    color: #687464;
  }
  .chart {
    padding: 20px 20px 0;
  }
  .chart svg {
    width: 100%;
    height: auto;
    max-height: 310px;
  }
  .chart svg text {
    font-size: 13px;
  }
  .data-table {
    width: 100%;
    border-collapse: collapse;
    text-align: left;
    font-size: 12px;
  }
  .data-table th {
    font-weight: 500;
    color: #687464;
    background: #f8faf4;
    padding: 13px 24px;
  }
  .data-table td {
    padding: 16px 24px;
    border-top: 1px solid #e9eee4;
  }
  .customer-name {
    display: flex;
    align-items: center;
    gap: 11px;
    white-space: nowrap;
    font-weight: 550;
  }
  .avatar {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 30px;
    height: 30px;
    border-radius: 50%;
    background: #edf2e6;
    color: #527446;
    font-size: 10px;
  }
  .table-scroll {
    overflow-x: auto;
  }
  @media (max-width: 760px) {
    .stats {
      grid-template-columns: 1fr 1fr;
    }
    .stats > div {
      padding: 16px 18px;
    }
    .stats > div:nth-child(2) {
      border-right: 0;
    }
    .stats > div:nth-child(-n + 2) {
      border-bottom: 1px solid var(--color-line);
    }
    .stats strong {
      font-size: 26px;
    }
    .stats small {
      font-size: 9px;
    }
    .panel-head {
      flex-wrap: wrap;
    }
    .panel-head > span > span {
      display: none;
    }
    .data-table td,
    .data-table th {
      padding: 14px 18px;
    }
    .chart {
      padding: 12px 4px 0;
    }
    .chart svg text {
      font-size: 25px;
    }
    .chart svg {
      overflow: visible;
    }
  }
</style>
