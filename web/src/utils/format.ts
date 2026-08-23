import type { Listing } from '@/api/types'

function plural(n: number, word: string) {
  return `${n} ${word}${n === 1 ? '' : 's'}`
}

// provision summarizes what an install creates — the single most useful thing
// on a card: "1 app · 2 databases · 4 volumes · 1 config".
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

// monogram is the icon fallback: templates without an icon URL still get a
// stable, readable tile instead of a hole in the grid.
export function monogram(name: string): string {
  return (name.trim()[0] ?? '?').toUpperCase()
}

// hueOf derives a stable accent per template so monogram tiles are
// distinguishable at a glance without hand-assigning colors.
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
