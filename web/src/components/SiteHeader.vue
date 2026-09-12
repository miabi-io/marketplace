<script setup lang="ts">
import { useThemeStore } from '@/stores/theme'
import AppIcon from '@/components/AppIcon.vue'

const theme = useThemeStore()
const docsUrl = 'https://docs.miabi.io'
const repoUrl = 'https://github.com/miabi-io/marketplace'
</script>

<template>
  <header class="site-header">
    <div class="container bar">
      <RouterLink
        to="/"
        class="brand"
      >
        <img
          src="/logo.svg"
          alt=""
          class="brand-mark"
          aria-hidden="true"
        >
        <span class="brand-text">
          <strong>Miabi</strong> Marketplace
        </span>
      </RouterLink>

      <nav class="links">
        <a
          :href="`${docsUrl}/marketplace/overview`"
          target="_blank"
          rel="noopener"
        >Docs</a>
        <a
          :href="`${repoUrl}/blob/main/CONTRIBUTING.md`"
          target="_blank"
          rel="noopener"
        >Contribute</a>
        <a
          href="/docs"
          target="_blank"
        >API</a>
        <RouterLink to="/about">
          About
        </RouterLink>
        <button
          class="btn btn-secondary btn-icon"
          type="button"
          :aria-label="theme.isDark ? 'Switch to light theme' : 'Switch to dark theme'"
          :title="theme.isDark ? 'Light theme' : 'Dark theme'"
          @click="theme.toggle()"
        >
          <AppIcon :name="theme.isDark ? 'sun' : 'moon'" />
        </button>
      </nav>
    </div>
  </header>
</template>

<style scoped>
.site-header {
  position: sticky;
  top: 0;
  z-index: 20;
  background: color-mix(in srgb, var(--bg-primary) 85%, transparent);
  backdrop-filter: blur(10px);
  border-bottom: 1px solid var(--border-primary);
}

.bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  height: 60px;
}

.brand {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  color: var(--text-primary);
  font-size: 16px;
}
.brand:hover {
  color: var(--text-primary);
}
.brand-mark {
  width: 26px;
  height: 26px;
}
.brand-text strong {
  font-weight: 700;
}

.links {
  display: flex;
  align-items: center;
  gap: 6px;
}
.links a {
  padding: 8px 12px;
  border-radius: var(--radius);
  color: var(--text-secondary);
  font-size: 14px;
  font-weight: 500;
}
.links a:hover,
.links a.router-link-exact-active {
  background: var(--bg-hover);
  color: var(--text-primary);
}

@media (max-width: 640px) {
  .links a {
    display: none;
  }
}
</style>
