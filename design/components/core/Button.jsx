import React from 'react';

/**
 * Atempo Button — calm, flat, warm. Three variants:
 *  - primary:  teal fill, espresso text (the one strong action: Run, Resume)
 *  - secondary: hairline outline on the espresso ground
 *  - ghost:    text-only until hover
 * Sizes: sm | md. Use the `glyph` prop for a leading functional mark (▶, →, +).
 */
export function Button({
  children,
  variant = 'secondary',
  size = 'md',
  glyph,
  disabled = false,
  onClick,
  style = {},
}) {
  const pad = size === 'sm' ? '9px 16px' : '11px 22px';
  const fontSize = size === 'sm' ? '13.5px' : '14.5px';

  const variants = {
    primary: {
      background: 'var(--accent)',
      color: 'var(--on-accent)',
      border: '1px solid transparent',
      fontWeight: 'var(--weight-semibold)',
    },
    secondary: {
      background: 'transparent',
      color: 'var(--text-heading)',
      border: '1px solid var(--line-16)',
      fontWeight: 'var(--weight-semibold)',
    },
    ghost: {
      background: 'transparent',
      color: 'var(--text-body)',
      border: '1px solid transparent',
      fontWeight: 'var(--weight-medium)',
    },
    accentOutline: {
      background: 'transparent',
      color: 'var(--accent)',
      border: '1px solid var(--teal-line-40)',
      fontWeight: 'var(--weight-semibold)',
    },
  };

  return (
    <button
      onClick={disabled ? undefined : onClick}
      disabled={disabled}
      style={{
        display: 'inline-flex',
        alignItems: 'center',
        gap: '8px',
        padding: pad,
        borderRadius: 'var(--radius-md)',
        fontFamily: 'var(--font-ui)',
        fontSize,
        lineHeight: 1,
        cursor: disabled ? 'not-allowed' : 'pointer',
        opacity: disabled ? 0.4 : 1,
        transition: 'background 140ms var(--ease-standard), color 140ms var(--ease-standard)',
        ...variants[variant],
        ...style,
      }}
    >
      {glyph && <span style={{ fontFamily: 'var(--font-mono)', fontSize: '0.92em' }}>{glyph}</span>}
      {children}
    </button>
  );
}
