# Demo a frontend before its backend is ready

These runnable Gobi examples help you exercise a storefront or support-desk interface with synthetic data, slow responses and deliberate failures. They are free under this repository’s MIT license. We maintain Gobi and also offer a separately paid implementation service, described at the end.

## Run a real API locally

From the root of this repository, build the current Go version specified in `go.mod` and launch one configuration:

```sh
go build -o gobi .
./gobi -config examples/demo-sprint/commerce.gobi.json
```

Use a free port; both templates default to 8080. If necessary, change `config.port` before starting. For the alternative example, stop the storefront with Ctrl+C and run:

```sh
./gobi -config examples/demo-sprint/support-desk.gobi.json
```

## Storefront routes

| Request | Behavior |
| --- | --- |
| GET /api/products | Eight generated objects, each containing a `product` object |
| GET /api/products/:id | A freshly generated product wrapper |
| GET /api/customers | Five customer wrappers |
| GET /api/orders | Five orders with nested customer and item data |
| POST /api/orders | Requires the top-level fields `customerId` and `items`; generated acknowledgment |
| GET /api/slow-products | Four product wrappers after a 1.5-second delay |
| GET /api/unavailable | A simulated HTTP 500 on every request |

This free example has seven routes. The separately paid pilot covers up to five **agreed client routes**, not every route in this sample.

```sh
curl http://localhost:8080/api/products
curl http://localhost:8080/api/slow-products
curl -i http://localhost:8080/api/unavailable
curl -i -H 'Content-Type: application/json' \
  -d '{"customerId":"demo","items":[]}' http://localhost:8080/api/orders
# Missing required items produces HTTP 400:
curl -i -H 'Content-Type: application/json' \
  -d '{"customerId":"demo"}' http://localhost:8080/api/orders
```

Successful ordinary mock responses use HTTP 200, including the POST acknowledgment. A payload value `"status":"created"` does not change the HTTP status to 201.

## Reuse shapes with Ref and Array

Open `commerce.gobi.json`. Its `ref.product` object defines a generated ID, company-style name, description, numeric price and boolean stock value. Routes reuse that object with `Ref(ref.product)`. Orders use `Array(len=2,type=Ref(ref.product))` to generate two nested items. A route-level `responseType: "Array(8)"` repeats its response object eight times.

The products route returns **wrappers**, so map `row.product` when rendering. Types such as numeric price and boolean stock remain JSON numbers and booleans rather than display strings.

## Connect a screen through a development proxy

`client.mjs` exports a small fetch client with JSON decoding, an HTTP-status check and a ten-second timeout. Import it in your frontend and call a same-origin proxy path:

```js
import { createClient } from './client.mjs';
const api = createClient('/gobi-api/api');
const rows = await api('/products');
const products = rows.map(row => row.product);
```

Copy the module into your frontend or adjust the import to its location. Your development server must forward `/gobi-api` to `http://127.0.0.1:8080` and strip that prefix. For a Vite-based project, add this proxy block to its existing configuration, preserving its current framework plugins:

```js
server: {
  proxy: {
    '/gobi-api': {
      target: 'http://127.0.0.1:8080',
      rewrite: path => path.replace(/^\/gobi-api/, '')
    }
  }
}
```

Then `/gobi-api/api/products` reaches Gobi’s `/api/products`. Restart the development server after changing its configuration. In containers, use a reachable service hostname rather than assuming loopback is the host. Gobi does not automatically enable cross-origin browser access; the same-origin proxy avoids that problem during development. This is a local setup example, not a production proxy deployment.

## Check the awkward states deliberately

- **Success:** use `/products`; confirm the screen renders the nested `product` fields correctly.
- **Loading:** call `/slow-products`; show a loading state while the promise is pending and clear it in `finally`.
- **Failure:** call `/unavailable`; show a useful error and allow a deliberate retry. This endpoint will keep failing until you change the route or configuration.
- **Empty:** set the frontend’s fixture/data source to an empty list to review its empty state. This sample does not include a dedicated empty-list HTTP route.
- **Validation:** omit `items` from a POST to confirm that the caller handles HTTP 400. Presence validation does not validate each item’s type or business rules.

## Support-desk alternative

The support-desk configuration offers generated tickets, assignees, nested messages and numeric overview fields, plus slow and failing responses. Try `/api/tickets`, `/api/agents` and `/api/overview`. Its POST `/api/tickets` requires `subject` and `email`.

## Limits worth understanding

These examples are stateless: POST does not save a record, and a path ID does not force a subsequent generated object to have that ID. Separate requests can return different values. Order totals are generated independently of line items, so they are not accounting fixtures. There is no real payment processing. Keep all data synthetic. See the repository’s PocketBase guide if the prototype genuinely needs persistence.

## Prefer an implementation handoff?

The [Frontend Demo Sprint](https://demo-sprint.164.68.99.24.sslip.io/) is a separate ₹7,500 pilot for one existing frontend screen, up to five agreed mock API routes, one agreed revision, integration source and run instructions. Scope and timing are confirmed after review. Gobi and these examples remain free MIT software. No customer outcomes or delivery guarantees are implied by this example.
