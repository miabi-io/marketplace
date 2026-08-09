// Shapes mirror the Go API (internal/api, internal/catalog). Keep field names in
// sync with the JSON tags there — the SPA is the only consumer that unwraps the
// Envelope, machine clients use /v1/export.

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
  db_only: boolean
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
  volume: string
  path: string
  readOnly?: boolean
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
  stack?: { description?: string; env?: Record<string, string> }
  applications?: ManifestApp[]
}

export interface VersionRef {
  version: string
  digest: string
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
  }
  readme?: string
  versions: VersionRef[]
  manifest: TemplateManifest
}

export interface Envelope<T> {
  success: boolean
  data: T
  error: { code: string; message: string; error?: string } | null
}
