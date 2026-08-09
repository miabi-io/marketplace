<script setup lang="ts">
import { ref } from 'vue'
import AppIcon from '@/components/AppIcon.vue'

const props = defineProps<{ value: string; label?: string }>()
const copied = ref(false)
let timer: ReturnType<typeof setTimeout> | undefined

async function copy() {
  try {
    await navigator.clipboard.writeText(props.value)
  } catch {
    // Clipboard API needs a secure context; fall back to a selectable prompt
    // rather than failing silently.
    const ta = document.createElement('textarea')
    ta.value = props.value
    document.body.appendChild(ta)
    ta.select()
    document.execCommand('copy')
    ta.remove()
  }
  copied.value = true
  if (timer) clearTimeout(timer)
  timer = setTimeout(() => (copied.value = false), 1600)
}
</script>

<template>
  <button
    type="button"
    class="btn btn-secondary btn-sm"
    :aria-label="label ?? 'Copy to clipboard'"
    @click="copy"
  >
    <AppIcon :name="copied ? 'check' : 'copy'" />
    <span>{{ copied ? 'Copied' : (label ?? 'Copy') }}</span>
  </button>
</template>
