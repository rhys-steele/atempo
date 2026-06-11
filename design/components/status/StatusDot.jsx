import React from 'react';

/**
 * StatusDot — the atomic state mark. teal = live, amber = attention, cream = neutral/cold.
 * `live` turns on the single slow `pulseDot` breath.
 */
export function StatusDot({ tone = 'live', live = false, size = 6, style = {} }) {
  const colors = {
    live: 'var(--accent)',
    attention: 'var(--amber-400)',
    steady: 'var(--cream-300)',
    cold: 'var(--cream-600)',
    neutral: 'var(--text-muted)',
  };
  return (
    <span style={{
      display: 'inline-block',
      width: `${size}px`, height: `${size}px`, borderRadius: '50%',
      background: colors[tone],
      animation: live ? 'pulseDot var(--pulse-duration) ease-in-out infinite' : 'none',
      ...style,
    }} />
  );
}
