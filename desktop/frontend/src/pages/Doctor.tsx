import { useQuery } from '@tanstack/react-query'
import { CheckCircle, XCircle, RefreshCw } from 'lucide-react'
import { runDoctor } from '../lib/tauri'

export default function Doctor() {
  const { data, isLoading, error, refetch, isFetching } = useQuery({
    queryKey: ['doctor'],
    queryFn: runDoctor,
  })

  return (
    <div className="p-6 max-w-lg">
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-lg font-semibold text-white">System Doctor</h1>
        <button
          onClick={() => refetch()}
          disabled={isFetching}
          className="flex items-center gap-1.5 text-xs text-white/40 hover:text-white/70 transition-colors disabled:opacity-40"
        >
          <RefreshCw size={13} className={isFetching ? 'animate-spin' : ''} />
          Re-check
        </button>
      </div>

      {isLoading && <div className="text-white/30 text-sm">Checking dependencies…</div>}
      {error && <div className="text-red-400 text-sm">Failed to run checks.</div>}

      {data && (
        <div className="space-y-2">
          {data.map(check => (
            <div
              key={check.name}
              className={`flex items-start gap-3 px-4 py-3 rounded-xl border ${
                check.ok
                  ? 'border-emerald-500/20 bg-emerald-500/5'
                  : 'border-red-500/20 bg-red-500/5'
              }`}
            >
              {check.ok
                ? <CheckCircle size={16} className="text-emerald-400 mt-0.5 shrink-0" />
                : <XCircle size={16} className="text-red-400 mt-0.5 shrink-0" />
              }
              <div>
                <div className="text-sm text-white">{check.name}</div>
                {!check.ok && check.hint && (
                  <div className="text-xs text-red-400/80 mt-0.5">{check.hint}</div>
                )}
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
