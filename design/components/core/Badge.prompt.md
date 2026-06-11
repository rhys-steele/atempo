A small monospace status capsule whose tone always carries meaning (tempo/steady/cold/attention); reach for it for live ports, tempo state, and counts.

```jsx
<Badge tone="tempo" dot live>live :4100</Badge>
<Badge tone="tempo">in tempo · 2h</Badge>
<Badge tone="attention">2 uncommitted</Badge>
<Badge tone="cold">held · 6w</Badge>
```

Tones: `tempo` (teal), `steady` (cream), `cold` (faint), `attention` (amber), `neutral` (hairline). `dot` prepends a status dot; add `live` to make it pulse.
