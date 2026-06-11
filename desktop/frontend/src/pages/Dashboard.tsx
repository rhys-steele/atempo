import { Plus } from 'lucide-react'
import { useNavigate } from 'react-router-dom'
import { useProjects } from '../hooks/useProjects'
import ProjectCard from '../components/ProjectCard'

export default function Dashboard() {
  const navigate = useNavigate()
  const { data: projects, isLoading, error } = useProjects()

  return (
    <div className="p-6">
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-lg font-semibold text-white">Projects</h1>
        <button
          onClick={() => navigate('/create')}
          className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-sm bg-violet-600 hover:bg-violet-500 text-white transition-colors"
        >
          <Plus size={14} />
          New Project
        </button>
      </div>

      {isLoading && (
        <div className="text-white/30 text-sm">Loading projects…</div>
      )}

      {error && (
        <div className="text-red-400 text-sm">Failed to load projects.</div>
      )}

      {projects && projects.length === 0 && (
        <div className="text-center py-20 text-white/25">
          <p className="text-base mb-2">No projects yet</p>
          <p className="text-sm">Create one to get started.</p>
        </div>
      )}

      {projects && projects.length > 0 && (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
          {projects.map(p => (
            <ProjectCard key={p.name} project={p} />
          ))}
        </div>
      )}
    </div>
  )
}
