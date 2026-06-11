# Archive — direction studies & provenance

These are the in-product direction studies that produced the Atempo design system. The **Cockpit** direction was chosen; its canonical, maintained implementation now lives in `../ui_kits/cockpit/`. Everything here is kept for history, not for active use.

- `The Cockpit.dc.html`, `The Sequencer.dc.html`, `The Studio.dc.html`, `Atempo Directions.dc.html` — the original design-tool exports (`<x-dc>` custom elements + `@dsCard` annotations).
- `support.js` — the generated design-tool runtime. It exists **only** to render the `.dc.html` files above; nothing in the live system depends on it.
- `uploads/` — the original product screenshots provided by the builder (provenance for the espresso/teal direction).
