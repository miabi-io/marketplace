<script setup lang="ts">
import AppIcon from '@/components/AppIcon.vue'

defineProps<{ icon?: string; title: string; hint?: string }>()
</script>

<template>
  <div class="empty">
    <AppIcon
      :name="icon ?? 'package'"
      :size="40"
      class="glyph"
    />
    <p class="title">{{ title }}</p>
    <p
      v-if="hint"
      class="hint"
    >{{ hint }}</p>
    <slot />
  </div>
</template>

<style scoped>
.empty {
  display: grid;
  justify-items: center;
  gap: 8px;
  padding: 64px 24px;
  text-align: center;
}
.glyph {
  color: var(--text-muted);
}
.title {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
}
.hint {
  font-size: 14px;
  color: var(--text-tertiary);
  max-width: 46ch;
}
</style>
