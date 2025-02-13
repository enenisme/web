<script setup>
import SecurityScan from './views/SecurityScan.vue'
</script>

<template>
  <router-view></router-view>
</template>

<style>
/* 确保应用容器占满整个视口 */
#app {
  height: 100vh;
  width: 100vw;
  overflow: hidden;
}

/* 移除所有可能的外边距 */
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

/* 自定义滚动条样式 */
::-webkit-scrollbar {
  width: 6px;
  height: 6px;
}

::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 3px;
}

::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 3px;
}

::-webkit-scrollbar-thumb:hover {
  background: #a8a8a8;
}
</style> 