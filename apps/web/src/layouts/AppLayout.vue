<template>
  <a-layout class="app-shell">
    <AppSidebar class="desktop-sidebar" />
    <a-layout>
      <AppHeader @toggle-menu="drawerOpen = true" />
      <a-layout-content class="app-content">
        <router-view />
      </a-layout-content>
    </a-layout>
    <a-drawer v-model:open="drawerOpen" placement="left" :width="232" :closable="false" class="mobile-drawer">
      <AppSidebar @navigate="drawerOpen = false" />
    </a-drawer>
  </a-layout>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import AppHeader from '@/layouts/AppHeader.vue'
import AppSidebar from '@/layouts/AppSidebar.vue'

const drawerOpen = ref(false)
</script>

<style scoped>
.app-shell {
  min-height: 100vh;
}

.desktop-sidebar {
  display: block;
}

.app-content {
  padding: 18px;
  overflow: auto;
}

@media (max-width: 768px) {
  .desktop-sidebar {
    display: none;
  }

  .app-content {
    padding: 12px;
  }
}
</style>
