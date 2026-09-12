<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { api } from '@/api/client'
import type { BuildInfo, Page } from '@/api/types'
import AppIcon from '@/components/AppIcon.vue'
import CopyButton from '@/components/CopyButton.vue'

const repoUrl = 'https://github.com/miabi-io/marketplace'
const contributingUrl = `${repoUrl}/blob/main/CONTRIBUTING.md`

// The origin, not a hard-coded marketplace.miabi.io: a self-hosted mirror should
// tell its readers to point Miabi at itself.
const envLine = `MIABI_MARKETPLACE_URL=${window.location.origin}`

const links = [
  { label: 'Source repository', icon: 'github', href: repoUrl },
  { label: 'Contributing guide', icon: 'source-branch', href: contributingUrl },
  { label: 'API reference', icon: 'api', href: '/docs' },
  { label: 'Catalog export', icon: 'download', href: '/v1/export' },
  { label: 'Report an issue', icon: 'bug', href: `${repoUrl}/issues` },
  { label: 'Miabi documentation', icon: 'book', href: 'https://docs.miabi.io' },
]

const build = ref<BuildInfo | null>(null)
const stats = ref<{ total?: number; official?: number; community?: number; categories?: number }>({})

const statCards = computed(() => [
  { label: 'Templates', value: stats.value.total },
  { label: 'Official', value: stats.value.official },
  { label: 'Community', value: stats.value.community },
  { label: 'Categories', value: stats.value.categories },
])

const isRelease = computed(() => !!build.value && build.value.version !== 'dev')
const versionLabel = computed(() => {
  if (!build.value) return '—'
  return isRelease.value ? `v${build.value.version}` : 'dev'
})
const releaseUrl = computed(() => `${repoUrl}/releases/tag/${versionLabel.value}`)
const commit = computed(() => (build.value && build.value.commit_id !== 'unknown' ? build.value.commit_id : ''))
const builtAt = computed(() => {
  const d = new Date(build.value?.build_date ?? '')
  if (Number.isNaN(d.getTime())) return null
  return {
    label: d.toLocaleDateString(undefined, { year: 'numeric', month: 'long', day: 'numeric' }),
    exact: d.toLocaleString(),
  }
})
const versionText = computed(() =>
  [
    'Miabi Marketplace',
    versionLabel.value,
    commit.value ? `(${commit.value})` : '',
    builtAt.value ? `built ${builtAt.value.exact}` : '',
  ]
    .filter(Boolean)
    .join(' '),
)

const totalOf = (r: PromiseSettledResult<Page>) => (r.status === 'fulfilled' ? r.value.total : undefined)

onMounted(async () => {
  document.title = 'About — Miabi Marketplace'
  const [info, all, official, community, categories] = await Promise.allSettled([
    api.version(),
    api.search({ per_page: 1 }),
    api.search({ source: 'official', per_page: 1 }),
    api.search({ source: 'community', per_page: 1 }),
    api.categories(),
  ])
  if (info.status === 'fulfilled') build.value = info.value
  stats.value = {
    total: totalOf(all),
    official: totalOf(official),
    community: totalOf(community),
    categories: categories.status === 'fulfilled' ? categories.value.length : undefined,
  }
})
</script>

