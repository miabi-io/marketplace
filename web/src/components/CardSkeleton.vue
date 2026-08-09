<template>
  <!-- Mirrors TemplateCard's box model so the grid does not reflow when real
       results replace it. -->
  <div
    class="card sk"
    aria-hidden="true"
  >
    <div class="row">
      <span class="bone tile" />
      <span class="bone pill" />
    </div>
    <span class="bone line w60" />
    <span class="bone line" />
    <span class="bone line w80" />
    <div class="row foot">
      <span class="bone line w40" />
      <span class="bone line w20" />
    </div>
  </div>
</template>

<style scoped>
.sk {
  display: flex;
  flex-direction: column;
  gap: 10px;
  padding: 18px;
}
.row {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.foot {
  margin-top: 4px;
  padding-top: 12px;
  border-top: 1px solid var(--border-secondary);
}
.bone {
  display: block;
  background: linear-gradient(
    90deg,
    var(--bg-tertiary) 25%,
    var(--bg-hover) 37%,
    var(--bg-tertiary) 63%
  );
  background-size: 400% 100%;
  animation: shimmer 1.4s ease infinite;
  border-radius: var(--radius-sm);
}
.tile {
  width: 44px;
  height: 44px;
  border-radius: var(--radius-md);
}
.pill {
  width: 62px;
  height: 18px;
  border-radius: 999px;
}
.line {
  height: 12px;
  width: 100%;
}
.w80 {
  width: 80%;
}
.w60 {
  width: 60%;
  height: 16px;
}
.w40 {
  width: 40%;
}
.w20 {
  width: 20%;
}

@keyframes shimmer {
  0% {
    background-position: 100% 50%;
  }
  100% {
    background-position: 0 50%;
  }
}
</style>
