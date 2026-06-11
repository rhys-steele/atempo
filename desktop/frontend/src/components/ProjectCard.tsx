import { Play, Square, ExternalLink, ChevronRight } from 'lucide-react'
import { useNavigate } from 'react-router-dom'
import type { Project } from '../lib/tauri'
import ServiceBadge from './ServiceBadge'
import { useStartProject, useStopProject } from '../hooks/useProjects'

const frameworkColors: Record<string, string> = {
  laravel: 'text-red-400',
  fastapi: 'text-teal-400',
  bun:     'text-yellow-400',
  python:  'text-blue-400',
}

interface Props {
  project: Project
}

export default function ProjectCard({ project }: Props) {
  const navigate = useNavigate()
  const start = useStartProject()
  const stop = useStopProject()

  const isRunning = project.status === 'running' || project.status === 'healthy'
  const fwColor = frameworkColors[project.framework] ?? 'text-white/60'
  const primaryUrl = project.urls?.[0]

  return (
    <div
      className="bg-white/4 border border-white/8 rounded-xl p-4 flex flex-col gap-3 hover:border-white/15 transition-colors cursor-pointer"
      onClick={() => navigate(`/project/${project.name}`)}
    >
      <div className="flex items-start justify-between gap-2">
        <div>
          <h3 className="font-medium text-white">{project.name}</h3>
          <span className={`text-xs ${fwColor}`}>{project.framework}{project.version ? ` ${project.version}` : ''}</span>
        </div>
        <ServiceBadge status={project.status || 'stopped'} />
      </div>

      {primaryUrl && (
        <a
          href={primaryUrl}
          className="text-xs text-white/40 hover:text-white/70 flex items-center gap-1 w-fit"
          onClick={e => e.stopPropagation()}
        >
          <ExternalLink size={11} />
          {primaryUrl}
        </a>
      )}

      <div className="flex items-center gap-2 mt-auto pt-1 border-t border-white/6">
        <button
          onClick={e => { e.stopPropagation(); isRunning ? stop.mutate(project.name) : start.mutate(project.name) }}
          disabled={start.isPending || stop.isPending}
          className="flex items-center gap-1.5 px-2.5 py-1.5 rounded-lg text-xs bg-white/6 hover:bg-white/10 disabled:opacity-40 transition-colors"
        >
          {isRunning ? <Square size={12} /> : <Play size={12} />}
          {isRunning ? 'Stop' : 'Start'}
        </button>
        <button
          onClick={e => { e.stopPropagation(); navigate(`/project/${project.name}`) }}
          className="ml-auto flex items-center gap-1 text-xs text-white/30 hover:text-white/60 transition-colors"
        >
          Details <ChevronRight size={12} />
        </button>
      </div>
    </div>
  )
}
