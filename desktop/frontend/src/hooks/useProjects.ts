import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { listProjects, describeProject, getStatus, startProject, stopProject, removeProject } from '../lib/tauri'

export const useProjects = () =>
  useQuery({
    queryKey: ['projects'],
    queryFn: listProjects,
    refetchInterval: 10_000,
  })

export const useProject = (name: string) =>
  useQuery({
    queryKey: ['project', name],
    queryFn: () => describeProject(name),
    enabled: !!name,
  })

export const useProjectStatus = (name: string) =>
  useQuery({
    queryKey: ['status', name],
    queryFn: () => getStatus(name),
    enabled: !!name,
    refetchInterval: 10_000,
  })

export const useStartProject = () => {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: startProject,
    onSuccess: () => qc.invalidateQueries({ queryKey: ['projects'] }),
  })
}

export const useStopProject = () => {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: stopProject,
    onSuccess: () => qc.invalidateQueries({ queryKey: ['projects'] }),
  })
}

export const useRemoveProject = () => {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: removeProject,
    onSuccess: () => qc.invalidateQueries({ queryKey: ['projects'] }),
  })
}
