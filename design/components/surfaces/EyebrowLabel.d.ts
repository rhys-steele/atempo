import * as React from 'react';

export interface EyebrowLabelProps {
  children: React.ReactNode;
  /** `muted` (default), `accent` (teal — for the recap label), `faint`. */
  tone?: 'muted' | 'accent' | 'faint';
  style?: React.CSSProperties;
}

/** The tracked-out uppercase mono label that titles every panel. */
export function EyebrowLabel(props: EyebrowLabelProps): React.ReactElement;
