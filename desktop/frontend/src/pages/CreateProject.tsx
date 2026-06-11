import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { open } from '@tauri-apps/plugin-dialog'
import { invoke } from '@tauri-apps/api/core'
import { listen } from '@tauri-apps/api/event'
import { FolderOpen } from 'lucide-react'
import LogStream from '../components/LogStream'

const FRAMEWORKS = [
  { id: 'laravel', label: 'Laravel', description: 'PHP · MySQL · Redis · Nginx' },
  { id: 'fastapi', label: 'FastAPI', description: 'Python 3.12 · PostgreSQL · Redis' },
  { id: 'bun',     label: 'Bun',     description: 'TypeScript · Bun runtime' },
  { id: 'python',  label: 'Python',  description: 'Python 3.12 · no ports' },
]

type Step = 'framework' | 'details' | 'progress'

interface LogLine {
  type: 'log' | 'success' | 'error' | 'done'
  data?: string
}

export default function CreateProject() {
  const navigate = useNavigate()
  const [step, setStep] = useState<Step>('framework')
  const [framework, setFramework] = useState('')
  const [name, setName] = useState('')
  const [version, setVersion] = useState('')
  const [directory, setDirectory] = useState('')
  const [logs, setLogs] = useState<LogLine[]>([])
  const [done, setDone] = useState(false)
  const [errored, setErrored] = useState(false)

  const pickDirectory = async () => {
    const selected = await open({ directory: true, multiple: false })
    if (typeof selected === 'string') setDirectory(selected)
  }

  const handleCreate = async () => {
    setStep('progress')
    setLogs([])
    setDone(false)
    setErrored(false)

    const unlisten = await listen<LogLine>('create:event', event => {
      const line = event.payload
      setLogs(prev => [...prev, line])
      if (line.type === 'done') { setDone(true); unlisten() }
      if (line.type === 'error') { setErrored(true); unlisten() }
    })

    try {
      await invoke('create_project', { framework, name, version: version || undefined, directory })
    } catch (e) {
      setLogs(prev => [...prev, { type: 'error', data: String(e) }])
      setErrored(true)
      unlisten()
    }
  }

  return (
    <div className="p-6 max-w-xl">
      <h1 className="text-lg font-semibold text-white mb-6">New Project</h1>

      {step === 'framework' && (
        <div className="space-y-3">
          <p className="text-sm text-white/50 mb-4">Choose a framework</p>
          {FRAMEWORKS.map(fw => (
            <button
              key={fw.id}
              onClick={() => { setFramework(fw.id); setStep('details') }}
              className="w-full flex items-center justify-between px-4 py-3 rounded-xl border border-white/8 bg-white/4 hover:border-white/20 hover:bg-white/6 transition-colors text-left"
            >
              <div>
                <div className="text-sm font-medium text-white">{fw.label}</div>
                <div className="text-xs text-white/40">{fw.description}</div>
              </div>
              <span className="text-white/20">→</span>
            </button>
          ))}
        </div>
      )}

      {step === 'details' && (
        <div className="space-y-4">
          <p className="text-sm text-white/50">
            Configure your <span className="text-white">{framework}</span> project
          </p>

          <div>
            <label className="block text-xs text-white/50 mb-1">Project name</label>
            <input
              type="text"
              value={name}
              onChange={e => setName(e.target.value)}
              placeholder="my-app"
              className="w-full bg-white/6 border border-white/10 rounded-lg px-3 py-2 text-sm text-white placeholder-white/20 focus:outline-none focus:border-violet-500"
            />
          </div>

          {framework === 'laravel' && (
            <div>
              <label className="block text-xs text-white/50 mb-1">Version (optional)</label>
              <input
                type="text"
                value={version}
                onChange={e => setVersion(e.target.value)}
                placeholder="e.g. 12"
                className="w-full bg-white/6 border border-white/10 rounded-lg px-3 py-2 text-sm text-white placeholder-white/20 focus:outline-none focus:border-violet-500"
              />
            </div>
          )}

          <div>
            <label className="block text-xs text-white/50 mb-1">Parent directory</label>
            <button
              onClick={pickDirectory}
              className="w-full flex items-center gap-2 bg-white/6 border border-white/10 rounded-lg px-3 py-2 text-sm text-left hover:border-white/20 transition-colors"
            >
              <FolderOpen size={14} className="text-white/40 shrink-0" />
              <span className={directory ? 'text-white' : 'text-white/25'}>{directory || 'Choose folder…'}</span>
            </button>
          </div>

          <div className="flex gap-2 pt-2">
            <button
              onClick={() => setStep('framework')}
              className="px-3 py-2 rounded-lg text-sm text-white/40 hover:text-white/70 transition-colors"
            >
              Back
            </button>
            <button
              onClick={handleCreate}
              disabled={!name || !directory}
              className="flex-1 px-3 py-2 rounded-lg text-sm bg-violet-600 hover:bg-violet-500 disabled:opacity-40 text-white transition-colors"
            >
              Create project
            </button>
          </div>
        </div>
      )}

      {step === 'progress' && (
        <div className="space-y-4">
          <p className="text-sm text-white/50">Creating <span className="text-white">{name}</span>…</p>
          <LogStream lines={logs} className="h-72" />
          {(done || errored) && (
            <div className="flex gap-2">
              {done && (
                <button
                  onClick={() => navigate(`/project/${name}`)}
                  className="flex-1 px-3 py-2 rounded-lg text-sm bg-violet-600 hover:bg-violet-500 text-white transition-colors"
                >
                  Open project
                </button>
              )}
              <button
                onClick={() => navigate('/')}
                className="px-3 py-2 rounded-lg text-sm text-white/40 hover:text-white/70 transition-colors"
              >
                {done ? 'Back to projects' : 'Cancel'}
              </button>
            </div>
          )}
        </div>
      )}
    </div>
  )
}