<template>
  <section class="head">
    <div class="container">
      <p class="eyebrow">
        About
      </p>
      <h1>The Miabi Marketplace</h1>
      <p class="lead">
        The registry and storefront for Miabi templates: apps, databases and full stacks you install in
        one click on Miabi, the open-source PaaS for Docker.
      </p>
      <div class="actions">
        <RouterLink
          class="btn btn-primary"
          to="/"
        >
          <AppIcon name="search" />
          Browse templates
        </RouterLink>
        <a
          class="btn btn-secondary"
          :href="contributingUrl"
          target="_blank"
          rel="noopener"
        >
          <AppIcon name="source-branch" />
          Publish a template
        </a>
      </div>
    </div>
  </section>

  <div class="container body">
    <div class="main">
      <dl class="card stats">
        <div
          v-for="s in statCards"
          :key="s.label"
          class="stat"
        >
          <dt>{{ s.label }}</dt>
          <dd>{{ s.value ?? '—' }}</dd>
        </div>
      </dl>

      <section class="card panel">
        <h2 class="section-title">
          What it is
        </h2>
        <p class="prose">
          Every template is a <code>miabi.io/v1</code> manifest describing a stack: its applications,
          databases, volumes and configuration. Miabi turns it into a running deployment. It asks for the
          inputs, generates the secrets, provisions the databases and wires them to the apps.
        </p>
        <p class="prose">
          Git is the database. Templates live in the
          <a
            :href="repoUrl"
            target="_blank"
            rel="noopener"
          >catalog repository</a>, and each tagged release publishes a new catalog version, the one
          shown below and in the footer.
        </p>
      </section>

      <section class="card panel">
        <h2 class="section-title">
          Official and community
        </h2>
        <div class="sources">
          <div>
            <RouterLink
              to="/?source=official"
              class="badge badge-official"
            >
              Official
            </RouterLink>
            <p>
              Curated by the Miabi maintainers. Every change needs core-maintainer review, and Miabi ships
              these templates inside its binary as an offline fallback.
            </p>
          </div>
          <div>
            <RouterLink
              to="/?source=community"
              class="badge badge-community"
            >
              Community
            </RouterLink>
            <p>
              Contributed through pull requests and reviewed by a maintainer before merge. A community
              template can be promoted to official.
            </p>
          </div>
        </div>
      </section>

      <section class="card panel">
        <h2 class="section-title">
          How templates are checked
        </h2>
        <ul class="checks">
          <li>
            <AppIcon name="shield" />
            <span>Community manifests are untrusted input: the validator rejects host bind mounts, privileged
              flags, unknown fields and malformed values.</span>
          </li>
          <li>
            <AppIcon name="check" />
            <span>CI validates every manifest against the
              <a
                href="/schema/template.schema.json"
                target="_blank"
              >template schema</a> and keeps names unique across official and community.</span>
          </li>
          <li>
            <AppIcon name="tag" />
            <span>Each version carries a content digest, so a published version cannot change silently.
              Image tags are pinned.</span>
          </li>
        </ul>
      </section>

      <section class="card panel">
        <h2 class="section-title">
          Use it from Miabi
        </h2>
        <p class="prose">
          Miabi syncs this catalog and lists it in the console's marketplace. To point an install at this
          registry, set:
        </p>
        <div class="env">
          <pre class="code-block"><code>{{ envLine }}</code></pre>
          <CopyButton :value="envLine" />
        </div>
        <p class="prose muted">
          Air-gapped or private? Fork the repository and sync from its committed <code>export.json</code>,
          or run this service on your own domain.
        </p>
      </section>
    </div>

    <aside class="side">
      <section class="card panel">
        <h2 class="section-title">
          Catalog version
        </h2>
        <dl class="kv">
          <div>
            <dt>Version</dt>
            <dd>
              <a
                v-if="isRelease"
                :href="releaseUrl"
                target="_blank"
                rel="noopener"
              >{{ versionLabel }}</a>
              <template v-else>
                {{ versionLabel }}
              </template>
            </dd>
          </div>
          <div>
            <dt>Commit</dt>
            <dd class="mono">
              <a
                v-if="commit"
                :href="`${repoUrl}/commit/${commit}`"
                target="_blank"
                rel="noopener"
              >{{ commit }}</a>
              <template v-else>
                —
              </template>
            </dd>
          </div>
          <div>
            <dt>Built</dt>
            <dd :title="builtAt?.exact">
              {{ builtAt?.label ?? '—' }}
            </dd>
          </div>
        </dl>
        <CopyButton
          v-if="build"
          :value="versionText"
          label="Copy version"
        />
      </section>

      <section class="card panel">
        <h2 class="section-title">
          Links
        </h2>
        <ul class="links">
          <li
            v-for="l in links"
            :key="l.href"
          >
            <a
              :href="l.href"
              target="_blank"
              rel="noopener"
            >
              <AppIcon :name="l.icon" />
              <span>{{ l.label }}</span>
              <AppIcon
                name="external"
                :size="14"
                class="ext"
              />
            </a>
          </li>
        </ul>
      </section>
    </aside>
  </div>
