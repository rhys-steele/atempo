Flat, warm action button; use `primary` (teal fill) for the single strong action on a surface and `secondary`/`ghost` for everything else.

```jsx
<Button variant="primary" glyph="▶">Run all — 3 steps</Button>
<Button variant="secondary">Open app</Button>
<Button variant="accentOutline" size="sm" glyph="▸">Run</Button>
<Button variant="ghost" glyph="→">Pick up in code</Button>
```

Variants: `primary` (teal fill, espresso text), `secondary` (hairline outline), `ghost` (text-only until hover), `accentOutline` (teal outline + teal text, for low-emphasis run actions). Sizes `sm | md`. `glyph` takes a Unicode mark — never an emoji or SVG icon.
