import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { api } from '@/api/client'
import type { CategoryFacet, Listing, Page } from '@/api/types'

const PER_PAGE = 24
const DEBOUNCE_MS = 180

// The catalog store owns search as a single piece of state: every filter change
// funnels through run(), which debounces, cancels the in-flight request, and
// drops responses that arrive out of order. Views never call the API directly.
export const useCatalogStore = defineStore('catalog', () => {
  const q = ref('')
  const source = ref('')
  const category = ref('')
  const tag = ref('')
  const page = ref(1)

  const items = ref<Listing[]>([])
  const total = ref(0)
  const totalPages = ref(0)
  const categories = ref<CategoryFacet[]>([])

  // `loading` is the first load (shows skeletons); `refreshing` is a keystroke
  // re-query, which keeps the current results on screen and only dims them —
  // swapping to skeletons on every character makes the grid strobe.
  const loading = ref(false)
  const refreshing = ref(false)
  const error = ref('')

  let timer: ReturnType<typeof setTimeout> | undefined
  let controller: AbortController | undefined
  let seq = 0

  const hasFilters = computed(() => !!(q.value || source.value || category.value || tag.value))
  const isEmpty = computed(() => !loading.value && items.value.length === 0)

  async function fetchPage() {
    controller?.abort()
    controller = new AbortController()
    const mine = ++seq

    if (items.value.length === 0) loading.value = true
    else refreshing.value = true
    error.value = ''

    try {
      const res: Page = await api.search(
        {
          q: q.value || undefined,
          source: source.value || undefined,
          category: category.value || undefined,
          tag: tag.value || undefined,
          page: page.value,
          per_page: PER_PAGE,
        },
        controller.signal,
      )
      if (mine !== seq) return // a newer query already won
      items.value = res.items ?? []
      total.value = res.total
      totalPages.value = res.total_pages
    } catch (e: unknown) {
      if (e instanceof DOMException && e.name === 'AbortError') return
      if (mine !== seq) return
      error.value = e instanceof Error ? e.message : 'Search failed'
      items.value = []
      total.value = 0
      totalPages.value = 0
    } finally {
      if (mine === seq) {
        loading.value = false
        refreshing.value = false
      }
    }
  }

  // run debounces typing but applies chip/facet clicks immediately — a click is
  // a deliberate act and waiting on it feels broken.
  function run(immediate = false) {
    if (timer) clearTimeout(timer)
    if (immediate) {
      void fetchPage()
      return
    }
    timer = setTimeout(() => void fetchPage(), DEBOUNCE_MS)
  }

  function setQuery(value: string) {
    q.value = value
    page.value = 1
    run()
  }

  function setSource(value: string) {
    source.value = source.value === value ? '' : value
    page.value = 1
    run(true)
  }

  function setCategory(value: string) {
    category.value = category.value === value ? '' : value
    page.value = 1
    run(true)
  }

  function setTag(value: string) {
    tag.value = tag.value === value ? '' : value
    page.value = 1
    run(true)
  }

  function setPage(value: number) {
    page.value = Math.max(1, value)
    run(true)
  }

  function clearFilters() {
    q.value = ''
    source.value = ''
    category.value = ''
    tag.value = ''
    page.value = 1
    run(true)
  }

  // hydrate applies URL state on first paint (deep links and back/forward) and
  // only re-queries when something actually changed.
  function hydrate(state: { q?: string; source?: string; category?: string; tag?: string; page?: number }) {
    const next = {
      q: state.q ?? '',
      source: state.source ?? '',
      category: state.category ?? '',
      tag: state.tag ?? '',
      page: state.page && state.page > 0 ? state.page : 1,
    }
    const changed =
      next.q !== q.value ||
      next.source !== source.value ||
      next.category !== category.value ||
      next.tag !== tag.value ||
      next.page !== page.value
    q.value = next.q
    source.value = next.source
    category.value = next.category
    tag.value = next.tag
    page.value = next.page
    if (changed || items.value.length === 0) run(true)
  }

  async function loadCategories() {
    try {
      categories.value = await api.categories()
    } catch {
      categories.value = []
    }
  }

  return {
    q,
    source,
    category,
    tag,
    page,
    items,
    total,
    totalPages,
    categories,
    loading,
    refreshing,
    error,
    hasFilters,
    isEmpty,
    setQuery,
    setSource,
    setCategory,
    setTag,
    setPage,
    clearFilters,
    hydrate,
    loadCategories,
  }
})
