<template>
  <div class="security-scan">
    <el-container>
      <el-aside width="200px">
        <div class="logo">
          <h2>Web家族漏洞扫描工具</h2>
        </div>
        <el-menu
          :default-active="activeModule"
          class="menu"
          @select="handleModuleChange"
        >
          <el-menu-item index="host">
            <el-icon><Monitor /></el-icon>
            <span>主机存活</span>
          </el-menu-item>
          <el-menu-item index="port">
            <el-icon><Connection /></el-icon>
            <span>端口扫描</span>
          </el-menu-item>
          <el-menu-item index="fingerprint">
            <el-icon><Stamp /></el-icon>
            <span>指纹识别</span>
          </el-menu-item>
          <el-menu-item index="subdomain">
            <el-icon><Share /></el-icon>
            <span>子域名扫描</span>
          </el-menu-item>
          <el-menu-item index="vulnerability">
            <el-icon><Warning /></el-icon>
            <span>漏洞检测</span>
          </el-menu-item>
          <el-menu-item index="history">
            <el-icon><DataAnalysis /></el-icon>
            <span>扫描历史</span>
          </el-menu-item>
        </el-menu>
      </el-aside>

      <el-container>
        <el-header>
          <div class="header-content">
            <h3>{{ getCurrentModuleTitle }}</h3>
            <div class="header-right">
              <el-dropdown @command="handleCommand">
                <span class="user-dropdown">
                  <el-icon><User /></el-icon>
                  管理员
                  <el-icon><CaretBottom /></el-icon>
                </span>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item command="profile">个人信息</el-dropdown-item>
                    <el-dropdown-item command="settings">系统设置</el-dropdown-item>
                    <el-dropdown-item divided command="logout">退出登录</el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </div>
          </div>
        </el-header>

        <el-main>
          <el-card v-if="activeModule !== 'history'" class="scan-card">
            <template #header>
              <div class="card-header">
                <span>扫描配置</span>
              </div>
            </template>
            
            <el-form :model="form" label-width="120px" class="scan-form">
              <el-form-item :label="getTargetLabel">
                <template v-if="activeModule === 'fingerprint'">
                  <el-radio-group v-model="form.inputType" class="input-type-group">
                    <el-radio label="url">URL输入</el-radio>
                    <el-radio label="file">文件上传</el-radio>
                  </el-radio-group>

                  <div class="target-input-container">
                    <el-input 
                      v-if="form.inputType === 'url'"
                      v-model="form.target" 
                      placeholder="请输入目标URL"
                    />

                    <el-upload
                      v-if="form.inputType === 'file'"
                      class="upload-demo"
                      action=""
                      :auto-upload="false"
                      :multiple="true"
                      :on-change="handleFileChange"
                      :file-list="fileList"
                    >
                      <el-button type="primary">选择文件</el-button>
                      <template #tip>
                        <div class="el-upload__tip">
                          支持上传多个txt文件，每个文件每行一个URL
                        </div>
                      </template>
                    </el-upload>
                  </div>
                </template>

                <template v-else>
                  <el-input 
                    v-model="form.target" 
                    :placeholder="getTargetPlaceholder"
                  />
                </template>
              </el-form-item>
              
              <!-- 端口扫描特有配置 -->
              <template v-if="activeModule === 'port'">
                <el-form-item label="扫描方式">
                  <el-radio-group v-model="form.portScanType">
                    <el-radio label="range">端口范围</el-radio>
                    <el-radio label="specific">指定端口</el-radio>
                  </el-radio-group>
                </el-form-item>

                <el-form-item v-if="form.portScanType === 'range'" label="端口范围">
                  <el-input-number 
                    v-model="form.portStart" 
                    :min="1" 
                    :max="65535"
                    placeholder="起始端口"
                  />
                  <span class="mx-2">至</span>
                  <el-input-number 
                    v-model="form.portEnd" 
                    :min="1" 
                    :max="65535"
                    placeholder="结束端口"
                  />
                </el-form-item>

                <el-form-item v-if="form.portScanType === 'specific'" label="指定端口">
                  <el-input
                    v-model="form.specificPorts"
                    placeholder="请输入端口号，多个端口用逗号分隔"
                  >
                    <template #append>
                      <el-select 
                        v-model="form.commonPort" 
                        placeholder="常用端口" 
                        @change="handleCommonPortSelect"
                        style="width: 120px"
                      >
                        <el-option label="HTTP (80)" value="80" />
                        <el-option label="HTTPS (443)" value="443" />
                        <el-option label="FTP (21)" value="21" />
                        <el-option label="SSH (22)" value="22" />
                        <el-option label="Telnet (23)" value="23" />
                        <el-option label="MySQL (3306)" value="3306" />
                        <el-option label="RDP (3389)" value="3389" />
                      </el-select>
                    </template>
                  </el-input>
                </el-form-item>
              </template>

              <!-- 扫描模式 - 只在非指纹识别模块显示 -->
              <el-form-item v-if="activeModule !== 'fingerprint'" label="扫描模式">
                <el-radio-group v-model="form.mode">
                  <el-radio label="fast">快速模式</el-radio>
                  <el-radio label="normal">标准模式</el-radio>
                  <el-radio label="deep">深度模式</el-radio>
                </el-radio-group>
              </el-form-item>

              <el-form-item label="高级选项">
                <el-checkbox-group v-model="form.options">
                  <el-checkbox label="saveResult">保存结果</el-checkbox>
                  <el-checkbox label="notification">扫描完成通知</el-checkbox>
                </el-checkbox-group>
              </el-form-item>

              <el-form-item class="form-actions">
                <el-button 
                  type="primary" 
                  :loading="scanning"
                  @click="startScan"
                  size="large"
                >
                  <el-icon><VideoPlay /></el-icon>
                  {{ scanning ? '扫描中...' : '开始扫描' }}
                </el-button>
                <el-button 
                  @click="resetForm"
                  size="large"
                >
                  <el-icon><RefreshRight /></el-icon>
                  重置
                </el-button>
              </el-form-item>
            </el-form>
          </el-card>

          <!-- 扫描结果展示 -->
          <el-card v-if="scanResults.length > 0" class="result-card">
            <template #header>
              <div class="card-header">
                <span>扫描结果</span>
                <el-button-group>
                  <el-button size="small" @click="exportResults">导出</el-button>
                  <el-button size="small" @click="saveToDatabase">保存</el-button>
                </el-button-group>
              </div>
            </template>
            
            <el-table :data="scanResults" style="width: 100%">
              <el-table-column prop="item" label="项目" />
              <el-table-column prop="details" label="详情" />
              <el-table-column prop="status" label="状态">
                <template #default="{ row }">
                  <el-tag :type="getResultStatusType(row.status)">
                    {{ row.status }}
                  </el-tag>
                </template>
              </el-table-column>
            </el-table>
          </el-card>

          <!-- 历史记录展示 -->
          <el-card v-if="activeModule === 'history'" class="history-card">
            <template #header>
              <div class="card-header">
                <span>扫描历史</span>
                <el-input
                  v-model="searchQuery"
                  placeholder="搜索历史记录"
                  style="width: 200px"
                >
                  <template #prefix>
                    <el-icon><Search /></el-icon>
                  </template>
                </el-input>
              </div>
            </template>
            
            <el-table :data="historyData" style="width: 100%">
              <el-table-column prop="date" label="日期" width="180" />
              <el-table-column prop="type" label="扫描类型" width="120" />
              <el-table-column prop="target" label="目标" />
              <el-table-column prop="status" label="状态" width="100">
                <template #default="{ row }">
                  <el-tag :type="getStatusType(row.status)">{{ row.status }}</el-tag>
                </template>
              </el-table-column>
              <el-table-column label="操作" width="150">
                <template #default="{ row }">
                  <el-button size="small" @click="viewDetail(row)">查看</el-button>
                  <el-button 
                    size="small" 
                    type="danger" 
                    @click="deleteRecord(row)"
                  >删除</el-button>
                </template>
              </el-table-column>
            </el-table>
          </el-card>
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessageBox, ElMessage, ElNotification } from 'element-plus'
import { 
  Monitor, 
  Connection, 
  Stamp, 
  Share, 
  Warning, 
  DataAnalysis,
  User, 
  CaretBottom,
  VideoPlay,
  RefreshRight,
  Search
} from '@element-plus/icons-vue'
import axios from 'axios'

