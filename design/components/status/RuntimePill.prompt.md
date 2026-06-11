Read-only runtime state shown as a result — live ports, ready, bypassed, or none — with no start/stop controls (Atempo runs things for you).

```jsx
<RuntimePill state="live" label="live :4100 :1420" />
<RuntimePill state="ready" label="ready — 1 service" />
<RuntimePill state="bypassed" label="static · skipped" />
<RuntimePill state="none" />
```

Critical: this never gains buttons. It reports; it does not operate. Use `bypassed` (dashed) for a part with no runtime, like a static site.
