import type { Listing } from '@/api/types'

function plural(n: number, word: string) {
  return `${n} ${word}${n === 1 ? '' : 's'}`
}

export function provision(l: Pick<Listing, 'applications' | 'databases' | 'volumes' | 'configs'>): string {
  const parts: string[] = []
  if (l.applications > 0) parts.push(plural(l.applications, 'app'))
  if (l.databases > 0) parts.push(plural(l.databases, 'database'))
  if (l.volumes > 0) parts.push(plural(l.volumes, 'volume'))
  // ?? 0: a listing from an older marketplace build carries no config count.
  if ((l.configs ?? 0) > 0) parts.push(plural(l.configs, 'config'))
  return parts.length ? parts.join(' · ') : 'no dependencies'
}

export function isURL(s?: string): boolean {
  return !!s && (s.startsWith('http://') || s.startsWith('https://'))
}

export function monogram(name: string): string {
  return (name.trim()[0] ?? '?').toUpperCase()
}

export function hueOf(seed: string): number {
  let h = 0
  for (let i = 0; i < seed.length; i++) h = (h * 31 + seed.charCodeAt(i)) % 360
  return h
}

export function engineLabel(engine: string): string {
  const map: Record<string, string> = {
    postgres: 'PostgreSQL',
    mysql: 'MySQL',
    mariadb: 'MariaDB',
    mongodb: 'MongoDB',
    redis: 'Redis',
    libsql: 'libSQL',
  }
  return map[engine] ?? engine
}

// Dates are calendar days (YYYY-MM-DD); parsing them as UTC avoids an off-by-one in negative offsets.
function parseDay(day: string): Date | null {
  const d = new Date(`${day}T00:00:00Z`)
  return Number.isNaN(d.getTime()) ? null : d
}

export function shortDate(day?: string): string {
  const d = day ? parseDay(day) : null
  if (!d) return ''
  return d.toLocaleDateString(undefined, { year: 'numeric', month: 'short', day: 'numeric', timeZone: 'UTC' })
}

export function relativeDay(day?: string, now = new Date()): string {
  const d = day ? parseDay(day) : null
  if (!d) return ''
  const today = Date.UTC(now.getFullYear(), now.getMonth(), now.getDate())
  const days = Math.round((today - d.getTime()) / 86_400_000)
  if (days <= 0) return 'today'
  if (days === 1) return 'yesterday'
  if (days < 30) return `${days} days ago`
  return shortDate(day)
}