const router = useRouter()
const activeModule = ref('host')
const scanning = ref(false)
const scanResults = ref([])
const searchQuery = ref('')
const fileList = ref([])

const form = ref({
  target: '',
  mode: 'normal',
  options: [],
  portScanType: 'range',
  portStart: 1,
  portEnd: 1000,
  specificPorts: '',
  commonPort: '',
  inputType: 'url',
  multipleTargets: '',
  fileUrls: []
})

// 模拟历史数据
const historyData = ref([
  {
    date: '2024-01-20 10:30:00',
    type: '端口扫描',
    target: '192.168.1.1',
    status: '完成'
  },
  {
    date: '2024-01-19 15:45:00',
    type: '漏洞检测',
    target: 'example.com',
    status: '完成'
  }
])

const getCurrentModuleTitle = computed(() => {
  const titles = {
    host: '主机存活',
    port: '端口扫描',
    fingerprint: '指纹识别',
    subdomain: '子域名扫描',
    vulnerability: '漏洞检测',
    history: '扫描历史'
  }
  return titles[activeModule.value]
})

const getTargetLabel = computed(() => {
  const labels = {
    host: 'IP地址',
    port: '目标主机',
    fingerprint: '目标URL',
    subdomain: '域名',
    vulnerability: '目标URL'
  }
  return labels[activeModule.value] || '目标'
})

