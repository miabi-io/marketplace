<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { api, ApiError } from '@/api/client'
import type { TemplateDetail, VersionRef } from '@/api/types'
import { engineLabel, hueOf, isURL, monogram, provision, shortDate } from '@/utils/format'
import CopyButton from '@/components/CopyButton.vue'
import MarkdownView from '@/components/MarkdownView.vue'
import EmptyState from '@/components/EmptyState.vue'
import AppIcon from '@/components/AppIcon.vue'

const route = useRoute()
const name = computed(() => String(route.params.name ?? ''))

const detail = ref<TemplateDetail | null>(null)
const yaml = ref('')
const loading = ref(true)
const notFound = ref(false)
const error = ref('')
const version = ref('')
const showYAML = ref(false)
// See TemplateCard: third-party icon URLs can fail, so fall back to the monogram.
const iconFailed = ref(false)

const hue = computed(() => hueOf(name.value))
const entry = computed(() => detail.value?.entry)
const manifest = computed(() => detail.value?.manifest)
interface TagChange {
  app: string
  from?: string
  to: string
}

// versions arrive newest first, so each entry's predecessor is the next one in the list.
function imageChanges(v: VersionRef, older?: VersionRef): TagChange[] {
  if (!older?.images || !v.images) return []
  const before = new Map(older.images.map((i) => [i.app, i.tag]))
  return v.images
    .filter((i) => before.get(i.app) !== i.tag)
    .map((i) => ({ app: i.app, from: before.get(i.app), to: i.tag }))
}

const history = computed(() =>
  (detail.value?.versions ?? []).map((v, i, all) => ({ ...v, changes: imageChanges(v, all[i + 1]) })),
)
const hasHistory = computed(() =>
  history.value.some((v) => v.release?.notes || v.release?.released || v.changes.length),
)

const freshnessHint = computed(() => {
  const f = entry.value?.freshness
  if (!f || (f.status !== 'major' && f.status !== 'update')) return ''
  const checked = f.checked_at ? ` (checked ${shortDate(f.checked_at)})` : ''
  if (f.status === 'major') {
    const app = f.apps.find((a) => a.status === 'major')
    const latest = app?.latest ? `: ${app.image} ${app.latest}` : ''
    return `Upstream has a newer major release${latest}${checked}.`
  }
  return `A newer upstream release is available${checked}.`
})

const configFileCount = computed(() =>
  (manifest.value?.configs ?? []).reduce((n, c) => n + Object.keys(c.files ?? {}).length, 0),
)

async function load(v?: string) {
  loading.value = true
  error.value = ''
  notFound.value = false
  try {
    const d = await api.template(name.value, v)
    detail.value = d
    version.value = d.manifest.metadata.version
    document.title = `${d.entry.display_name} — Miabi Marketplace`
    yaml.value = await api.manifestYAML(name.value, version.value)
  } catch (e: unknown) {
    if (e instanceof ApiError && e.status === 404) notFound.value = true
    else error.value = e instanceof Error ? e.message : 'Failed to load template'
  } finally {
    loading.value = false
  }
}

watch(name, () => void load(), { immediate: true })

function selectVersion(v: string) {
  if (v !== version.value) void load(v)
}

const manifestURL = computed(() => (version.value ? api.manifestURL(name.value, version.value) : ''))
const absoluteManifestURL = computed(() =>
  manifestURL.value ? new URL(manifestURL.value, window.location.origin).toString() : '',
)
</script>

