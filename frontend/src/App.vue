<template>
  <div id="app">
    <router-view />
  </div>
</template>

<script setup lang="ts">
import './assets/theme.css'
</script>

<style>
#app {
  min-height: 100vh;
}

/* 滚动条 */
::-webkit-scrollbar {
  width: 6px;
}
::-webkit-scrollbar-track {
  background: transparent;
}
::-webkit-scrollbar-thumb {
  background: #d1d5db;
  border-radius: 3px;
}
::-webkit-scrollbar-thumb:hover {
  background: #9ca3af;
}
</style>
