<script setup lang="ts">
import { computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useCatalogStore } from '@/stores/catalog'
import SearchField from '@/components/SearchField.vue'
import FilterChips from '@/components/FilterChips.vue'
import TemplateCard from '@/components/TemplateCard.vue'
import CardSkeleton from '@/components/CardSkeleton.vue'
import EmptyState from '@/components/EmptyState.vue'
import AppIcon from '@/components/AppIcon.vue'

const store = useCatalogStore()
const route = useRoute()
const router = useRouter()

const query = computed({
  get: () => store.q,
  set: (v: string) => store.setQuery(v),
})

// The URL is the shareable source of truth: every filter lands in the query
// string (replace, so typing does not fill the back stack), and the store
// hydrates from it on load and on back/forward.
function syncURL() {
  const q: Record<string, string> = {}
  if (store.q) q.q = store.q
  if (store.source) q.source = store.source
  if (store.category) q.category = store.category
  if (store.tag) q.tag = store.tag
  if (store.page > 1) q.page = String(store.page)
  void router.replace({ name: 'home', query: q })
}

watch(() => [store.q, store.source, store.category, store.tag, store.page], syncURL)

watch(
  () => route.query,
  (q) => {
    if (route.name !== 'home') return
    store.hydrate({
      q: (q.q as string) || '',
      source: (q.source as string) || '',
      category: (q.category as string) || '',
      tag: (q.tag as string) || '',
      page: Number(q.page) || 1,
    })
  },
)

onMounted(() => {
  document.title = 'Miabi Marketplace — official & community app templates'
  store.hydrate({
    q: (route.query.q as string) || '',
    source: (route.query.source as string) || '',
    category: (route.query.category as string) || '',
    tag: (route.query.tag as string) || '',
    page: Number(route.query.page) || 1,
  })
  void store.loadCategories()
})

const countLabel = computed(() => {
  if (store.loading) return 'Searching…'
  const n = store.total
  const noun = n === 1 ? 'template' : 'templates'
  return store.hasFilters ? `${n} ${noun} match` : `${n} ${noun}`
})
</script>

<template>
  <section class="hero">
    <div class="container">
      <h1>Deploy anything on Miabi</h1>
      <p class="sub">
        Official and community templates — apps, databases and full stacks, installed in one click
        with managed dependencies wired for you.
      </p>

      <div class="search-wrap">
        <SearchField
          v-model="query"
          :busy="store.refreshing"
        />
      </div>

      <FilterChips
        :categories="store.categories"
        :category="store.category"
        :source="store.source"
        @category="store.setCategory($event)"
        @source="store.setSource($event)"
      />
    </div>
  </section>

  <section class="container results">
    <div class="bar">
      <p
        class="count"
        aria-live="polite"
        aria-atomic="true"
      >
        {{ countLabel }}
      </p>
      <button
        v-if="store.hasFilters"
        type="button"
        class="btn btn-secondary btn-sm"
        @click="store.clearFilters()"
      >
        <AppIcon name="filter-remove" />
        Clear filters
      </button>
    </div>

    <p
      v-if="store.error"
      class="error"
      role="alert"
    >
      <AppIcon name="alert" />
      {{ store.error }}
    </p>

    <div
      v-if="store.loading"
      class="grid"
    >
      <CardSkeleton
        v-for="n in 8"
        :key="n"
      />
    </div>

    <div
      v-else-if="store.items.length"
      class="grid"
      :class="{ dim: store.refreshing }"
    >
      <TemplateCard
        v-for="t in store.items"
        :key="t.name"
        :item="t"
      />
    </div>

    <EmptyState
      v-else-if="store.isEmpty"
      icon="no-results"
      title="No templates match"
      hint="Try a shorter term, or clear the filters to see the whole catalog."
    >
      <button
        v-if="store.hasFilters"
        type="button"
        class="btn btn-primary btn-sm"
        @click="store.clearFilters()"
      >
        Clear filters
      </button>
    </EmptyState>

    <nav
      v-if="store.totalPages > 1"
      class="pager"
      aria-label="Pagination"
    >
      <button
        class="btn btn-secondary btn-sm"
        type="button"
        :disabled="store.page <= 1"
        @click="store.setPage(store.page - 1)"
      >
        <AppIcon name="chevron-left" /> Previous
      </button>
      <span class="pos">Page {{ store.page }} of {{ store.totalPages }}</span>
      <button
        class="btn btn-secondary btn-sm"
        type="button"
        :disabled="store.page >= store.totalPages"
        @click="store.setPage(store.page + 1)"
      >
        Next <AppIcon name="chevron-right" />
      </button>
    </nav>
  </section>
</template>

<style scoped>
.hero {
  padding: 56px 0 28px;
  background:
    radial-gradient(60% 120% at 50% 0%, color-mix(in srgb, var(--primary-500) 12%, transparent), transparent 70%),
    var(--bg-primary);
  border-bottom: 1px solid var(--border-primary);
}
.hero h1 {
  font-size: clamp(28px, 4vw, 40px);
  font-weight: 700;
  letter-spacing: -0.02em;
  line-height: 1.15;
}
.sub {
  margin-top: 10px;
  max-width: 62ch;
  color: var(--text-secondary);
  font-size: 16px;
}
.search-wrap {
  margin: 26px 0 18px;
  max-width: 620px;
}

.results {
  padding: 26px 0 8px;
}
.bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 16px;
}
.count {
  font-size: 14px;
  color: var(--text-tertiary);
  font-variant-numeric: tabular-nums;
}

.error {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 14px;
  margin-bottom: 16px;
  border: 1px solid var(--border-primary);
  border-radius: var(--radius);
  background: var(--bg-primary);
  color: var(--text-secondary);
  font-size: 14px;
}

.grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(258px, 1fr));
  gap: 16px;
  transition: opacity var(--transition);
}
/* A re-query dims the current results instead of tearing them down, so the page
   never jumps while you type. */
.grid.dim {
  opacity: 0.55;
}

.pager {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 16px;
  padding: 32px 0 8px;
}
.pos {
  font-size: 13px;
  color: var(--text-tertiary);
  font-variant-numeric: tabular-nums;
}
</style>