const getTargetPlaceholder = computed(() => {
  const placeholders = {
    host: '请输入IP地址，例如：192.168.1.1',
    port: '请输入主机地址，例如：example.com',
    fingerprint: '请输入URL，例如：http://example.com',
    subdomain: '请输入域名，例如：example.com',
    vulnerability: '请输入URL，例如：http://example.com'
  }
  return placeholders[activeModule.value] || '请输入目标'
})

const handleModuleChange = (index) => {
  activeModule.value = index
  resetForm()
}

const startScan = () => {
  scanning.value = true
  
  if (form.value.inputType === 'url') {
    // 处理单个URL输入
    processUrls([form.value.target])
  } else {
    // 处理多文件上传
    const promises = fileList.value.map(file => {
      return new Promise((resolve) => {
        const reader = new FileReader()
        reader.onload = (e) => {
          const urls = e.target.result
            .split('\n')
            .map(url => url.trim())
            .filter(url => url && url.length > 0)
          resolve(urls)
        }
        reader.readAsText(file.raw)
      })
    })

    Promise.all(promises).then(urlArrays => {
      const allUrls = urlArrays.flat()
      processUrls(allUrls)
    })
  }
}

const processUrls = async (urls) => {
  try {
    // 创建 FormData 对象
    const formData = new FormData()
    
    if (form.value.inputType === 'url') {
      // 单个 URL 模式
      formData.append('inputType', 'url')
      formData.append('target', form.value.target)
    } else {
      // 文件上传模式
      formData.append('inputType', 'file')
      fileList.value.forEach(file => {
        formData.append('files', file.raw)
      })
    }

    // 发送请求到后端
    const response = await axios.post('/api/fingerprint/scan', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })

    // 处理响应数据
    if (response.data) {
      scanResults.value = [{
        item: form.value.target,
        details: response.data.data.join(', ') || '未识别',
        status: response.data.data.length > 0 ? 'identified' : 'unknown'
      }]
      
      ElMessage.success(response.data.message || '扫描完成')
    }
  } catch (error) {
    ElMessage.error(error.response?.data?.error || '扫描失败，请重试')
  } finally {
    scanning.value = false
  }
}

const resetForm = () => {
  form.value = {
    target: '',
    // 只在非指纹识别模块时设置 mode
    ...(activeModule.value !== 'fingerprint' && { mode: 'normal' }),
    options: [],
    portScanType: 'range',
    portStart: 1,
    portEnd: 1000,
    specificPorts: '',
    commonPort: '',
    // 如果是指纹识别模块，重置相关字段
    ...(activeModule.value === 'fingerprint' ? {
      inputType: 'url',
      multipleTargets: '',
      fileUrls: []
    } : {})
  }
  scanResults.value = []
}

const getStatusType = (status) => {
  return status === '完成' ? 'success' : 'warning'
}

const getResultStatusType = (status) => {
  const types = {
    online: 'success',
    offline: 'danger',
    open: 'success',
    closed: 'info',
    identified: 'success',
    active: 'success',
    inactive: 'info',
    high: 'danger',
    medium: 'warning',
    low: 'info'
  }
  return types[status] || 'info'
}

