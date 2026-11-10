---
name: Gobi
description: A quiet developer workbench for mock APIs.
colors:
  paper: "#fafaf7"
  ink: "#202b25"
  muted: "#59635c"
  line: "#dce1d9"
  green: "#236b48"
  wash: "#eef4ec"
  panel: "#ffffff"
  code-bg: "#182a22"
  code-text: "#e7efdf"
  green-hover: "#195535"
  focus: "#2c8557"
typography:
  display:
    fontFamily: "DM Sans Variable, sans-serif"
    fontSize: "clamp(46px, 4.8vw, 69px)"
    fontWeight: 650
    lineHeight: 1.08
    letterSpacing: "-0.025em"
  headline:
    fontFamily: "DM Sans Variable, sans-serif"
    fontSize: "43px"
    fontWeight: 650
    lineHeight: 1.08
    letterSpacing: "-0.025em"
  title:
    fontFamily: "DM Sans Variable, sans-serif"
    fontSize: "25px"
    fontWeight: 650
    letterSpacing: "-0.025em"
  body:
    fontFamily: "DM Sans Variable, sans-serif"
    fontSize: "15px"
    lineHeight: 1.85
  label:
    fontFamily: "DM Sans Variable, sans-serif"
    fontSize: "12px"
    fontWeight: 650
  code:
    fontFamily: "JetBrains Mono Variable, monospace"
    fontSize: "13px"
    lineHeight: 1.85
rounded:
  badge: "4px"
  tab: "6px"
  control: "7px"
  notice: "9px"
  panel: "12px"
spacing:
  xs: "4px"
  sm: "8px"
  md: "12px"
  lg: "16px"
  xl: "24px"
  2xl: "32px"
  3xl: "48px"
components:
  button-primary:
    backgroundColor: "{colors.green}"
    textColor: "{colors.panel}"
    rounded: "{rounded.control}"
    padding: "10px 16px"
    height: "42px"
  button-primary-hover:
    backgroundColor: "{colors.green-hover}"
  button-secondary:
    backgroundColor: "transparent"
    textColor: "{colors.ink}"
    rounded: "{rounded.control}"
    padding: "10px 16px"
  input:
    backgroundColor: "{colors.panel}"
    textColor: "{colors.ink}"
    rounded: "{rounded.control}"
    padding: "10px 12px"
  panel:
    backgroundColor: "{colors.panel}"
    rounded: "{rounded.panel}"
  code-panel:
    backgroundColor: "{colors.code-bg}"
    textColor: "{colors.code-text}"
    typography: "{typography.code}"
    rounded: "{rounded.panel}"
---

# Design System: Gobi

## Overview

**Creative North Star: "A quiet developer workbench"**

The user approved warm white surfaces, charcoal text and restrained evergreen accents, with code-first implementation. Self-hosted DM Sans supports reading and interface tasks; JetBrains Mono identifies code and measurements. The thesis is to make the transformation from a small JSON template to useful frontend data visible and repeatable.

This direction was explicitly pinned and confirmed before implementation. The user-pinned-direction exemption applies; this is not a retrospective FORM seed. The concept-seed tool ran (key 45686601); its alternatives do not override the confirmed workbench direction. No approved visual comp exists. This document records the implemented docs/demo artifact rather than inventing visual approval.

**Key Characteristics:**
- Warm paper, charcoal reading text and evergreen interaction states.
- Fine separators and flat, readable work areas.
- Visible schema/result relationships and explicit data provenance.
- Reading-first guides and task-focused demos.

## Colors

The palette pairs warm neutral surfaces with a restrained evergreen accent and a dark evergreen code surface.

### Primary
- **Evergreen:** primary actions, links, active navigation and the brand mark; the darker hover token strengthens action feedback.
- **Pale wash:** contextual notices and quiet supporting surfaces.

### Neutral
- **Warm paper:** page canvas.
- **Charcoal ink:** primary text.
- **Muted charcoal:** supporting text and descriptions.
- **Warm line:** panel borders and structural separators.
- **White panel:** fields and work areas.
- **Code evergreen / pale code ink:** dark schema and response panels.

**The Restrained Accent Rule.** Use evergreen for actions, links and active state; let neutral surfaces carry the page.

## Typography

**Display and Body Font:** DM Sans Variable, sans-serif.
**Label/Mono Font:** JetBrains Mono Variable, monospace for code and measurements; interface labels use DM Sans.

