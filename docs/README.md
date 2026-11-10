# Gobi docs and interactive demos

Astro + MDX for guides, Svelte 5 for interactive pages, Tailwind CSS 4 for styling. Fonts are self-hosted through Fontsource. No third-party analytics or runtime account is required.

## Develop

Use Node 24 and Go 1.27.3+. Start the backend from the repository root:

```sh
go build -o gobi .
./gobi -config examples/docs-demo.gobi.json
```

Then start the site:

```sh
cd docs
npm ci
npm run dev
```

Visit localhost:4321. The default source is **Sample fixture**. Choose **Live Gobi server** and send/refresh to make real requests through `/gobi-api`. Set `GOBI_API_URL` when starting Astro to override the backend proxy target (default localhost:8080). The static output has no development proxy; deploy `compose.docs.yml` or configure your hosting reverse proxy for live requests.

## Pages

- `/`: introduction and schema/response flow.
- `/docs/quick-start/`: start of the MDX guide, including configuration, routes, generators, Array, Ref, auth/validation, PocketBase and deployment.
- `/playground/`: complete downloadable schemas, sample or live responses, POST body editing, errors and measured request timing.
- `/dashboard/`: data-bound summary figures, an accessible daily-visits chart, searchable customer table and visible source schema.
- `/builder/`: one-endpoint schema builder with editable generator expressions, reusable profile templates, arrays, latency/failure settings, validation and JSON export.

## Check and build

```sh
npm run check
npm test
npm run build
npx playwright install --with-deps chromium
npm run test:e2e
```

To include browser integration tests against a running demo API:

```sh
GOBI_E2E_LIVE=1 GOBI_API_URL=http://127.0.0.1:8080 npm run test:e2e
```

Guide content is in `src/content/*.mdx`; guide navigation is in `src/lib/navigation.ts`. `examples/docs-demo.gobi.json` and `public/demo.gobi.json` must stay in sync (a unit test checks them). The starter builder does not replace the full JSON configuration format.

## Docker

From the repository root:

```sh
docker compose -f compose.docs.yml up --build -d
```

The docs are served on localhost:4321, with live requests proxied to the backend container. No existing Coolify applications are changed by this configuration.

## Dependency audit

The initial audit reports three low-severity entries stemming from one `deepmerge` advisory in the Svelte Vite build-tool dependency chain. Upstream currently resolves to deepmerge 4.3.1, with no compatible patched version. It is used during the build; the static docs do not expose a merge API. Do not apply npm's suggested downgrade to the historical Astro Svelte 0.1.2 integration, which is incompatible with this stack. Recheck with `npm audit` as upstream releases fixes.