const handleCommand = (command) => {
  switch (command) {
    case 'logout':
      ElMessageBox.confirm(
        '确定要退出登录吗？',
        '提示',
        {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning',
        }
      ).then(() => {
        localStorage.removeItem('token')
        router.push('/login')
        ElMessage.success('已成功退出登录')
      }).catch(() => {})
      break
    case 'profile':
      ElMessage.info('功能开发中...')
      break
    case 'settings':
      ElMessage.info('功能开发中...')
      break
  }
}

const exportResults = () => {
  ElMessage.success('导出成功')
}

const saveToDatabase = () => {
  ElMessage.success('保存成功')
}

const viewDetail = (row) => {
  console.log('查看详情', row)
}

const deleteRecord = (row) => {
  ElMessageBox.confirm('确定要删除这条记录吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    ElMessage.success('删除成功')
  })
}

// 添加处理常用端口选择的方法
const handleCommonPortSelect = (value) => {
  if (!value) return
  
  const currentPorts = form.value.specificPorts
    .split(',')
    .map(port => port.trim())
    .filter(port => port !== '')
  
  // 检查端口是否已经存在
  if (!currentPorts.includes(value)) {
    const newPorts = [...currentPorts, value]
    form.value.specificPorts = newPorts.join(', ')
  }
  
  // 重置选择
  form.value.commonPort = ''
}

const handleFileChange = (uploadFile) => {
  fileList.value.push(uploadFile)
}
</script>

<style scoped>
.security-scan {
  height: 100vh;
  width: 100vw;
  position: fixed;
  top: 0;
  left: 0;
  overflow: hidden;
  background: #f5f7fa;
  display: flex;
}

.el-container {
  height: 100vh;
  width: 100vw;
  margin: 0;
  padding: 0;
}

.el-aside {
  height: 100vh;
  background: #304156;
  color: #fff;
  overflow-x: hidden;
  overflow-y: auto;
  padding: 0;
  margin: 0;
  box-shadow: none;
}

.logo {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #2b3a51;
  margin: 0;
  padding: 0;
}

.logo h2 {
  color: #fff;
  margin: 0;
  font-size: 16px;
  white-space: nowrap;
  text-shadow: 1px 1px 2px rgba(0, 0, 0, 0.3);
}

.menu {
  border-right: none;
  background: #304156;
  height: calc(100vh - 60px);
  margin: 0;
  padding: 0;
}

.el-header {
  height: 60px;
  background: #fff;
  border-bottom: 1px solid #e6e6e6;
  padding: 0 20px;
  margin: 0;
}

.el-main {
  height: calc(100vh - 60px);
  padding: 0;
  margin: 0;
  overflow-y: auto;
  background: #f5f7fa;
}

.scan-card {
  margin: 0;
  border-radius: 0;
}

:deep(.el-card) {
  border: none;
  box-shadow: none;
}

:deep(.el-card__header) {
  padding: 15px 20px;
  border-bottom: 1px solid #ebeef5;
}

:deep(.el-card__body) {
  padding: 20px;
}

.header-content {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0;
  margin: 0;
}

.header-right {
  display: flex;
  align-items: center;
  margin: 0;
}

.user-dropdown {
  display: flex;
  align-items: center;
  cursor: pointer;
  padding: 0 12px;
  height: 40px;
  color: #303133;
  transition: all 0.3s;
}

.scan-form {
  padding: 20px;
  margin: 0;
}

:deep(.el-form-item) {
  margin-bottom: 18px;
}

:deep(.el-form-item:last-child) {
  margin-bottom: 0;
}

.form-actions {
  display: flex;
  justify-content: center;
  gap: 20px;
}

