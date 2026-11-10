# Product
<!-- impeccable:product-schema 1 -->

## Platform
web

## Stack
Delegated by the user: Astro for Markdown/MDX documentation, Svelte 5 for interactive demos and schema building, Tailwind CSS 4. The Go application remains the mock API engine.

## Users
Developers learning and evaluating Gobi, integrating mock API data into frontends, and creating .gobi.json configurations. Audience and goal confirmed by the user.

## Product Purpose
Make Gobi understandable and directly usable through clear documentation, visual examples, a working dashboard demo, and a starter schema builder that exports a valid configuration.

## Capabilities and Constraints
Explain actual Go generator syntax, HTTP methods, Array and Ref, auth, latency/failures, hot reload and optional PocketBase persistence. Interactive examples must distinguish illustrative browser fixtures from live Gobi responses. The builder exports configuration; it does not remotely execute arbitrary configurations. No invented customer or performance claims.

## Brand Commitments
Gobi; free MIT software by ronny8877. User approved a light developer workbench: warm white surfaces, charcoal text, restrained green accents, clean UI. User explicitly chose code-first implementation.

## Evidence on Hand
Go source, README, example configurations, API behavior and test suite in this repository. Existing images show the terminal CLI, not the new documentation site.

## Product Principles
- Show the schema alongside its result.
- Teach real supported features with runnable examples.
- Make first steps possible without a backend, clearly label sample data.
- Export a complete usable configuration.
