import * as React from 'react';

export interface RuntimePillProps {
  /** `live` (teal + ports), `ready`, `bypassed` (dashed inert), `none`. */
  state?: 'live' | 'ready' | 'bypassed' | 'none';
  /** Override text, e.g. "live :4100 :1420" or "ready — 1 service". */
  label?: string;
  style?: React.CSSProperties;
}

/**
 * Read-only runtime state shown as a *result*, never a control. Atempo brings things
 * up for you and confirms they're live — there are no start/stop affordances here.
 */
export function RuntimePill(props: RuntimePillProps): React.ReactElement;
