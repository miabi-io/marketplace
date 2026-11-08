<script setup lang="ts">
import { computed, ref } from 'vue'
import type { Listing } from '@/api/types'
import { hueOf, isURL, monogram, provision, relativeDay, shortDate } from '@/utils/format'

const props = defineProps<{ item: Listing }>()

// Icons are third-party URLs; fall back to the monogram when one fails to load.
const iconFailed = ref(false)
const iconOk = computed(() => isURL(props.item.icon) && !iconFailed.value)
const hue = computed(() => hueOf(props.item.name))
const summary = computed(() => provision(props.item))
</script>

<template>
  <RouterLink
    class="card tpl"
    :to="{ name: 'template', params: { name: item.name } }"
  >
    <div class="head">
      <span
        class="tile"
        :style="{ '--tile-hue': hue }"
        aria-hidden="true"
      >
        <img
          v-if="iconOk"
          :src="item.icon"
          alt=""
          loading="lazy"
          decoding="async"
          @error="iconFailed = true"
        >
        <template v-else>{{ monogram(item.display_name) }}</template>
      </span>
      <div class="head-badges">
        <span
          v-if="item.deprecated"
          class="badge badge-warning"
        >deprecated</span>
        <span
          class="badge"
          :class="item.source === 'community' ? 'badge-community' : 'badge-official'"
        >{{ item.source }}</span>
      </div>
    </div>

    <h3 class="name">{{ item.display_name }}</h3>
    <p class="desc">{{ item.description }}</p>

    <div class="foot">
      <span class="provision">{{ summary }}</span>
      <span class="version">
        <span
          v-if="item.updated"
          class="updated"
          :title="`Updated ${shortDate(item.updated)}`"
        >{{ relativeDay(item.updated) }} · </span>v{{ item.version }}
      </span>
    </div>
  </RouterLink>
</template>

<style scoped>
.tpl {
  display: flex;
  flex-direction: column;
  gap: 10px;
  padding: 18px;
  color: inherit;
  transition: transform var(--transition), box-shadow var(--transition), border-color var(--transition);
}
.tpl:hover {
  color: inherit;
  transform: translateY(-2px);
  border-color: var(--primary-200);
  box-shadow: var(--shadow-md);
}
[data-theme='dark'] .tpl:hover {
  border-color: var(--primary-100);
}

.head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 10px;
}

.head-badges {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
  justify-content: flex-end;
}

.tile {
  width: 44px;
  height: 44px;
  font-size: 18px;
}

.name {
  font-size: 16px;
  font-weight: 600;
  line-height: 1.3;
}

.desc {
  flex: 1;
  font-size: 13.5px;
  line-height: 1.55;
  color: var(--text-secondary);
  /* Two lines keeps every card the same height, so the grid stays a grid. */
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.foot {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  padding-top: 10px;
  border-top: 1px solid var(--border-secondary);
  font-size: 12px;
  color: var(--text-tertiary);
}
.version {
  font-variant-numeric: tabular-nums;
  white-space: nowrap;
}
.updated {
  color: var(--text-muted);
}
</style>
