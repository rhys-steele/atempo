import * as React from 'react';

export interface BadgeProps {
  children: React.ReactNode;
  /** Maps to Atempo's meaning system. Default `neutral`. */
  tone?: 'tempo' | 'steady' | 'cold' | 'attention' | 'neutral';
  /** Prepend a status dot. */
  dot?: boolean;
  /** Make the dot breathe (the single `pulseDot` animation) — for live state. */
  live?: boolean;
  style?: React.CSSProperties;
}

/**
 * A small mono status capsule — `in tempo · 2h`, `live :4100`, `2 uncommitted`.
 * Never decorative; the tone always carries meaning.
 */
export function Badge(props: BadgeProps): React.ReactElement;
