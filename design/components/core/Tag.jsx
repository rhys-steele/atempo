import React from 'react';

/**
 * Tag — a hairline tech/stack chip (`go`, `astro`, `ts`, `swift`). Mono, low-emphasis.
 * Distinct from Badge: Tag never carries status, only labels a stack or kind.
 */
export function Tag({ children, style = {} }) {
  return (
    <span style={{
      display: 'inline-flex', alignItems: 'center',
      fontFamily: 'var(--font-mono)', fontSize: '11.5px', lineHeight: 1,
      color: 'var(--text-body)',
      border: '1px solid var(--line-12)',
      borderRadius: 'var(--radius-sm)', padding: '4px 10px',
      ...style,
    }}>
      {children}
    </span>
  );
}

/**
 * KeyCap — a keyboard shortcut shown in a faint hairline cap (e.g. ⌘K).
 */
export function KeyCap({ children, style = {} }) {
  return (
    <span style={{
      display: 'inline-flex', alignItems: 'center',
      fontFamily: 'var(--font-mono)', fontSize: '11px', lineHeight: 1,
      color: 'var(--text-faint)',
      border: '1px solid var(--line-12)',
      borderRadius: '5px', padding: '2px 7px',
      ...style,
    }}>
      {children}
    </span>
  );
}