<template>
  <div
    v-if="loading"
    class="container state"
  >
    <p class="muted">Loading…</p>
  </div>

  <div
    v-else-if="notFound"
    class="container"
  >
    <EmptyState
      icon="package-gone"
      title="Template not found"
      hint="It may have been renamed or removed from the catalog."
    >
      <RouterLink
        class="btn btn-primary btn-sm"
        to="/"
      >
        Browse the catalog
      </RouterLink>
    </EmptyState>
  </div>

  <div
    v-else-if="error"
    class="container state"
  >
    <p
      class="muted"
      role="alert"
    >{{ error }}</p>
  </div>

  <template v-else-if="detail && entry && manifest">
    <section class="head">
      <div class="container">
        <RouterLink
          class="back"
          to="/"
        >
          <AppIcon name="arrow-left" /> All templates
        </RouterLink>

        <div class="hero">
          <span
            class="tile"
            :style="{ '--tile-hue': hue }"
            aria-hidden="true"
          >
            <img
              v-if="isURL(entry.icon) && !iconFailed"
              :src="entry.icon"
              alt=""
              @error="iconFailed = true"
            >
            <template v-else>{{ monogram(entry.display_name) }}</template>
          </span>

          <div class="titles">
            <h1>{{ entry.display_name }}</h1>
            <p class="desc">{{ entry.description }}</p>
            <div class="badges">
              <span
                class="badge"
                :class="entry.source === 'community' ? 'badge-community' : 'badge-official'"
              >{{ entry.source }}</span>
              <span
                v-if="entry.category"
                class="badge badge-neutral"
              >{{ entry.category }}</span>
              <span class="badge badge-neutral">{{ provision(entry) }}</span>
              <span
                v-if="entry.deprecated"
                class="badge badge-warning"
              >deprecated</span>
            </div>
            <p
              v-if="entry.updated || freshnessHint"
              class="meta-line"
            >
              <span v-if="entry.updated">Last updated {{ shortDate(entry.updated) }}</span>
              <span
                v-if="freshnessHint"
                class="fresh"
              >{{ freshnessHint }}</span>
            </p>
          </div>

          <div class="actions">
            <a
              class="btn btn-primary"
              :href="manifestURL"
              download
            >
              <AppIcon name="download" />
              template.yaml
            </a>
            <CopyButton
              :value="absoluteManifestURL"
              label="Copy manifest URL"
            />
          </div>
        </div>
      </div>
    </section>

    <section class="container body">
      <div class="main">
        <div
          v-if="entry.deprecated"
          class="deprecated"
          role="status"
        >
          <AppIcon name="alert" />
          <div>
            <p class="deprecated-title">This template is deprecated</p>
            <p v-if="entry.deprecation_note">{{ entry.deprecation_note }}</p>
            <p v-if="entry.replaced_by">
              Use
              <RouterLink :to="{ name: 'template', params: { name: entry.replaced_by } }">{{ entry.replaced_by }}</RouterLink>
              instead.
            </p>
          </div>
        </div>

        <div class="card panel install">
          <h2 class="section-title">Install</h2>
          <ol class="steps">
            <li>Open your Miabi workspace and go to <strong>Marketplace</strong>.</li>
            <li>Search for <strong>{{ entry.display_name }}</strong> and choose <strong>Install</strong>.</li>
            <li>Answer the inputs below — managed databases and volumes are created for you.</li>
          </ol>
          <p class="note">
            Self-hosting this catalog or installing by hand? Import the manifest instead:
          </p>
          <div class="cmd">
            <code>curl -fsSL {{ absoluteManifestURL }}</code>
            <CopyButton
              :value="`curl -fsSL ${absoluteManifestURL}`"
              label="Copy"
            />
          </div>
        </div>

        <div
          v-if="manifest.inputs?.length"
          class="card panel"
        >
          <h2 class="section-title">Inputs</h2>
          <ul class="inputs">
            <li
              v-for="i in manifest.inputs"
              :key="i.key"
            >
              <div class="input-head">
                <code>{{ i.key }}</code>
                <span
                  v-if="i.required"
                  class="req"
                >required</span>
                <span
                  v-if="i.generate"
                  class="gen"
                >auto-generated</span>
              </div>
              <p class="input-label">{{ i.label ?? i.key }}</p>
              <p
                v-if="i.help"
                class="input-help"
              >{{ i.help }}</p>
            </li>
          </ul>
        </div>

        <div
          v-if="detail.readme"
          class="card panel"
        >
          <h2 class="section-title">Readme</h2>
          <MarkdownView :source="detail.readme" />
        </div>

        <div
          v-if="hasHistory"
          class="card panel"
        >
          <h2 class="section-title">Version history</h2>
          <ol class="history">
            <li
              v-for="v in history"
              :key="v.version"
              class="rel"
            >
              <div class="rel-head">
                <button
                  type="button"
                  class="rel-version"
                  :class="{ active: v.version === version }"
                  @click="selectVersion(v.version)"
                >
                  v{{ v.version }}
                </button>
                <span
                  v-if="v.release?.released"
                  class="rel-date"
                >{{ shortDate(v.release.released) }}</span>
                <span
                  v-if="v.release?.breaking"
                  class="badge badge-warning"
                >breaking</span>
                <span
                  v-if="v.release?.deprecated"
                  class="badge badge-neutral"
                >deprecated</span>
              </div>
              <p
                v-if="v.release?.requires_from"
                class="rel-requires"
              >
                Requires upgrading from v{{ v.release.requires_from }} or later first.
              </p>
              <p
                v-if="v.release?.deprecation_note"
                class="rel-requires"
              >
                {{ v.release.deprecation_note }}
              </p>
              <MarkdownView
                v-if="v.release?.notes"
                class="rel-notes"
                :source="v.release.notes"
              />
              <ul
                v-else-if="v.changes.length"
                class="rel-changes"
              >
                <li
                  v-for="c in v.changes"
                  :key="c.app"
                >
                  <code>{{ c.app }}</code>: {{ c.from ?? 'new' }} → {{ c.to }}
                </li>
              </ul>
            </li>
          </ol>
        </div>

        <div class="card panel">
          <div class="panel-head">
            <h2 class="section-title">Manifest</h2>
            <div class="panel-actions">
              <CopyButton
                :value="yaml"
                label="Copy YAML"
              />
              <button
                type="button"
                class="btn btn-secondary btn-sm"
                :aria-expanded="showYAML"
                @click="showYAML = !showYAML"
              >
                {{ showYAML ? 'Hide' : 'Show' }}
              </button>
            </div>
          </div>
          <pre
            v-if="showYAML"
            class="code-block yaml"
          >{{ yaml }}</pre>
        </div>
      </div>

      <aside class="side">
        <div class="card panel">
          <h2 class="section-title">Version</h2>
          <div class="versions">
            <button
              v-for="v in detail.versions"
              :key="v.version"
              type="button"
              class="ver"
              :class="{ active: v.version === version }"
              @click="selectVersion(v.version)"
            >
              <span>v{{ v.version }}</span>
              <span
                v-if="v.version === detail.versions[0].version"
                class="tag"
              >latest</span>
            </button>
          </div>
          <p class="digest">
            <span>Digest</span>
            <code>{{ detail.versions.find((v) => v.version === version)?.digest?.slice(0, 19) }}…</code>
          </p>
        </div>

        <div class="card panel">
          <h2 class="section-title">What it provisions</h2>
          <ul class="facts">
            <li>
              <AppIcon name="cube" />
              <span>{{ entry.applications }} application{{ entry.applications === 1 ? '' : 's' }}</span>
            </li>
            <li v-if="manifest.databases?.length">
              <AppIcon name="database" />
              <span>{{ manifest.databases.map((d) => engineLabel(d.engine)).join(', ') }}</span>
            </li>
            <li v-if="entry.volumes">
              <AppIcon name="harddisk" />
              <span>{{ entry.volumes }} volume{{ entry.volumes === 1 ? '' : 's' }}</span>
            </li>
            <li v-if="manifest.configs?.length">
              <AppIcon name="file" />
              <span>{{ configFileCount }} config file{{ configFileCount === 1 ? '' : 's' }}</span>
            </li>
            <li v-if="manifest.metadata.minMiabi">
              <AppIcon name="tag" />
              <span>Miabi {{ manifest.metadata.minMiabi }}+</span>
            </li>
          </ul>
        </div>

        <div class="card panel">
          <h2 class="section-title">Links</h2>
          <ul class="links">
            <li v-if="entry.homepage">
              <a
                :href="entry.homepage"
                target="_blank"
                rel="noopener"
              >
                <AppIcon name="external" /> Project homepage
              </a>
            </li>
            <li v-if="detail.meta.source_repo">
              <a
                :href="detail.meta.source_repo"
                target="_blank"
                rel="noopener"
              >
                <AppIcon name="source-branch" /> Template source
              </a>
            </li>
            <li v-if="entry.author?.website">
              <a
                :href="entry.author.website"
                target="_blank"
                rel="noopener"
              >
                <AppIcon name="account" /> {{ entry.author.name }}
              </a>
            </li>
          </ul>
        </div>

        <div
          v-if="entry.tags?.length"
          class="card panel"
        >
          <h2 class="section-title">Tags</h2>
          <div class="tags">
            <RouterLink
              v-for="t in entry.tags"
              :key="t"
              class="tag-link"
              :to="{ name: 'home', query: { tag: t } }"
            >
              {{ t }}
            </RouterLink>
          </div>
        </div>
      </aside>
    </section>
  </template>