:deep(.form-actions .el-button) {
  min-width: 140px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.result-card, .history-card {
  margin-top: 20px;
}

/* 自定义滚动条样式 */
:deep(*::-webkit-scrollbar) {
  width: 6px;
  height: 6px;
}

:deep(*::-webkit-scrollbar-track) {
  background: #f1f1f1;
}

:deep(*::-webkit-scrollbar-thumb) {
  background: #c1c1c1;
  border-radius: 3px;
}

:deep(*::-webkit-scrollbar-thumb:hover) {
  background: #a8a8a8;
}

.mx-2 {
  margin: 0 12px;
}

:deep(.el-input-number) {
  width: 160px;
}

:deep(.el-menu) {
  border-right: none;
  background: #304156;
}

:deep(.el-menu-item) {
  height: 56px;
  line-height: 56px;
  color: #bfcbd9;
  text-shadow: 1px 1px 1px rgba(0, 0, 0, 0.2);
  position: relative;
  transition: all 0.3s ease;
}

:deep(.el-menu-item.is-active) {
  color: #fff;
  text-shadow: 1px 1px 2px rgba(0, 0, 0, 0.3);
  background: linear-gradient(90deg, #1890ff 0%, rgba(24, 144, 255, 0.8) 100%);
}

:deep(.el-menu-item.is-active::after) {
  content: '';
  position: absolute;
  top: 0;
  right: 0;
  width: 4px;
  height: 100%;
  background: #fff;
  box-shadow: -1px 0 4px rgba(0, 0, 0, 0.2);
}

:deep(.el-menu-item:hover) {
  color: #fff;
  text-shadow: 1px 1px 2px rgba(0, 0, 0, 0.3);
  background: linear-gradient(90deg, rgba(24, 144, 255, 0.7) 0%, rgba(24, 144, 255, 0.5) 100%);
}

:deep(.el-menu-item .el-icon) {
  color: inherit;
  filter: drop-shadow(1px 1px 1px rgba(0, 0, 0, 0.2));
  margin-right: 16px;
  font-size: 18px;
  vertical-align: middle;
}

:deep(.el-menu-item span) {
  font-size: 14px;
  letter-spacing: 0.5px;
}

/* 添加选中项的动画效果 */
:deep(.el-menu-item.is-active),
:deep(.el-menu-item:hover) {
  transition: all 0.3s cubic-bezier(0.645, 0.045, 0.355, 1);
}

/* 优化菜单间距 */
:deep(.el-menu-item + .el-menu-item) {
  margin-top: 4px;
}

/* 端口选择相关样式 */
:deep(.el-input-group__append) {
  padding: 0;
  background-color: #fff;
}

:deep(.el-input-group__append .el-select) {
  margin: 0;
}

:deep(.el-select .el-input) {
  width: 120px;
}

/* 自定义通知样式 */
.custom-notification {
  background: linear-gradient(to right, #ffffff, #f8f9fa);
  border-left: 4px solid var(--el-notification-color);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  border-radius: 4px;
  padding: 16px;
  transition: all 0.3s ease;
}

.custom-notification:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.12);
}

.custom-notification .el-notification__title {
  font-weight: 600;
  font-size: 16px;
  margin-bottom: 8px;
}

.custom-notification .el-notification__content {
  font-size: 14px;
  line-height: 1.5;
  color: #606266;
}

/* 自定义警告弹窗样式 */
.custom-alert {
  border-radius: 8px;
  padding: 20px;
}

.custom-alert .el-message-box__header {
  padding-bottom: 15px;
}

.custom-alert .el-message-box__title {
  font-size: 18px;
  color: #e6a23c;
}

.custom-alert .el-message-box__content {
  padding: 20px 0;
  font-size: 16px;
  color: #606266;
}

.custom-alert .el-message-box__btns {
  padding-top: 15px;
}

.custom-alert .el-button {
  padding: 12px 24px;
  font-size: 14px;
}

.input-type-group {
  margin-bottom: 15px;
  display: block;
}

.target-input-container {
  width: 100%;
}

.upload-container {
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  padding: 20px;
  background-color: #fff;
}

.upload-area {
  display: flex;
  align-items: center;
  gap: 15px;
}

.upload-demo {
  width: 100%;
}

.el-upload__tip {
  color: #909399;
  font-size: 12px;
  margin: 0;
  line-height: 1.4;
}

.upload-file-list {
  margin-top: 15px;
}

.url-preview {
  margin-top: 10px;
}

.url-count {
  margin-bottom: 8px;
  color: #606266;
  font-size: 13px;
}

:deep(.el-textarea__inner) {
  font-family: monospace;
}

:deep(.el-upload) {
  width: 100%;
}

:deep(.el-upload-list) {
  display: none;
}
</style> 