<script setup lang="ts">
import { computed } from 'vue'
import {
  mdiAccountOutline,
  mdiAlertCircleOutline,
  mdiApi,
  mdiArrowLeft,
  mdiBookOpenPageVariantOutline,
  mdiBugOutline,
  mdiCheck,
  mdiChevronLeft,
  mdiChevronRight,
  mdiClose,
  mdiCompassOffOutline,
  mdiContentCopy,
  mdiCubeOutline,
  mdiDatabaseOutline,
  mdiDownload,
  mdiFileDocumentOutline,
  mdiFilterRemoveOutline,
  mdiGithub,
  mdiHarddisk,
  mdiMagnify,
  mdiMagnifyClose,
  mdiOpenInNew,
  mdiPackageVariant,
  mdiPackageVariantRemove,
  mdiShieldCheckOutline,
  mdiSourceBranch,
  mdiTagOutline,
  mdiWeatherNight,
  mdiWeatherSunny,
} from '@mdi/js'

// Inline SVG paths from @mdi/js rather than the @mdi/font webfont: the storefront
// uses ~20 glyphs, and the font ships ~740 KB of woff2/eot plus a 337 KB
// stylesheet. Vite tree-shakes these imports down to the paths actually used.
const PATHS: Record<string, string> = {
  'account': mdiAccountOutline,
  'alert': mdiAlertCircleOutline,
  'api': mdiApi,
  'arrow-left': mdiArrowLeft,
  'book': mdiBookOpenPageVariantOutline,
  'bug': mdiBugOutline,
  'check': mdiCheck,
  'chevron-left': mdiChevronLeft,
  'chevron-right': mdiChevronRight,
  'close': mdiClose,
  'compass-off': mdiCompassOffOutline,
  'copy': mdiContentCopy,
  'cube': mdiCubeOutline,
  'database': mdiDatabaseOutline,
  'download': mdiDownload,
  'external': mdiOpenInNew,
  'file': mdiFileDocumentOutline,
  'filter-remove': mdiFilterRemoveOutline,
  'github': mdiGithub,
  'harddisk': mdiHarddisk,
  'moon': mdiWeatherNight,
  'no-results': mdiMagnifyClose,
  'package': mdiPackageVariant,
  'package-gone': mdiPackageVariantRemove,
  'search': mdiMagnify,
  'shield': mdiShieldCheckOutline,
  'source-branch': mdiSourceBranch,
  'sun': mdiWeatherSunny,
  'tag': mdiTagOutline,
}

const props = withDefaults(defineProps<{ name: string; size?: number }>(), { size: 18 })
const path = computed(() => PATHS[props.name] ?? PATHS.package)
</script>

<template>
  <svg
    class="icon"
    :width="size"
    :height="size"
    viewBox="0 0 24 24"
    fill="currentColor"
    aria-hidden="true"
    focusable="false"
  >
    <path :d="path" />
  </svg>
</template>

<style scoped>
.icon {
  flex-shrink: 0;
  display: inline-block;
  vertical-align: -0.15em;
}
</style>
