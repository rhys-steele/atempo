/* global React */

/**
 * TempoMeter — the signature mark. A 4-bar equalizer reading status as pace.
 *  - in:     tall, full teal
 *  - steady: mid, cream, last bars fading
 *  - cold:   short, faint brown
 * Also the Atempo logo glyph (use size="lg"). Never animate it — pace is read at a glance.
 */
window.AtempoDS = window.AtempoDS || {};
window.AtempoDS.TempoMeter = function TempoMeter({ tempo = 'in', size = 'md', style = {} }) {
  const unit = size === 'lg' ? { w: 8, gap: 4, h: [34, 24, 30, 18] }
            : size === 'sm' ? { w: 5, gap: 3, h: [18, 14, 16, 11] }
            : { w: 6, gap: 3, h: [22, 17, 20, 13] };

  const palette = {
    in:     ['var(--teal-500)', 'var(--teal-500)', 'var(--teal-500)', 'var(--teal-500)'],
    steady: ['var(--cream-300)', 'var(--cream-300)', 'rgba(201,185,163,0.4)', 'rgba(201,185,163,0.25)'],
    cold:   ['rgba(124,109,92,0.6)', 'rgba(124,109,92,0.4)', 'rgba(124,109,92,0.25)', 'rgba(124,109,92,0.15)'],
  };

  // cold compresses every bar toward the floor
  const heights = tempo === 'cold'
    ? unit.h.map((h, i) => (i === 0 ? Math.round(h * 0.5) : 9))
    : tempo === 'steady'
      ? unit.h.map((h, i) => (i < 2 ? h : 9))
      : unit.h;

  return (
    <span style={{
      display: 'inline-flex', alignItems: 'flex-end', gap: `${unit.gap}px`,
      height: `${Math.max(...unit.h)}px`, ...style,
    }}>
      {heights.map((h, i) => (
        <span key={i} style={{
          width: `${unit.w}px`, height: `${h}px`,
          borderRadius: '2px', background: palette[tempo][i],
        }} />
      ))}
    </span>
  );
};
