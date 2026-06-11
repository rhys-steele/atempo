import React from 'react';

/**
 * Card — the base surface. Flat espresso with a hairline; no drop shadow (shadow is
 * reserved for the floating window). `emphasis` raises the border + tints the fill:
 *  - flat:   neutral hairline
 *  - active: teal tint + teal hairline (in tempo / live)
 *  - cold:   faded hairline, recedes
 *  - hero:   large radius for feature/recap surfaces
 */
export function Card({ children, emphasis = 'flat', radius = 'xl', padding = 24, style = {} }) {
  const emph = {
    flat:   { background: 'transparent', border: '1px solid var(--line-08)' },
    active: { background: 'var(--teal-tint-04)', border: '1px solid var(--teal-line-20)' },
    cold:   { background: 'transparent', border: '1px solid var(--line-04)' },
    hero:   { background: 'var(--teal-tint-04)', border: '1px solid var(--teal-line-20)' },
  };
  const radii = {
    lg: 'var(--radius-lg)', xl: 'var(--radius-xl)', '2xl': 'var(--radius-2xl)',
  };
  return (
    <div style={{
      borderRadius: emphasis === 'hero' ? 'var(--radius-2xl)' : radii[radius],
      padding: typeof padding === 'number' ? `${padding}px` : padding,
      boxSizing: 'border-box',
      ...emph[emphasis],
      ...style,
    }}>
      {children}
    </div>
  );
}