</template>

<style scoped>
.head {
  padding: 40px 0 32px;
  background: var(--bg-primary);
  border-bottom: 1px solid var(--border-primary);
}
.eyebrow {
  font-size: 13px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.04em;
  color: var(--primary-600);
  margin-bottom: 6px;
}
.head h1 {
  font-size: 30px;
  font-weight: 700;
  letter-spacing: -0.01em;
}
.lead {
  margin-top: 8px;
  max-width: 64ch;
  font-size: 16px;
  color: var(--text-secondary);
}
.actions {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  margin-top: 20px;
}

.body {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 300px;
  gap: 20px;
  padding-top: 24px;
}
.main,
.side {
  display: grid;
  grid-template-columns: minmax(0, 1fr);
  gap: 16px;
  align-content: start;
}
.panel {
  padding: 18px;
}

.prose {
  font-size: 14px;
  line-height: 1.7;
  color: var(--text-secondary);
}
.prose + .prose,
.prose + .env,
.env + .prose {
  margin-top: 12px;
}
.muted {
  color: var(--text-tertiary);
}
.prose code {
  padding: 1px 5px;
  border-radius: var(--radius-sm);
  background: var(--bg-tertiary);
  font-size: 13px;
}

.stats {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
}
.stat {
  padding: 16px 18px;
}
.stat + .stat {
  border-left: 1px solid var(--border-primary);
}
.stat dt {
  font-size: 13px;
  color: var(--text-tertiary);
}
.stat dd {
  font-size: 24px;
  font-weight: 700;
  font-variant-numeric: tabular-nums;
  color: var(--text-primary);
}

.sources {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
}
.sources p {
  margin-top: 8px;
  font-size: 14px;
  line-height: 1.6;
  color: var(--text-secondary);
}

.checks {
  list-style: none;
  display: grid;
  gap: 12px;
}
.checks li {
  display: flex;
  gap: 10px;
  font-size: 14px;
  line-height: 1.6;
  color: var(--text-secondary);
}
.checks .icon {
  margin-top: 3px;
  color: var(--primary-600);
}

.env {
  display: flex;
  align-items: center;
  gap: 8px;
}
.env pre {
  flex: 1;
  min-width: 0;
}

.kv {
  display: grid;
  gap: 10px;
  margin-bottom: 16px;
}
.kv div {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  font-size: 14px;
}
.kv dt {
  color: var(--text-tertiary);
}
.kv dd {
  font-weight: 500;
  text-align: right;
  color: var(--text-primary);
}
.mono {
  font-family: 'SF Mono', ui-monospace, Menlo, Consolas, monospace;
  font-size: 13px;
}

.links {
  list-style: none;
  display: grid;
  gap: 2px;
}
.links a {
  display: flex;
  align-items: center;
  gap: 10px;
  margin: 0 -8px;
  padding: 7px 8px;
  border-radius: var(--radius);
  font-size: 14px;
  color: var(--text-secondary);
}
.links a:hover {
  background: var(--bg-hover);
  color: var(--text-primary);
}
.links span {
  flex: 1;
}
.ext {
  color: var(--text-muted);
}

@media (max-width: 880px) {
  .body {
    grid-template-columns: minmax(0, 1fr);
  }
}

@media (max-width: 640px) {
  .head h1 {
    font-size: 24px;
  }
  .stats {
    grid-template-columns: 1fr 1fr;
  }
  .stat:nth-child(3) {
    border-left: 0;
  }
  .stat:nth-child(n + 3) {
    border-top: 1px solid var(--border-primary);
  }
  .sources {
    grid-template-columns: 1fr;
  }
}
</style>
