interface Props {
  status: string
}

const statusStyles: Record<string, string> = {
  running:   'bg-emerald-500/15 text-emerald-400 border-emerald-500/30',
  healthy:   'bg-emerald-500/15 text-emerald-400 border-emerald-500/30',
  stopped:   'bg-white/8 text-white/40 border-white/10',
  unhealthy: 'bg-red-500/15 text-red-400 border-red-500/30',
}

export default function ServiceBadge({ status }: Props) {
  const style = statusStyles[status] ?? statusStyles.stopped
  return (
    <span className={`inline-flex items-center gap-1.5 px-2 py-0.5 rounded-full text-xs border ${style}`}>
      <span className="size-1.5 rounded-full bg-current" />
      {status}
    </span>
  )
}
