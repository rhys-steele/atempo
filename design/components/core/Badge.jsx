/* global React */

/**
 * Badge — a small status capsule. Tones map to Atempo's meaning system:
 *  - tempo:     teal tint (in tempo / live)
 *  - steady:    cream tint
 *  - cold:      faint, inert
 *  - attention: amber (open / uncommitted / note)
 *  - neutral:   hairline only
 * Pass `dot` to prepend a status dot; `live` makes the dot breathe.
 */
window.AtempoDS = window.AtempoDS || {};
window.AtempoDS.Badge = function Badge({ children, tone = 'neutral', dot = false, live = false, style = {} }) {
  const tones = {
    tempo:     { color: 'var(--accent)', bg: 'var(--teal-tint-10)', border: 'transparent', dotColor: 'var(--accent)' },
    steady:    { color: 'var(--cream-300)', bg: 'rgba(201,185,163,0.10)', border: 'transparent', dotColor: 'var(--cream-300)' },
    cold:      { color: 'var(--cream-600)', bg: 'rgba(242,233,221,0.04)', border: 'transparent', dotColor: 'var(--cream-600)' },
    attention: { color: 'var(--amber-400)', bg: 'rgba(224,164,88,0.10)', border: 'transparent', dotColor: 'var(--amber-400)' },
    neutral:   { color: 'var(--text-heading)', bg: 'transparent', border: 'var(--line-12)', dotColor: 'var(--text-muted)' },
  };
  const t = tones[tone];
  return (
    <span style={{
      display: 'inline-flex', alignItems: 'center', gap: '7px',
      fontFamily: 'var(--font-mono)', fontSize: '11.5px', lineHeight: 1,
      color: t.color, background: t.bg,
      border: `1px solid ${t.border}`,
      borderRadius: 'var(--radius-sm)', padding: '5px 10px',
      ...style,
    }}>
      {dot && <span style={{
        width: 'var(--dot-size-sm)', height: 'var(--dot-size-sm)', borderRadius: '50%',
        background: t.dotColor,
        animation: live ? 'pulseDot var(--pulse-duration) ease-in-out infinite' : 'none',
      }} />}
      {children}
    </span>
  );
};
