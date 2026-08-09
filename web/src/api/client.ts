import type { CategoryFacet, Envelope, Page, TemplateDetail } from './types'

const BASE = '/v1'

// ApiError carries the API's stable error code so callers can branch on it
// (a 404 renders the "not found" view rather than a failure banner).
export class ApiError extends Error {
  status: number
  code: string

  constructor(status: number, code: string, message: string) {
    super(message)
    this.name = 'ApiError'
    this.status = status
    this.code = code
  }
}

async function get<T>(path: string, params?: Record<string, string | number | undefined>, signal?: AbortSignal): Promise<T> {
  const url = new URL(BASE + path, window.location.origin)
  for (const [k, v] of Object.entries(params ?? {})) {
    if (v !== undefined && v !== '') url.searchParams.set(k, String(v))
  }
  const res = await fetch(url, { signal, headers: { Accept: 'application/json' } })
  const body = (await res.json().catch(() => null)) as Envelope<T> | null
  if (!res.ok || !body?.success) {
    throw new ApiError(res.status, body?.error?.code ?? 'REQUEST_FAILED', body?.error?.message ?? res.statusText)
  }
  return body.data
}

export interface SearchParams {
  q?: string
  source?: string
  category?: string
  tag?: string
  page?: number
  per_page?: number
}

export const api = {
  search(params: SearchParams, signal?: AbortSignal) {
    return get<Page>('/templates', { ...params }, signal)
  },
  categories(signal?: AbortSignal) {
    return get<CategoryFacet[]>('/categories', undefined, signal)
  },
  template(name: string, version?: string, signal?: AbortSignal) {
    return get<TemplateDetail>(`/templates/${encodeURIComponent(name)}`, { version }, signal)
  },
  // The raw manifest is served as YAML (not enveloped), with the digest as ETag.
  manifestURL(name: string, version: string) {
    return `${BASE}/templates/${encodeURIComponent(name)}/versions/${encodeURIComponent(version)}/manifest`
  },
  async manifestYAML(name: string, version: string, signal?: AbortSignal): Promise<string> {
    const res = await fetch(this.manifestURL(name, version), { signal })
    if (!res.ok) throw new ApiError(res.status, 'MANIFEST_FAILED', res.statusText)
    return res.text()
  },
}
