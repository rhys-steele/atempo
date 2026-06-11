import { invoke } from '@tauri-apps/api/core'

export interface Port {
  service: string
  external: number
  internal: number
  protocol: string
}

export interface Service {
  name: string
  status: string
  url?: string
}

export interface Project {
  name: string
  path: string
  framework: string
  version: string
  status: string
  created_at: string
  last_accessed: string
  ports: Port[]
  urls: string[]
  services: Service[]
  git_branch?: string
  git_status?: string
}

export interface DoctorCheck {
  name: string
  ok: boolean
  hint: string
}

export interface StreamEvent {
  type: 'log' | 'success' | 'error' | 'done'
  data?: string
}

export const listProjects = (): Promise<Project[]> =>
  invoke('list_projects')

export const describeProject = (name: string): Promise<Project> =>
  invoke('describe_project', { name })

export const startProject = (name: string): Promise<void> =>
  invoke('start_project', { name })

export const stopProject = (name: string): Promise<void> =>
  invoke('stop_project', { name })

export const removeProject = (name: string): Promise<void> =>
  invoke('remove_project', { name })

export const getStatus = (name: string): Promise<Service[]> =>
  invoke('get_status', { name })

export const runDoctor = (): Promise<DoctorCheck[]> =>
  invoke('run_doctor')

export const reconfigure = (name: string): Promise<void> =>
  invoke('reconfigure_project', { name })
