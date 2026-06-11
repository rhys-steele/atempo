import { NavLink } from 'react-router-dom'
import { LayoutDashboard, Plus, Stethoscope } from 'lucide-react'

const links = [
  { to: '/', icon: LayoutDashboard, label: 'Projects' },
  { to: '/create', icon: Plus, label: 'New Project' },
  { to: '/doctor', icon: Stethoscope, label: 'Doctor' },
]

export default function Sidebar() {
  return (
    <aside className="w-56 shrink-0 border-r border-white/8 flex flex-col">
      <div className="px-5 py-5">
        <span className="text-sm font-semibold tracking-widest text-white/40 uppercase">Atempo</span>
      </div>
      <nav className="flex-1 px-2 space-y-0.5">
        {links.map(({ to, icon: Icon, label }) => (
          <NavLink
            key={to}
            to={to}
            end={to === '/'}
            className={({ isActive }) =>
              `flex items-center gap-2.5 px-3 py-2 rounded-lg text-sm transition-colors ${
                isActive
                  ? 'bg-white/10 text-white'
                  : 'text-white/50 hover:text-white/80 hover:bg-white/5'
              }`
            }
          >
            <Icon size={16} />
            {label}
          </NavLink>
        ))}
      </nav>
    </aside>
  )
}