</template>

<style scoped>
.state {
  padding: 80px 20px;
  text-align: center;
}
.muted {
  color: var(--text-tertiary);
}

.head {
  padding: 20px 0 28px;
  background: var(--bg-primary);
  border-bottom: 1px solid var(--border-primary);
}
.back {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 14px;
  color: var(--text-tertiary);
  margin-bottom: 18px;
}
.back:hover {
  color: var(--primary-600);
}

.hero {
  display: flex;
  align-items: flex-start;
  gap: 18px;
  flex-wrap: wrap;
}
.hero .tile {
  width: 66px;
  height: 66px;
  font-size: 26px;
}
.titles {
  flex: 1;
  min-width: 240px;
}
.titles h1 {
  font-size: 26px;
  font-weight: 700;
  letter-spacing: -0.01em;
}
.desc {
  margin-top: 4px;
  color: var(--text-secondary);
  max-width: 70ch;
}
.badges {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 12px;
}
.meta-line {
  display: flex;
  flex-wrap: wrap;
  gap: 4px 14px;
  margin-top: 10px;
  font-size: 13px;
  color: var(--text-tertiary);
}
.fresh {
  color: var(--text-muted);
}
.deprecated {
  display: flex;
  gap: 12px;
  padding: 14px 16px;
  border: 1px solid var(--warning-200);
  border-radius: var(--radius-lg);
  background: var(--warning-50);
  color: var(--text-secondary);
  font-size: 14px;
}
.deprecated > .icon {
  flex-shrink: 0;
  margin-top: 2px;
  color: var(--warning-700);
}
.deprecated-title {
  font-weight: 600;
  color: var(--text-primary);
}
.history {
  list-style: none;
  display: grid;
}
.rel {
  padding: 14px 0;
  border-top: 1px solid var(--border-secondary);
}
.rel:first-child {
  padding-top: 0;
  border-top: none;
}
.rel-head {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 8px;
}
.rel-version {
  padding: 0;
  border: none;
  background: none;
  font-family: inherit;
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
  cursor: pointer;
}
.rel-version:hover,
.rel-version.active {
  color: var(--primary-600);
}
.rel-date {
  font-size: 12.5px;
  color: var(--text-tertiary);
}
.rel-requires {
  margin-top: 6px;
  font-size: 13px;
  color: var(--text-secondary);
}
.rel-notes {
  margin-top: 8px;
  font-size: 13.5px;
}
.rel-changes {
  margin-top: 8px;
  padding-left: 18px;
  font-size: 13px;
  color: var(--text-secondary);
}
.rel-changes code {
  font-size: 12.5px;
}
.actions {
  display: flex;
  gap: 8px;
  align-items: center;
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
  gap: 16px;
  align-content: start;
}
.panel {
  padding: 18px;
}
.panel-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}
.panel-head .section-title {
  margin-bottom: 0;
}
.panel-actions {
  display: flex;
  gap: 8px;
}

