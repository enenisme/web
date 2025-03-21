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

              <el-form-item class="form-actions">
                <el-button 
                  type="primary" 
                  @click="startScan" 
                  :loading="loading"
                  :disabled="loading"
                  size="large"
                >
                  <el-icon><VideoPlay /></el-icon>
                  {{ loading ? '扫描中...' : '开始扫描' }}
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
                <span>扫描结果 (共 {{ total }} 条)</span>
                <el-button-group>
                  <el-button size="small" @click="exportResults">导出</el-button>
                  <el-button size="small" @click="saveToDatabase">保存</el-button>
                </el-button-group>
              </div>
            </template>
            
            <div v-if="loading" class="loading-container">
              <el-skeleton :rows="5" animated />
            </div>
            
            <template v-else>
              <el-table :data="paginatedResults" style="width: 100%">
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

              <div class="pagination-container">
                <el-pagination
                  v-model:current-page="currentPage"
                  v-model:page-size="pageSize"
                  :page-sizes="[10, 20, 50, 100]"
                  :total="total"
                  layout="total, sizes, prev, pager, next, jumper"
                  @size-change="handleSizeChange"
                  @current-change="handlePageChange"
                  :hide-on-single-page="false"
                />
              </div>
            </template>
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
                  clearable
                >
                  <template #prefix>
                    <el-icon><Search /></el-icon>
                  </template>
                </el-input>
              </div>
            </template>
            
            <div v-if="historyLoading" class="loading-container">
              <el-skeleton :rows="5" animated />
            </div>
            
            <template v-else>
              <el-table 
                v-if="historyData.length" 
                :data="historyData" 
                style="width: 100%"
                border
              >
                <el-table-column prop="date" label="日期" width="180" />
                <el-table-column prop="type" label="扫描类型" width="160" />
                <el-table-column prop="target" label="目标" width="250"/>
                <el-table-column label="结果" >
                  <template #default="{ row }">
                    <el-tooltip 
                      :content="getResultSummary(row)" 
                      placement="top" 
                      :show-after="500"
                    >
                      <div class="result-summary">{{ getResultSummary(row) }}</div>
                    </el-tooltip>
                  </template>
                </el-table-column>
                <el-table-column prop="status" label="状态" width="100">
                  <template #default="{ row }">
                    <el-tag :type="getStatusType(row.status)">{{ row.status }}</el-tag>
                  </template>
                </el-table-column>
                <el-table-column label="操作" width="150" fixed="right">
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

              <el-empty 
                v-else 
                description="暂无历史记录" 
              />

              <div class="pagination-container">
                <el-pagination
                  v-model:current-page="historyCurrentPage"
                  v-model:page-size="historyPageSize"
                  :page-sizes="[10, 20, 50, 100]"
                  :total="historyTotal"
                  layout="total, sizes, prev, pager, next, jumper"
                  @size-change="handleHistorySizeChange"
                  @current-change="handleHistoryPageChange"
                />
              </div>
            </template>
          </el-card>
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessageBox, ElMessage, ElNotification } from 'element-plus'
import request from '../utils/request'
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
import { debounce } from 'lodash'

const router = useRouter()
const activeModule = ref('host')
const scanning = ref(false)
const scanResults = ref([])
const searchQuery = ref('')
const fileList = ref([])

// 添加分页相关的响应式变量
const currentPage = ref(1)
const pageSize = ref(10)

// 添加计算属性获取分页后的数据
const paginatedResults = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return scanResults.value.slice(start, end)
})

// 添加总数计算属性
const total = computed(() => scanResults.value.length)

const form = ref({
  target: '',
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
})

// 添加 loading 状态变量
const loading = ref(false)

// 修改历史数据相关的响应式变量
const historyData = ref([])
const historyTotal = ref(0)
const historyCurrentPage = ref(1)
const historyPageSize = ref(10)
const historyLoading = ref(false)

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

const handleModuleChange = async (index) => {
  activeModule.value = index
  resetForm()
  
  // 当切换到历史记录模块时，自动加载历史数据
  if (index === 'history') {
    await loadHistoryData()
  }
}