The pairing is clear and practical: proportional reading text beside visually distinct machine-readable specimens. Frontmatter records the desktop roles; the homepage display becomes 49px on phones, article titles become 34px, and operation titles move from 42px to 34px. Descriptions use 17px, with the homepage description becoming 16px on phones. Article subheads use the title role and 18px third-level headings. Buttons use 14px, compact controls 12px. Article text uses the body role; supporting interface text ranges from 12–14px.

**The Readable Code Rule.** Nested code inherits its preformatted panel's size. Code panels use the code role; mobile guide and homepage examples use 12px. Do not compound the global inline-code scaling inside a preformatted block.

## Layout

The shared container is capped at 1400px with 40px side padding. Above 1500px padding becomes 48px; at 1150px it becomes 28px, and at 760px it becomes 20px. The masthead is 82px tall and becomes 68px on phones. It remains in normal page flow.

Guides use a 208px sidebar, a central article capped at 760px, a 160px contents rail and 48px gaps. At 1150px the contents rail disappears and the sidebar becomes 190px with 38px gaps. At 760px guides become one column and expose a native navigation disclosure. Article rhythm uses 48px before second-level sections and 16px paragraph spacing.

Operation surfaces pair equal-width work areas with 24px gaps. The builder uses 32px gaps and stacks at 850px, releasing its sticky preview; other split panels stack at 760px. The homepage uses an asymmetric two-column introduction, followed by separated teaching sections. Its schema-flow image switches to a separate vertical SVG at 760px. Dashboard metrics change from four columns to two; chart labels increase in SVG coordinates on phones, and tables retain horizontal scrolling.

The visitor story is guide → schema/result comparison → dashboard → configuration export. First viewports expose purpose and navigation; operation pages identify data sources before controls. These are confirmed surface strategies, not a requirement that every future page copy the homepage composition.

## Elevation & Depth

The system uses flat surfaces, borders and tonal separation rather than ornamental shadows. Code panels establish depth through dark evergreen fill. The homepage specimen alone has a subtle half-degree rotation on desktop, removed on phones and with reduced motion. No shadow token is established.

**The Separator Rule.** Establish panel boundaries with fine lines and surface tone.

## Shapes

Controls have gently curved corners using the control token; large panels use the panel token. Small method badges and inline code use the badge token. Notices, prose code blocks and mobile menus use the notice token. Circular status dots and customer initials are functional silhouettes. SVG icons use geometric strokes rather than text glyphs.

## Components

### Buttons

Compact, clear actions. Primary buttons use evergreen with white text, a matching border and semibold type. Secondary buttons are transparent with a warm gray border, changing to a pale neutral wash on hover. Small variants reduce the minimum height to 34px and padding to 6px 11px. All interactive controls share a 3px green focus outline offset by 4px. Disabled buttons use 0.55 opacity and a waiting cursor. Background transitions last 0.15s.

### Inputs / Fields

White fields use a fine muted border, the control radius and persistent text labels above them. Inputs, selects and textareas share the same treatment and focus outline. Errors use a pale warm red surface and dark rust text; no unobserved per-field invalid border state is established.

### Cards / Containers

White, bordered panels contain a separated header with 15px 20px padding and a body with 24px padding, reduced to 18px on phones. Code panels invert the palette, retain the border-based structure and horizontally scroll long lines. Copy and download controls sit in the header; copy feedback is announced as status text.

### Navigation

The masthead uses muted 14px links with an evergreen active underline and stronger active text. Guide navigation uses softly rounded active washes. On phones native details disclosures expose product and guide links. Search filters guide links and exposes a textual empty result.

### Selection and provenance

Example buttons use a quiet green wash and border when aria-pressed is true. Method badges use compact monospace text and a pale green surface. Status indicators pair a small circle with explicit source text, including the separate sample treatment. Color does not replace the provenance label.

### Data and teaching specimens

Schema/result pairs are the signature pattern. Charts include labeled axes, an accessible description and a table disclosure; customer data is a scrollable table with a search-empty state. Pending, error and empty branches are explicit in the source. Reduced-motion preferences disable animation, transitions and smooth scrolling. These are source-observed states; this document does not claim visual inspection of every branch.

## Do's and Don'ts

### Do:
- **Do** pair schemas with their results and keep sample/live provenance explicit.
- **Do** use the recorded font roles and preserve readable code sizing.
- **Do** retain labeled fields, visible keyboard focus, data alternatives and mobile disclosures.
- **Do** use borders and restrained surface changes to organize work.

### Don't:
- **Don't** introduce gradients, glass or ornamental shadows into the approved workbench.
- **Don't** invent customer proof, performance claims or production-traffic metrics.
- **Don't** hide operation controls when panels stack on phones.
