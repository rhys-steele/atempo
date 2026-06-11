import * as React from 'react';

export interface KeyCapProps {
  children: React.ReactNode;
  style?: React.CSSProperties;
}

/** A keyboard-shortcut cap shown in a faint hairline (⌘K). For command hints. */
export function KeyCap(props: KeyCapProps): React.ReactElement;
