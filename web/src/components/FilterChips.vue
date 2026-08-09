<script setup lang="ts">
import { computed, ref } from 'vue'
import type { CategoryFacet } from '@/api/types'

const props = defineProps<{
  categories: CategoryFacet[]
  category: string
  source: string
}>()
const emit = defineEmits<{ category: [string]; source: [string] }>()

// Categories are chips rather than a <select>: the whole taxonomy is visible and
// one click away. Beyond a dozen they'd wrap into a wall, so the tail collapses.
const VISIBLE = 8
const expanded = ref(false)
const shown = computed(() =>
  expanded.value || props.categories.length <= VISIBLE
    ? props.categories
    : props.categories.slice(0, VISIBLE),
)
const hiddenCount = computed(() => Math.max(0, props.categories.length - shown.value.length))
</script>

<template>
  <div class="filters">
    <div
      class="sources"
      role="group"
      aria-label="Filter by source"
    >
      <button
        v-for="opt in [
          { value: '', label: 'All' },
          { value: 'official', label: 'Official' },
          { value: 'community', label: 'Community' },
        ]"
        :key="opt.value"
        type="button"
        class="seg"
        :class="{ active: source === opt.value }"
        :aria-pressed="source === opt.value"
        @click="emit('source', opt.value)"
      >
        {{ opt.label }}
      </button>
    </div>

    <div
      class="cats"
      role="group"
      aria-label="Filter by category"
    >
      <button
        v-for="c in shown"
        :key="c.category"
        type="button"
        class="chip"
        :class="{ active: category === c.category }"
        :aria-pressed="category === c.category"
        @click="emit('category', c.category)"
      >
        {{ c.category }}
        <span class="count">{{ c.count }}</span>
      </button>
      <button
        v-if="hiddenCount > 0 || expanded"
        type="button"
        class="chip more"
        @click="expanded = !expanded"
      >
        {{ expanded ? 'Show less' : `+${hiddenCount} more` }}
      </button>
    </div>
  </div>
</template>

<style scoped>
.filters {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 12px 16px;
}

.sources {
  display: inline-flex;
  padding: 3px;
  background: var(--bg-tertiary);
  border: 1px solid var(--border-primary);
  border-radius: 999px;
}
.seg {
  padding: 6px 14px;
  border: 0;
  border-radius: 999px;
  background: transparent;
  color: var(--text-secondary);
  font-family: inherit;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all var(--transition);
}
.seg:hover {
  color: var(--text-primary);
}
.seg.active {
  background: var(--bg-primary);
  color: var(--primary-600);
  box-shadow: var(--shadow-sm);
}

.cats {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}
.chip {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 5px 12px;
  border: 1px solid var(--border-primary);
  border-radius: 999px;
  background: var(--bg-primary);
  color: var(--text-secondary);
  font-family: inherit;
  font-size: 13px;
  cursor: pointer;
  transition: all var(--transition);
}
.chip:hover {
  border-color: var(--primary-200);
  color: var(--text-primary);
}
.chip.active {
  background: var(--primary-50);
  border-color: var(--primary-200);
  color: var(--primary-600);
  font-weight: 500;
}
.count {
  font-size: 11px;
  color: var(--text-muted);
  font-variant-numeric: tabular-nums;
}
.chip.active .count {
  color: var(--primary-500);
}
.more {
  color: var(--text-tertiary);
}
</style>
