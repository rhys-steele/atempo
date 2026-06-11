import * as React from 'react';

export interface CardProps {
  children: React.ReactNode;
  /** `flat` (hairline), `active` (teal tint, in-tempo), `cold` (receded), `hero` (large radius + tint). */
  emphasis?: 'flat' | 'active' | 'cold' | 'hero';
  radius?: 'lg' | 'xl' | '2xl';
  padding?: number | string;
  style?: React.CSSProperties;
}

/**
 * The base Atempo surface — flat espresso + hairline, no drop shadow. Emphasis encodes
 * tempo: `active` for in-tempo/live, `cold` for receded. Compose rows and panels from this.
 *
 * @startingPoint section="Surfaces" subtitle="Flat espresso surface — flat / active / cold / hero" viewport="700x180"
 */
export function Card(props: CardProps): React.ReactElement;
