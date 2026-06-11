import * as React from 'react';

export interface ButtonProps {
  children: React.ReactNode;
  /** Visual weight. `primary` is the single strong teal action; default `secondary`. */
  variant?: 'primary' | 'secondary' | 'ghost' | 'accentOutline';
  size?: 'sm' | 'md';
  /** Leading functional glyph — a Unicode mark like ▶, →, +. No icon fonts. */
  glyph?: string;
  disabled?: boolean;
  onClick?: () => void;
  style?: React.CSSProperties;
}

/**
 * The Atempo action button — flat, warm, restrained. Reserve `primary` (teal fill)
 * for the one orchestration/re-entry action on a surface (Run, Resume with Claude).
 *
 * @startingPoint section="Core" subtitle="Flat warm button — primary / secondary / ghost" viewport="700x140"
 */
export function Button(props: ButtonProps): React.ReactElement;
