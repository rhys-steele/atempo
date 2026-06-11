/* global React */

/**
 * Tag — a hairline tech/stack chip (`go`, `astro`, `ts`, `swift`). Mono, low-emphasis.
 * Distinct from Badge: Tag never carries status, only labels a stack or kind.
 * (Keyboard-shortcut caps live in the sibling KeyCap component.)
 */
window.AtempoDS = window.AtempoDS || {};
window.AtempoDS.Tag = function Tag({ children, style = {} }) {
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
};
