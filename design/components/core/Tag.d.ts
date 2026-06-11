import * as React from 'react';

export interface TagProps {
  children: React.ReactNode;
  style?: React.CSSProperties;
}

/** A hairline mono stack/kind chip — `go`, `astro`, `ts`. Labels, never status. */
export function Tag(props: TagProps): React.ReactElement;
