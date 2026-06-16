/* Atempo site — decorative motion only.
   Two pieces: the equalizer-bar parallax backdrop, and the hero auto-demo that
   cycles a "woken" fleet row. Dependency-free, transform-only, and fully
   disabled under prefers-reduced-motion. Nothing here is load-bearing — the
   page reads and works identically with JS off. */
(() => {
  'use strict';
  const prefersReduced = matchMedia('(prefers-reduced-motion: reduce)').matches;

  /* ---------------------------------------------------------------
     Equalizer-bar parallax backdrop.
     Static (no motion) under reduced-motion — the bars still add depth,
     they just don't move.
  --------------------------------------------------------------- */
  function initBackdrop() {
    const field = document.createElement('div');
    field.className = 'bg-eq';
    field.setAttribute('aria-hidden', 'true');

    const count = Math.max(5, Math.min(18, Math.round(innerWidth / 70)));
    const bars = [];
    for (let i = 0; i < count; i++) {
      const bar = document.createElement('span');
      bar.className = 'bg-eq__bar' + (i % 3 === 0 ? ' bg-eq__bar--cream' : '');
      const x = (i + 0.5) / count;                       // 0..1 across the width
      bar.style.left = (x * 100).toFixed(2) + '%';
      bar.style.height = (28 + 56 * Math.abs(Math.sin(i * 1.7))).toFixed(1) + 'vh';
      bars.push({ el: bar, x, phase: i * 0.9 });
      field.appendChild(bar);
    }
    document.body.prepend(field);

    if (prefersReduced) return;   // depth without motion

    let mx = 0.5, my = 0.5, cmx = 0.5, cmy = 0.5;   // target + eased cursor
    let running = true;

    addEventListener('pointermove', (e) => {
      mx = e.clientX / innerWidth;
      my = e.clientY / innerHeight;
    }, { passive: true });
    document.addEventListener('visibilitychange', () => {
      running = !document.hidden;
      if (running) requestAnimationFrame(frame);
    });

    function frame() {
      if (!running) return;
      cmx += (mx - cmx) * 0.06;                       // ease toward cursor — calm follow
      cmy += (my - cmy) * 0.06;
      const tiltX = (cmx - 0.5) * 16;                 // gentle horizontal parallax
      const driftY = -window.scrollY * 0.04 + (cmy - 0.5) * 10;  // scroll drift + subtle vertical tilt
      field.style.transform = `translate3d(${tiltX.toFixed(2)}px, ${driftY.toFixed(2)}px, 0)`;

      const t = performance.now() / 1000;
      for (const b of bars) {
        const near = Math.max(0, 1 - Math.abs(b.x - cmx) * 3);   // proximity to cursor
        const breathe = 0.05 * Math.sin(t * 0.6 + b.phase);      // slow equalizer breath
        const sy = 1 + breathe + near * 0.18;
        b.el.style.transform = `scaleY(${sy.toFixed(3)})`;
      }
      requestAnimationFrame(frame);
    }
    requestAnimationFrame(frame);
  }

  /* ---------------------------------------------------------------
     Hero wake-a-project auto-demo.
     The wake itself is pure CSS (:hover/:focus-within). This just cycles
     a woken row so the demo is alive on touch and before any interaction.
     Pauses while the pointer is over the panel — hover wins.
  --------------------------------------------------------------- */
  function initFleetDemo() {
    if (prefersReduced) return;
    const demo = document.querySelector('.hero__demo');
    if (!demo) return;
    const rows = Array.from(demo.querySelectorAll('.fleet-row'));
    if (rows.length < 2) return;

    let i = Math.max(0, rows.findIndex((r) => r.classList.contains('is-awake')));
    let paused = false;
    // Hand off to :hover while the pointer is in the panel — clear the auto-woken
    // row so only the hovered one is awake; resume cycling on leave.
    demo.addEventListener('pointerenter', () => {
      paused = true;
      rows.forEach((r) => r.classList.remove('is-awake'));
    });
    demo.addEventListener('pointerleave', () => {
      paused = false;
      rows[i].classList.add('is-awake');
    });

    setInterval(() => {
      if (paused) return;
      rows[i].classList.remove('is-awake');
      i = (i + 1) % rows.length;
      rows[i].classList.add('is-awake');
    }, 2800);
  }

  function start() { initBackdrop(); initFleetDemo(); }
  if (document.readyState === 'loading') addEventListener('DOMContentLoaded', start);
  else start();
})();
