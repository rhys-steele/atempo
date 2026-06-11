/* global React */

/**
 * KeyCap — a keyboard shortcut shown in a faint hairline cap (e.g. ⌘K).
 * Used for command hints in the command bar / palette affordance. Mono, low-emphasis;
 * never carries status (that's Badge) and never labels a stack (that's Tag).
 */
window.AtempoDS = window.AtempoDS || {};
window.AtempoDS.KeyCap = function KeyCap({ children, style = {} }) {
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
};
