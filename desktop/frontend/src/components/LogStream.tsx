import { useEffect, useRef } from 'react'

interface LogLine {
  type: 'log' | 'success' | 'error' | 'done'
  data?: string
}

interface Props {
  lines: LogLine[]
  className?: string
}

export default function LogStream({ lines, className = '' }: Props) {
  const bottomRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [lines])

  return (
    <div className={`font-mono text-xs bg-black/30 rounded-lg overflow-y-auto p-3 space-y-0.5 ${className}`}>
      {lines.map((line, i) => (
        <div
          key={i}
          className={
            line.type === 'error'   ? 'text-red-400' :
            line.type === 'success' ? 'text-emerald-400' :
            'text-white/60'
          }
        >
          {line.data ?? ''}
        </div>
      ))}
      <div ref={bottomRef} />
    </div>
  )
}
