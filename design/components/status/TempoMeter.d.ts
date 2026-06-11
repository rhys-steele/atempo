import * as React from 'react';

export interface TempoMeterProps {
  /** Status as pace: `in` (tall teal), `steady` (mid cream), `cold` (short faint). */
  tempo?: 'in' | 'steady' | 'cold';
  /** `lg` doubles as the Atempo logo glyph. */
  size?: 'sm' | 'md' | 'lg';
  style?: React.CSSProperties;
}

/**
 * The signature Atempo mark — a 4-bar equalizer reading tempo as a feeling of pace.
 * Bars descend left-to-right; cold compresses them toward the floor.
 *
 * @startingPoint section="Status" subtitle="The 4-bar tempo equalizer / logo glyph" viewport="700x140"
 */
export function TempoMeter(props: TempoMeterProps): React.ReactElement;