.steps {
  padding-left: 20px;
  display: grid;
  gap: 6px;
  font-size: 14px;
  color: var(--text-secondary);
}
.note {
  margin-top: 14px;
  font-size: 13px;
  color: var(--text-tertiary);
}
.cmd {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-top: 8px;
}
.cmd code {
  flex: 1;
  min-width: 0;
  padding: 10px 12px;
  border-radius: var(--radius);
  background: var(--bg-code);
  color: #e2e0f5;
  font-size: 12.5px;
  overflow-x: auto;
  white-space: nowrap;
}

.yaml {
  margin-top: 12px;
  max-height: 480px;
  overflow: auto;
  white-space: pre;
}

.inputs {
  list-style: none;
  display: grid;
  gap: 14px;
}
.input-head {
  display: flex;
  align-items: center;
  gap: 8px;
}
.input-head code {
  font-size: 13px;
  color: var(--text-primary);
}
.req,
.gen {
  font-size: 11px;
  padding: 1px 7px;
  border-radius: 999px;
  border: 1px solid var(--border-primary);
  color: var(--text-tertiary);
}
.input-label {
  font-size: 14px;
  color: var(--text-primary);
}
.input-help {
  font-size: 13px;
  color: var(--text-tertiary);
}

.versions {
  display: grid;
  gap: 6px;
}
.ver {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 12px;
  border: 1px solid var(--border-primary);
  border-radius: var(--radius);
  background: var(--bg-primary);
  color: var(--text-secondary);
  font-family: inherit;
  font-size: 13px;
  cursor: pointer;
  transition: all var(--transition);
}
.ver:hover {
  border-color: var(--primary-200);
}
.ver.active {
  background: var(--primary-50);
  border-color: var(--primary-200);
  color: var(--primary-600);
  font-weight: 500;
}
.tag {
  font-size: 11px;
  color: var(--text-muted);
}
.digest {
  display: flex;
  justify-content: space-between;
  gap: 8px;
  margin-top: 12px;
  font-size: 12px;
  color: var(--text-tertiary);
}

.facts {
  list-style: none;
  display: grid;
  gap: 10px;
  font-size: 14px;
  color: var(--text-secondary);
}
.facts li,
.links li {
  display: flex;
  align-items: center;
  gap: 10px;
}
.facts .icon,
.links .icon {
  color: var(--text-muted);
}
.links {
  list-style: none;
  display: grid;
  gap: 10px;
  font-size: 14px;
}

.tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}
.tag-link {
  padding: 4px 10px;
  border: 1px solid var(--border-primary);
  border-radius: 999px;
  background: var(--bg-primary);
  color: var(--text-secondary);
  font-size: 12.5px;
}
.tag-link:hover {
  border-color: var(--primary-200);
  color: var(--primary-600);
}

@media (max-width: 900px) {
  .body {
    grid-template-columns: 1fr;
  }
}
</style>
