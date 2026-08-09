<script setup lang="ts">
import { onMounted, onBeforeUnmount, ref } from 'vue'
import AppIcon from '@/components/AppIcon.vue'

const props = defineProps<{ modelValue: string; busy?: boolean; count?: number }>()
const emit = defineEmits<{ 'update:modelValue': [string] }>()

const input = ref<HTMLInputElement>()

// "/" and ⌘K focus the field from anywhere; Escape clears it. There is no submit
// button — results follow typing.
function onKeydown(e: KeyboardEvent) {
  const target = e.target as HTMLElement | null
  const typing = target && (target.tagName === 'INPUT' || target.tagName === 'TEXTAREA')
  if ((e.key === '/' && !typing) || ((e.metaKey || e.ctrlKey) && e.key.toLowerCase() === 'k')) {
    e.preventDefault()
    input.value?.focus()
    input.value?.select()
  }
}

function onEscape() {
  if (props.modelValue) emit('update:modelValue', '')
  else input.value?.blur()
}

onMounted(() => window.addEventListener('keydown', onKeydown))
onBeforeUnmount(() => window.removeEventListener('keydown', onKeydown))
</script>

<template>
  <div
    class="search"
    :class="{ busy }"
  >
    <AppIcon
      name="search"
      :size="20"
      class="lead"
    />
    <input
      ref="input"
      type="search"
      class="input"
      :value="modelValue"
      placeholder="Search templates, tags or descriptions…"
      autocomplete="off"
      spellcheck="false"
      aria-label="Search templates"
      aria-describedby="search-hint"
      @input="emit('update:modelValue', ($event.target as HTMLInputElement).value)"
      @keydown.esc.prevent="onEscape"
    >
    <span
      v-if="busy"
      class="spinner"
      aria-hidden="true"
    />
    <button
      v-else-if="modelValue"
      type="button"
      class="clear"
      aria-label="Clear search"
      @click="emit('update:modelValue', '')"
    >
      <AppIcon name="close" />
    </button>
    <kbd
      v-else
      class="hint"
      aria-hidden="true"
    >/</kbd>
  </div>
  <p
    id="search-hint"
    class="sr-only"
  >
    Results update as you type.
  </p>
</template>

<style scoped>
.search {
  position: relative;
  display: flex;
  align-items: center;
  background: var(--bg-input);
  border: 1px solid var(--border-input);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-sm);
  transition: border-color var(--transition), box-shadow var(--transition);
}
.search:focus-within {
  border-color: var(--border-focus);
  box-shadow: var(--shadow-focus);
}

.lead {
  padding-left: 16px;
  font-size: 20px;
  color: var(--text-muted);
}

.input {
  flex: 1;
  min-width: 0;
  padding: 14px 12px;
  border: 0;
  background: transparent;
  color: var(--text-primary);
  font-family: inherit;
  font-size: 16px;
}
.input:focus {
  outline: none;
}
.input::placeholder {
  color: var(--text-muted);
}
.input::-webkit-search-cancel-button {
  display: none;
}

.clear {
  display: grid;
  place-items: center;
  width: 30px;
  height: 30px;
  margin-right: 10px;
  border: 0;
  border-radius: 50%;
  background: var(--bg-tertiary);
  color: var(--text-tertiary);
  cursor: pointer;
}
.clear:hover {
  background: var(--bg-hover);
  color: var(--text-primary);
}

.hint {
  margin-right: 14px;
  padding: 2px 8px;
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-sm);
  background: var(--bg-tertiary);
  color: var(--text-muted);
  font-family: inherit;
  font-size: 12px;
}

.spinner {
  width: 16px;
  height: 16px;
  margin-right: 16px;
  border: 2px solid var(--border-primary);
  border-top-color: var(--primary-500);
  border-radius: 50%;
  animation: spin 700ms linear infinite;
}
@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}
</style>
