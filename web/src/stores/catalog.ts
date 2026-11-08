import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { api } from '@/api/client'
import type { CategoryFacet, Listing, Page } from '@/api/types'

const PER_PAGE = 24
const DEBOUNCE_MS = 180

// Every filter change funnels through run(), which debounces, cancels in-flight requests and drops stale responses.
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

  // Swapping to skeletons on every keystroke makes the grid strobe, so re-queries only dim the results.
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

  // Clicks apply immediately; only typing is debounced.
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
