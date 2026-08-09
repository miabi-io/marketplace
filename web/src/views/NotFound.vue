<script setup lang="ts">
import EmptyState from '@/components/EmptyState.vue'
</script>

<template>
  <div class="container">
    <EmptyState
      icon="compass-off"
      title="Page not found"
      hint="That URL is not part of the marketplace."
    >
      <RouterLink
        class="btn btn-primary btn-sm"
        to="/"
      >
        Browse the catalog
      </RouterLink>
    </EmptyState>
  </div>
</template>
