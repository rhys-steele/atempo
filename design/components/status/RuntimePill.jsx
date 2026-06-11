import React from 'react';
import { StatusDot } from './StatusDot.jsx';

/**
 * RuntimePill — shows a part's runtime state as a result, never a control.
 *  - live:     teal pill with breathing dot + ports
 *  - ready:    faint mono "ready — N service(s)"
 *  - bypassed: dashed inert pill ("static · skipped")
 *  - none:     faint "no runtime"
 * This is read-only by design — Atempo runs things for you; you don't operate them here.
 */
export function RuntimePill({ state = 'ready', label, style = {} }) {
  if (state === 'live') {
    return (
      <span style={{
        display: 'inline-flex', alignItems: 'center', gap: '7px',
        fontFamily: 'var(--font-mono)', fontSize: '11.5px', color: 'var(--accent)',
        background: 'var(--teal-tint-10)', borderRadius: 'var(--radius-sm)', padding: '5px 10px',
        ...style,
      }}>
        <StatusDot tone="live" live size={5} />
        {label || 'live'}
      </span>
    );
  }
  if (state === 'bypassed') {
    return (
      <span style={{
        display: 'inline-flex', alignItems: 'center',
        fontFamily: 'var(--font-mono)', fontSize: '11.5px', color: 'var(--text-faint)',
        border: '1px dashed var(--line-16)', borderRadius: 'var(--radius-sm)', padding: '5px 10px',
        ...style,
      }}>
        {label || 'static · skipped'}
      </span>
    );
  }
  return (
    <span style={{
      fontFamily: 'var(--font-mono)', fontSize: '11.5px',
      color: state === 'none' ? 'var(--text-faint)' : 'var(--text-muted)', ...style,
    }}>
      {label || (state === 'none' ? 'no runtime' : 'ready')}
    </span>
  );
}
