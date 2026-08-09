<script setup lang="ts">
import { computed } from 'vue'
import { marked } from 'marked'
import DOMPurify from 'dompurify'

const props = defineProps<{ source: string }>()

// READMEs are community-contributed, so the rendered HTML is sanitized before it
// reaches the DOM — never trust catalog content with v-html.
const html = computed(() => {
  const raw = marked.parse(props.source, { async: false, gfm: true, breaks: false }) as string
  return DOMPurify.sanitize(raw, { USE_PROFILES: { html: true } })
})
</script>

<template>
  <!-- eslint-disable-next-line vue/no-v-html -- sanitized above -->
  <div
    class="md"
    v-html="html"
  />
</template>

<style scoped>
.md {
  font-size: 14.5px;
  line-height: 1.7;
  color: var(--text-secondary);
  overflow-wrap: anywhere;
}
.md :deep(h1),
.md :deep(h2),
.md :deep(h3) {
  color: var(--text-primary);
  line-height: 1.3;
  margin: 24px 0 10px;
}
.md :deep(h1) {
  font-size: 20px;
}
.md :deep(h2) {
  font-size: 17px;
}
.md :deep(h3) {
  font-size: 15px;
}
.md :deep(p),
.md :deep(ul),
.md :deep(ol),
.md :deep(blockquote),
.md :deep(table) {
  margin-bottom: 14px;
}
.md :deep(ul),
.md :deep(ol) {
  padding-left: 22px;
}
.md :deep(li) {
  margin-bottom: 4px;
}
.md :deep(code) {
  padding: 2px 6px;
  border-radius: var(--radius-sm);
  background: var(--bg-tertiary);
  font-size: 0.9em;
}
.md :deep(pre) {
  padding: 14px 16px;
  border-radius: var(--radius);
  background: var(--bg-code);
  color: #e2e0f5;
  overflow-x: auto;
  margin-bottom: 14px;
}
.md :deep(pre code) {
  background: transparent;
  padding: 0;
  color: inherit;
}
.md :deep(blockquote) {
  padding-left: 14px;
  border-left: 3px solid var(--border-primary);
  color: var(--text-tertiary);
}
.md :deep(img) {
  border-radius: var(--radius);
}
.md :deep(table) {
  width: 100%;
  border-collapse: collapse;
  display: block;
  overflow-x: auto;
}
.md :deep(th),
.md :deep(td) {
  padding: 8px 12px;
  border: 1px solid var(--border-primary);
  text-align: left;
}
</style>
