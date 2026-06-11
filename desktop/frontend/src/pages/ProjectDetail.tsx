import { useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { listen } from '@tauri-apps/api/event'
import { invoke } from '@tauri-apps/api/core'
import { openUrl } from '@tauri-apps/plugin-opener'
import { Play, Square, RotateCcw, Trash2, ExternalLink, ArrowLeft } from 'lucide-react'
import { useProject, useStartProject, useStopProject, useRemoveProject } from '../hooks/useProjects'
import ServiceBadge from '../components/ServiceBadge'
import LogStream from '../components/LogStream'

interface LogLine { type: 'log' | 'success' | 'error' | 'done'; data?: string }

export default function ProjectDetail() {
  const { name } = useParams<{ name: string }>()
  const navigate = useNavigate()
  const { data: project, isLoading } = useProject(name!)
  const start = useStartProject()
  const stop = useStopProject()
  const remove = useRemoveProject()

  const [logs, setLogs] = useState<LogLine[]>([])
  const [streaming, setStreaming] = useState(false)

  const streamLogs = async () => {
    setLogs([])
    setStreaming(true)
    const unlisten = await listen<LogLine>('logs:event', event => {
      const line = event.payload
      setLogs(prev => [...prev, line])
      if (line.type === 'done' || line.type === 'error') { setStreaming(false); unlisten() }
    })
    try {
      await invoke('stream_logs', { name })
    } catch (e) {
      setLogs(prev => [...prev, { type: 'error', data: String(e) }])
      setStreaming(false)
      unlisten()
    }
  }

  const handleRemove = async () => {
    if (!confirm(`Remove ${name}? This deregisters the project but does not delete files.`)) return
    await remove.mutateAsync(name!)
    navigate('/')
  }

  if (isLoading) return <div className="p-6 text-white/30 text-sm">Loading…</div>
  if (!project) return <div className="p-6 text-white/30 text-sm">Project not found.</div>

  const isRunning = project.status === 'running' || project.status === 'healthy'

  return (
    <div className="p-6 space-y-6 max-w-2xl">
      <div className="flex items-center gap-3">
        <button onClick={() => navigate('/')} className="text-white/30 hover:text-white/70 transition-colors">
          <ArrowLeft size={16} />
        </button>
        <div>
          <h1 className="text-lg font-semibold text-white">{project.name}</h1>
          <span className="text-xs text-white/40">{project.framework}{project.version ? ` ${project.version}` : ''} · {project.path}</span>
        </div>
        <ServiceBadge status={project.status || 'stopped'} />
      </div>

      {/* Actions */}
      <div className="flex flex-wrap gap-2">
        <button
          onClick={() => isRunning ? stop.mutate(name!) : start.mutate(name!)}
          disabled={start.isPending || stop.isPending}
          className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-sm bg-white/6 hover:bg-white/10 disabled:opacity-40 transition-colors"
        >
          {isRunning ? <Square size={13} /> : <Play size={13} />}
          {isRunning ? 'Stop' : 'Start'}
        </button>
        <button
          onClick={streamLogs}
          disabled={streaming}
          className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-sm bg-white/6 hover:bg-white/10 disabled:opacity-40 transition-colors"
        >
          Logs
        </button>
        <button
          onClick={() => invoke('reconfigure_project', { name })}
          className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-sm bg-white/6 hover:bg-white/10 transition-colors"
        >
          <RotateCcw size={13} /> Reconfigure
        </button>
        <button
          onClick={handleRemove}
          className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-sm text-red-400/70 hover:text-red-400 hover:bg-red-500/10 transition-colors ml-auto"
        >
          <Trash2 size={13} /> Remove
        </button>
      </div>

      {/* Services */}
      <section>
        <h2 className="text-xs font-medium text-white/40 uppercase tracking-wider mb-2">Services</h2>
        <div className="border border-white/8 rounded-xl overflow-hidden">
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-white/6 text-white/30 text-xs">
                <th className="text-left px-4 py-2 font-normal">Service</th>
                <th className="text-left px-4 py-2 font-normal">Status</th>
                <th className="text-left px-4 py-2 font-normal">URL</th>
              </tr>
            </thead>
            <tbody>
              {(project.services ?? []).map(svc => (
                <tr key={svc.name} className="border-b border-white/4 last:border-0">
                  <td className="px-4 py-2.5 text-white/80">{svc.name}</td>
                  <td className="px-4 py-2.5"><ServiceBadge status={svc.status} /></td>
                  <td className="px-4 py-2.5">
                    {svc.url ? (
                      <button
                        onClick={() => openUrl(svc.url!)}
                        className="flex items-center gap-1 text-white/40 hover:text-white/70 text-xs transition-colors"
                      >
                        <ExternalLink size={11} /> {svc.url}
                      </button>
                    ) : <span className="text-white/20">—</span>}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </section>

      {/* Ports */}
      {(project.ports ?? []).length > 0 && (
        <section>
          <h2 className="text-xs font-medium text-white/40 uppercase tracking-wider mb-2">Ports</h2>
          <div className="space-y-1">
            {project.ports.map((p, i) => (
              <div key={i} className="flex items-center gap-2 text-sm text-white/60">
                <span className="text-white">{p.external}</span>
                <span className="text-white/20">→</span>
                <span>{p.internal}</span>
                <span className="text-white/30 text-xs">{p.service}</span>
              </div>
            ))}
          </div>
        </section>
      )}

      {/* Log stream */}
      {logs.length > 0 && (
        <section>
          <h2 className="text-xs font-medium text-white/40 uppercase tracking-wider mb-2">Logs</h2>
          <LogStream lines={logs} className="h-64" />
        </section>
      )}
    </div>
  )
}
