<template>
  <div class="scan-result-card">
    <el-card>
      <template v-if="loading">
        <div class="loading-container">
          <el-skeleton :rows="3" animated />
        </div>
      </template>
      
      <template v-else-if="data && data.length">
        <el-table :data="data" border stripe>
          <el-table-column
            v-for="(column, index) in getColumns"
            :key="index"
            :prop="column.prop"
            :label="column.label"
          />
        </el-table>
      </template>
      
      <template v-else>
        <el-empty description="暂无扫描数据" />
      </template>
    </el-card>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  loading: {
    type: Boolean,
    default: false
  },
  data: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['update:urls'])

// 根据不同类型的数据动态生成表格列
const getColumns = computed(() => {
  if (!props.data.length) return []
  
  return Object.keys(props.data[0]).map(key => ({
    prop: key,
    label: key.charAt(0).toUpperCase() + key.slice(1)
  }))
})

const handleFileChange = (file) => {
  const reader = new FileReader()
  reader.onload = (e) => {
    const content = e.target.result
    const urls = content.split('\n')
      .map(url => url.trim())
      .filter(url => url && url.length > 0)
    emit('update:urls', urls)
  }
  reader.readAsText(file.raw)
}
</script>

<style scoped>
.scan-result-card {
  padding: 10px;
}

.loading-container {
  padding: 20px;
}

.table-actions {
  margin-bottom: 16px;
}

.upload-demo {
  display: inline-block;
}
</style> 