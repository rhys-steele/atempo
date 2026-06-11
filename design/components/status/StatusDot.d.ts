import * as React from 'react';

export interface StatusDotProps {
  /** teal `live`, amber `attention`, cream `steady`/`cold`, muted `neutral`. */
  tone?: 'live' | 'attention' | 'steady' | 'cold' | 'neutral';
  /** Turn on the slow pulse breath (for live runtime). */
  live?: boolean;
  size?: number;
  style?: React.CSSProperties;
}

/** The atomic status mark — a small round dot that optionally breathes when live. */
export function StatusDot(props: StatusDotProps): React.ReactElement;
