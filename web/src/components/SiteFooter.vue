<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { api } from '@/api/client'
import type { BuildInfo } from '@/api/types'

const year = new Date().getFullYear()
const siteUrl = 'https://miabi.io'
const docsUrl = 'https://docs.miabi.io'
const githubUrl = 'https://github.com/miabi-io/miabi'
const repoUrl = 'https://github.com/miabi-io/marketplace'

const build = ref<BuildInfo | null>(null)
const isRelease = computed(() => !!build.value && build.value.version !== 'dev')
const versionLabel = computed(() => (isRelease.value ? `v${build.value!.version}` : 'dev'))
const releaseUrl = computed(() => `${repoUrl}/releases/tag/${versionLabel.value}`)
const buildTitle = computed(() => {
  const b = build.value
  if (!b) return ''
  const parts = [b.commit_id !== 'unknown' ? `Commit ${b.commit_id}` : '']
  const date = new Date(b.build_date)
  if (!Number.isNaN(date.getTime())) parts.push(`built ${date.toLocaleDateString()}`)
  return parts.filter(Boolean).join(', ')
})

onMounted(async () => {
  try {
    build.value = await api.version()
  } catch {
    build.value = null
  }
})
</script>

<template>
  <footer class="site-footer">
    <div class="container">
      <div class="cols">
        <div class="brand-col">
          <a
            :href="siteUrl"
            class="brand"
          >
            <img
              src="/logo-white.svg"
              alt="Miabi"
              class="brand-mark"
            >
            <span>Miabi</span>
          </a>
          <p class="blurb">
            The open-source, self-hosted Platform-as-a-Service for Docker. Deploy apps with automatic
            SSL, managed databases, backups, scaling and monitoring — from one web interface.
          </p>
        </div>

        <div class="col">
          <h2>Marketplace</h2>
          <ul>
            <li><RouterLink to="/">
              Browse templates
            </RouterLink></li>
            <li><a
              :href="`${repoUrl}/blob/main/CONTRIBUTING.md`"
              target="_blank"
              rel="noopener"
            >Publish a template</a></li>
            <li><a
              href="/docs"
              target="_blank"
            >API reference</a></li>
            <li><a
              href="/v1/export"
              target="_blank"
            >Catalog export</a></li>
          </ul>
        </div>

        <div class="col">
          <h2>Project</h2>
          <ul>
            <li><a
              :href="docsUrl"
              target="_blank"
              rel="noopener"
            >Documentation</a></li>
            <li><a
              :href="githubUrl"
              target="_blank"
              rel="noopener"
            >GitHub</a></li>
            <li><a
              :href="repoUrl"
              target="_blank"
              rel="noopener"
            >Catalog repo</a></li>
          </ul>
        </div>
      </div>

      <div class="bottom">
        <p class="legal">
          <span>&copy; {{ year }} Miabi. Catalog content is Apache-2.0 licensed.</span>
          <template v-if="build">
            <a
              v-if="isRelease"
              :href="releaseUrl"
              class="version"
              :title="buildTitle || undefined"
              target="_blank"
              rel="noopener"
            >Catalog {{ versionLabel }}</a>
            <span
              v-else
              class="version"
              :title="buildTitle || undefined"
            >Catalog {{ versionLabel }}</span>
          </template>
        </p>
        <!-- This storefront runs on Miabi. -->
        <a
          :href="siteUrl"
          class="badge-link"
          aria-label="Running on Miabi"
          target="_blank"
          rel="noopener"
        >
          <img
            src="/badges/running-on-miabi-white.svg"
            alt="Running on Miabi"
            width="150"
            height="32"
          >
        </a>
      </div>
    </div>
  </footer>
</template>

<style scoped>
/* The footer is dark in both themes — it matches miabi.io, and the badge asset
   is the white-on-dark variant. */
.site-footer {
  margin-top: 64px;
  padding: 48px 0 32px;
  background: var(--bg-footer);
  color: #9ca3af;
  border-top: 1px solid var(--border-primary);
}

.cols {
  display: grid;
  grid-template-columns: 2fr 1fr 1fr;
  gap: 40px;
}

.brand {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  color: #fff;
  font-size: 18px;
  font-weight: 700;
  margin-bottom: 12px;
}
.brand-mark {
  width: 28px;
  height: 28px;
}

.blurb {
  font-size: 14px;
  line-height: 1.7;
  color: #6b7280;
  max-width: 42ch;
}

.col h2 {
  font-size: 13px;
  font-weight: 600;
  color: #fff;
  margin-bottom: 14px;
  text-transform: none;
}
.col ul {
  list-style: none;
  display: grid;
  gap: 10px;
}
.col a,
.col :deep(a) {
  font-size: 14px;
  color: #9ca3af;
}
.col a:hover,
.col :deep(a:hover) {
  color: var(--primary-400);
}

.bottom {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 20px;
  flex-wrap: wrap;
  margin-top: 36px;
  padding-top: 24px;
  border-top: 1px solid rgba(255, 255, 255, 0.08);
  font-size: 13px;
  color: #6b7280;
}

.legal {
  display: flex;
  align-items: center;
  gap: 8px 12px;
  flex-wrap: wrap;
}

.version {
  padding: 2px 8px;
  border: 1px solid rgba(255, 255, 255, 0.12);
  border-radius: 999px;
  font-size: 12px;
  font-variant-numeric: tabular-nums;
  color: #9ca3af;
  white-space: nowrap;
}
a.version:hover {
  color: var(--primary-400);
  border-color: var(--primary-400);
}

.badge-link {
  opacity: 0.85;
  transition: opacity var(--transition);
}
.badge-link:hover {
  opacity: 1;
}
.badge-link img {
  height: 32px;
  width: auto;
  display: block;
}

@media (max-width: 820px) {
  .cols {
    grid-template-columns: 1fr 1fr;
  }
  .brand-col {
    grid-column: 1 / -1;
  }
}
</style>
