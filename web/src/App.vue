<script setup lang="ts">
import SiteHeader from '@/components/SiteHeader.vue'
import SiteFooter from '@/components/SiteFooter.vue'
</script>

<template>
  <a
    class="skip-link"
    href="#main"
  >Skip to content</a>
  <SiteHeader />
  <main id="main">
    <RouterView />
  </main>
  <SiteFooter />
</template>

<style scoped>
#app,
main {
  min-height: 60vh;
}

.skip-link {
  position: absolute;
  left: -9999px;
  top: 0;
  z-index: 100;
  padding: 10px 16px;
  background: var(--primary-600);
  color: var(--text-on-primary);
  border-radius: 0 0 var(--radius) 0;
}
.skip-link:focus {
  left: 0;
}
</style>
