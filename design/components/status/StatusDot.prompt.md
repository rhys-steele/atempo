The atomic status mark — a small round dot, teal when live (optionally pulsing), amber for attention, cream for steady/cold.

```jsx
<StatusDot tone="live" live />
<StatusDot tone="attention" />
```

Pair it inline with mono runtime text (`<StatusDot tone="live" live /> live :4100`). Use `live` only for genuinely-running services.
