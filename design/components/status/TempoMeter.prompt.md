The signature Atempo mark — a 4-bar equalizer that reads tempo as pace; also the logo glyph at `size="lg"`.

```jsx
<TempoMeter tempo="in" />
<TempoMeter tempo="steady" />
<TempoMeter tempo="cold" size="sm" />
<TempoMeter tempo="in" size="lg" />  {/* logo lockup: place beside "atempo" */}
```

`tempo`: `in` (tall, full teal) / `steady` (mid, cream, fading) / `cold` (short, faint brown). Never animate it — tempo changes by decay, not motion.
