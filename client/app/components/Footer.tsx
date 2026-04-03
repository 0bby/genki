import { ExternalLinkIcon } from 'lucide-react'
import { useQuery } from '@tanstack/react-query'
import { getSyncStatus } from 'api/api'
import pkg from '../../package.json'

function timeAgo(dateStr: string): string {
    const seconds = Math.floor((Date.now() - new Date(dateStr).getTime()) / 1000)
    if (seconds < 60) return 'just now'
    const minutes = Math.floor(seconds / 60)
    if (minutes < 60) return `${minutes}m ago`
    const hours = Math.floor(minutes / 60)
    if (hours < 24) return `${hours}h ago`
    return `${Math.floor(hours / 24)}d ago`
}

export default function Footer() {
    const { data: syncStatus } = useQuery({
        queryKey: ['sync-status-footer'],
        queryFn: getSyncStatus,
        refetchInterval: 60000,
        staleTime: 30000,
    })

    const wgerCursor = syncStatus?.sources?.find(s => s.source === 'wger')
    const fitbitCursor = syncStatus?.sources?.find(s => s.source === 'fitbit')

    return (
        <div className="mx-auto py-10 pt-20 color-fg-tertiary text-sm">
            <ul className="flex flex-col items-center w-sm justify-around">
                {(wgerCursor || fitbitCursor) && (
                    <li className="text-xs mb-2">
                        {wgerCursor && <span>wger: {timeAgo(wgerCursor.last_synced_at)}</span>}
                        {wgerCursor && fitbitCursor && <span> · </span>}
                        {fitbitCursor && <span>fitbit: {timeAgo(fitbitCursor.last_synced_at)}</span>}
                    </li>
                )}
                <li>Genki {import.meta.env.VITE_GENKI_VERSION || pkg.version}</li>
                <li><a href="https://github.com/0bby/genki" target="_blank" className="link-underline">View the source on GitHub <ExternalLinkIcon className='inline mb-1' size={14}/></a></li>
            </ul>
        </div>
    )
}
