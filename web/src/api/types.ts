// Keep field names in sync with the Go JSON tags in internal/api and internal/catalog.

export type Source = 'official' | 'community'

export interface Author {
  name: string
  email?: string
  website?: string
}

// Listing is the card view of a template (latest version).
export interface Listing {
  name: string
  display_name: string
  description: string
  category: string
  icon?: string
  tags?: string[]
  homepage?: string
  author?: Author
  source: Source
  featured?: boolean
  version: string
  versions: string[]
  applications: number
  databases: number
  volumes: number
  configs: number
  db_only: boolean
  updated?: string
  deprecated?: boolean
  deprecation_note?: string
  replaced_by?: string
  freshness?: Freshness
}

export type FreshnessStatus = 'current' | 'update' | 'major' | 'held' | 'unknown' | 'error'

export interface AppFreshness {
  app: string
  image: string
  tag: string
  latest_in_major?: string
  latest?: string
  status: FreshnessStatus
  error?: string
  hold?: string
}

export interface Freshness {
  version: string
  status: FreshnessStatus
  checked_at?: string
  apps: AppFreshness[]
}

export interface Page {
  items: Listing[]
  page: number
  per_page: number
  total: number
  total_pages: number
}

export interface CategoryFacet {
  category: string
  count: number
}

export interface TemplateInput {
  key: string
  label?: string
  help?: string
  type?: 'string' | 'password' | 'bool' | 'select' | 'number'
  default?: string
  placeholder?: string
  pattern?: string
  options?: string[]
  required?: boolean
  generate?: boolean
  length?: number
}

export interface ManifestDatabase {
  name: string
  engine: string
  version?: string
  placement?: 'auto' | 'dedicated' | 'shared'
}

export interface ManifestPort {
  container: number
  scheme?: string
}

export interface ManifestMount {
  volume?: string
  config?: string
  key?: string
  mode?: string
  path: string
  readOnly?: boolean
}

export interface ManifestConfig {
  name: string
  files: Record<string, string>
  mode?: string
  sensitive?: boolean
  delimiters?: string[]
}

export interface ManifestApp {
  name: string
  primary?: boolean
  image: string
  tag?: string
  command?: string[]
  ports?: ManifestPort[]
  env?: Record<string, string>
  secretEnv?: string[]
  mounts?: ManifestMount[]
  healthcheck?: { type?: string; path?: string; command?: string; port?: number }
}

export interface TemplateManifest {
  apiVersion: string
  kind: string
  metadata: {
    name: string
    displayName: string
    version: string
    description?: string
    category?: string
    icon?: string
    homepage?: string
    author?: Author
    tags?: string[]
    minMiabi?: string
  }
  inputs?: TemplateInput[]
  databases?: ManifestDatabase[]
  volumes?: { name: string }[]
  configs?: ManifestConfig[]
  stack?: { description?: string; env?: Record<string, string> }
  applications?: ManifestApp[]
}

export interface Release {
  released?: string
  notes?: string
  breaking?: boolean
  requires_from?: string
  deprecated?: boolean
  deprecation_note?: string
  images?: Record<string, string>
}

export interface VersionImage {
  app: string
  image: string
  tag: string
}

export interface VersionRef {
  version: string
  digest: string
  release?: Release
  images?: VersionImage[]
}

export interface TemplateDetail {
  entry: Listing
  source: Source
  featured?: boolean
  meta: {
    featured?: boolean
    screenshots?: string[]
    source_repo?: string
    maintainer?: string
    deprecated?: boolean
    deprecation_note?: string
    replaced_by?: string
  }
  readme?: string
  versions: VersionRef[]
  manifest: TemplateManifest
}

// BuildInfo is the running server build; version is "dev" outside a release.
export interface BuildInfo {
  version: string
  commit_id: string
  build_date: string
}

export interface Envelope<T> {
  success: boolean
  data: T
  error: { code: string; message: string; error?: string } | null
}