const startScan = async () => {
  if (!form.value.target) {
    ElMessage.warning('请输入目标地址')
    return
  }

  loading.value = true
  scanResults.value = []

  try {
    let response
    
    switch (activeModule.value) {
      case 'host':
        response = await request.post('/host/alive', {
          target: form.value.target
        })
        
        if (response.data.data) {
          // 处理主机存活结果
          scanResults.value = response.data.data.map(result => ({
            item: result.IP,
            details: result.Alive ? '主机在线' : '主机离线',
            status: result.Alive ? 'online' : 'offline'
          }))
          
          ElMessage.success(response.data.message || '扫描完成')
        }
        break
        
      case 'port':
        // 构建端口范围
        let portRange
        if (form.value.portScanType === 'range') {
          portRange = `${form.value.portStart}-${form.value.portEnd}`
        } else {
          portRange = form.value.specificPorts
        }
        
        response = await request.post('/port/scan', {
          target: form.value.target,
          portRange: portRange
        })
        
        if (response.data.data) {
          // 处理端口扫描结果
          scanResults.value = response.data.data.ports.map(port => ({
            item: `端口 ${port.port}`,
            details: `服务: ${port.service || '未知'}`,
            status: 'open'
          }))
          
          ElMessage.success(`扫描完成，发现 ${response.data.data.portCount} 个开放端口`)
        }
        break
        
      case 'fingerprint':
        // 保持原有的指纹识别逻辑
        if (form.value.inputType === 'url') {
          processUrls([form.value.target])
        } else {
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
        break
        
      case 'subdomain':
        response = await request.post('/subdomain/scan', {
          domain: form.value.target
        })
        
        if (response.data.data) {
          // 处理子域名扫描结果
          scanResults.value = response.data.data.map(subdomain => ({
            item: subdomain.domain,
            details: subdomain.ip || '未解析',
            status: subdomain.status || 'found'
          }))
          
          ElMessage.success(`扫描完成，发现 ${scanResults.value.length} 个子域名`)
        }
        break
      
      case 'vulnerability':
        response = await request.post('/vulnerability/scan', {
          target: form.value.target
        })

        if (response.data.data) {
          // 处理漏洞扫描结果
          scanResults.value = response.data.data.map(vuln => ({
            item: vuln.VulnName,
            details: vuln.Info,
            status: 'vulnerable'
          }))

          ElMessage.success(`扫描完成，发现 ${scanResults.value.length} 个漏洞`)
        }
        break
        
      // ... 其他 case ...
      // ... 其他模块的处理逻辑 ...
    }
  } catch (error) {
    console.error('扫描出错:', error)
    ElMessage.error(error.response?.data?.error || '扫描失败，请重试')
  } finally {
    loading.value = false
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
    const response = await request.post('/fingerprint/scan', formData, {
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
    loading.value = false
  }
}

const resetForm = () => {
  form.value = {
    target: '',
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

const viewDetail = async (row) => {
  try {
    // 直接使用已有的结果数据，无需再次请求
    const resultData = row.result
    
    // 根据不同的任务类型解析结果
    switch (row.type) {
      case '主机存活':
        // 解析主机存活结果
        try {
          const matches = resultData.match(/\[(.*?)\]/)[1].split(' ')
          scanResults.value = [{
            item: matches[0],
            details: matches[1] === 'true' ? '主机在线' : '主机离线',
            status: matches[1] === 'true' ? 'online' : 'offline'
          }]
        } catch (e) {
          scanResults.value = [{
            item: row.target,
            details: resultData,
            status: 'unknown'
          }]
        }
        break
        
      case '端口扫描':
        // 解析端口扫描结果
        try {
          // 尝试提取端口信息
          const portsMatch = resultData.match(/ports:\[(.*?)\]/)[1]
          const portEntries = portsMatch.match(/map\[(.*?)\]/g)
          
          scanResults.value = portEntries.map(entry => {
            const port = entry.match(/port:(\d+)/)[1]
            const state = entry.match(/state:(\w+)/)[1]
            const service = entry.match(/service:(.*?) /)[1] || '未知'
            
            return {
              item: `端口 ${port}`,
              details: `服务: ${service}`,
              status: state
            }
          })
        } catch (e) {
          scanResults.value = [{
            item: row.target,
            details: resultData,
            status: 'unknown'
          }]
        }
        break
        
      case '子域名扫描':
        // 解析子域名扫描结果
        try {
          const subdomainEntries = resultData.match(/\{(.*?)\}/g)
          
          scanResults.value = subdomainEntries.map(entry => {
            const parts = entry.replace(/[{}]/g, '').split(' ').filter(p => p)
            return {
              item: parts[0] || '未知子域名',
              details: parts.slice(1).join(' ') || '无详细信息',
              status: 'found'
            }
          })
        } catch (e) {
          scanResults.value = [{
            item: row.target,
            details: resultData,
            status: 'unknown'
          }]
        }
        break
        
      default:
        // 默认情况下直接显示原始结果
        scanResults.value = [{
          item: row.target,
          details: resultData,
          status: 'unknown'
        }]
    }
    
    // 切换到对应的模块
    const moduleMap = {
      '主机存活': 'host',
      '端口扫描': 'port',
      '指纹识别': 'fingerprint',
      '子域名扫描': 'subdomain',
      '漏洞检测': 'vulnerability'
    }
    
    activeModule.value = moduleMap[row.type] || 'host'
    ElMessage.success('加载历史记录详情成功')
  } catch (error) {
    console.error('解析历史记录详情失败:', error)
    ElMessage.error('解析历史记录详情失败')
  }
}

const deleteRecord = async (row) => {
  ElMessageBox.confirm('确定要删除这条记录吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await request.delete(`/history/${row.id}`)
      await loadHistoryData() // 重新加载历史数据
      ElMessage.success('删除成功')
    } catch (error) {
      console.error('删除历史记录失败:', error)
      ElMessage.error(error.response?.data?.error || '删除历史记录失败')
    }
  }).catch(() => {})
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

// 添加分页改变处理函数
const handlePageChange = (page) => {
  currentPage.value = page
}

const handleSizeChange = (size) => {
  pageSize.value = size
  currentPage.value = 1
}

// 修改加载历史数据的方法
const loadHistoryData = async () => {
  historyLoading.value = true
  try {
    const response = await request.get('/history', {
      params: {
        page: historyCurrentPage.value,
        pageSize: historyPageSize.value,
        search: searchQuery.value
      }
    })
    
    if (response.data && response.data.data) {
      // 设置总数
      historyTotal.value = response.data.total
      
      // 处理历史数据
      historyData.value = response.data.data.map(item => {
        // 转换任务类型为中文显示
        const typeMap = {
          'host_alive': '主机存活',
          'port_scan': '端口扫描',
          'fingerprint': '指纹识别',
          'subdomain': '子域名扫描',
          'vulnerability': '漏洞检测'
        }
        
        // 转换状态为中文显示
        const statusMap = {
          'completed': '完成',
          'running': '进行中',
          'failed': '失败',
          'pending': '等待中'
        }
        
        // 格式化日期
        const startTime = new Date(item.startTime)
        const formattedDate = startTime.toLocaleString('zh-CN')
        
        return {
          id: item.id,
          date: formattedDate,
          type: typeMap[item.taskType] || item.taskType,
          target: item.target,
          status: statusMap[item.status] || item.status,
          result: item.result
        }
      })
    }
  } catch (error) {
    console.error('加载历史数据失败:', error)
    ElMessage.error(error.response?.data?.error || '加载历史数据失败')
  } finally {
    historyLoading.value = false
  }
}

// 添加历史记录分页处理函数
const handleHistoryPageChange = (page) => {
  historyCurrentPage.value = page
  loadHistoryData()
}

const handleHistorySizeChange = (size) => {
  historyPageSize.value = size
  historyCurrentPage.value = 1
  loadHistoryData()
}

// 添加搜索处理函数
const handleSearch = () => {
  historyCurrentPage.value = 1
  loadHistoryData()
}

// 添加防抖处理
const debouncedSearch = debounce(handleSearch, 300)

// 监听搜索输入变化
watch(searchQuery, () => {
  debouncedSearch()
})

// 修改 getResultSummary 方法，确保正确解析和展示结果
const getResultSummary = (row) => {
  try {
    switch (row.type) {
      case '主机存活':
        try {
          const hostMatches = row.result.match(/\[(.*?)\]/)[1].split(' ')
          return hostMatches[1] === 'true' ? '主机在线' : '主机离线'
        } catch (e) {
          return row.result
        }
        
      case '端口扫描':
        try {
          const portCount = row.result.match(/portCount:(\d+)/)[1]
          return `发现 ${portCount} 个开放端口`
        } catch (e) {
          return row.result
        }
        
      case '子域名扫描':
        try {
          const subdomainCount = (row.result.match(/\{/g) || []).length
          return `发现 ${subdomainCount} 个子域名`
        } catch (e) {
          return row.result
        }
        
      case '指纹识别':
        return row.result.length > 30 ? 
          row.result.substring(0, 30) + '...' : 
          row.result
        
      case '漏洞检测':
        try {
          const vulnCount = (row.result.match(/\{/g) || []).length
          return `发现 ${vulnCount} 个漏洞`
        } catch (e) {
          return row.result
        }
        
      default:
        return row.result.length > 30 ? 
          row.result.substring(0, 30) + '...' : 
          row.result
    }
  } catch (error) {
    console.error('解析结果摘要失败:', error)
    return '无法解析结果'
  }
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

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}

:deep(.el-pagination) {
  padding: 0;
  margin: 0;
}

.result-card {
  margin: 20px;
}

:deep(.el-card__body) {
  padding: 20px;
}

.loading-container {
  padding: 20px;
  background: #fff;
  border-radius: 4px;
}

:deep(.el-skeleton__paragraph) {
  padding: 0;
}

:deep(.el-skeleton__item) {
  margin-top: 12px;
}

/* 按钮加载状态样式 */
:deep(.el-button.is-loading) {
  pointer-events: none;
}

:deep(.el-button.is-loading .el-icon) {
  display: none;
}

.result-summary {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 180px;
}
</style> 