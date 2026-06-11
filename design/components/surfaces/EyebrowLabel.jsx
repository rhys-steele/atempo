import React from 'react';

/**
 * EyebrowLabel — the tracked-out mono section label that titles every panel
 * (WHERE YOU LEFT OFF, NEXT BEAT, OPEN THREADS, LAUNCH RAMP, PARTS · 3).
 */
export function EyebrowLabel({ children, tone = 'muted', style = {} }) {
  const colors = {
    muted: 'var(--text-muted)',
    accent: 'var(--accent)',
    faint: 'var(--text-faint)',
  };
  return (
    <div style={{
      fontFamily: 'var(--font-mono)', fontSize: 'var(--label-size)',
      letterSpacing: 'var(--label-tracking)', textTransform: 'uppercase',
      color: colors[tone], ...style,
    }}>
      {children}
    </div>
  );
}
