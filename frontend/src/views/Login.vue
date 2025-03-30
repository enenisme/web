<template>
  <div class="login-container">
    <div class="login-box">
      <div class="login-header">
        <h2>Web家族漏洞扫描工具</h2>
        <p class="subtitle">Web Family Vulnerability Scanner</p>
      </div>
      
      <el-form
        ref="loginFormRef"
        :model="loginForm"
        :rules="loginRules"
        class="login-form"
      >
        <el-form-item prop="username">
          <el-input
            v-model="loginForm.username"
            placeholder="用户名"
          >
            <template #prefix>
              <el-icon><User /></el-icon>
            </template>
          </el-input>
        </el-form-item>
        
        <el-form-item prop="password">
          <el-input
            v-model="loginForm.password"
            type="password"
            placeholder="密码"
          >
            <template #prefix>
              <el-icon><Lock /></el-icon>
            </template>
          </el-input>
        </el-form-item>

        <el-form-item prop="captcha">
          <div class="captcha-container">
            <el-input
              v-model="loginForm.captcha"
              placeholder="验证码"
              @keyup.enter="handleLogin"
            >
              <template #prefix>
                <el-icon><Key /></el-icon>
              </template>
            </el-input>
            <div class="captcha-image" @click="refreshCaptcha">
              <canvas ref="captchaCanvas" width="120" height="40"></canvas>
            </div>
          </div>
        </el-form-item>
        
        <el-form-item class="remember-me">
          <el-checkbox v-model="loginForm.remember">记住我</el-checkbox>
          <el-link type="primary" :underline="false">忘记密码？</el-link>
        </el-form-item>
        
        <el-form-item>
          <el-button
            type="primary"
            :loading="loading"
            class="login-button"
            @click="handleLogin"
          >
            {{ loading ? '登录中...' : '登录' }}
          </el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { User, Lock, Key } from '@element-plus/icons-vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElNotification } from 'element-plus'

const router = useRouter()
const loginFormRef = ref(null)
const loading = ref(false)

const captchaCanvas = ref(null)
const captchaText = ref('')

const loginForm = reactive({
  username: '',
  password: '',
  remember: false,
  captcha: ''
})

const loginRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '长度在 3 到 20 个字符', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 20, message: '长度在 6 到 20 个字符', trigger: 'blur' }
  ],
  captcha: [
    { required: true, message: '请输入验证码', trigger: 'blur' },
    { min: 4, max: 4, message: '验证码长度为4位', trigger: 'blur' }
  ]
}

const generateCaptcha = () => {
  const canvas = captchaCanvas.value
  const ctx = canvas.getContext('2d')
  
  ctx.fillStyle = '#f5f5f5'
  ctx.fillRect(0, 0, canvas.width, canvas.height)
  
  const digits = '0123456789'
  let code = ''
  for (let i = 0; i < 4; i++) {
    code += digits[Math.floor(Math.random() * digits.length)]
  }
  captchaText.value = code
  
  ctx.fillStyle = '#333'
  ctx.font = '24px Arial'
  ctx.textBaseline = 'middle'
  
  for (let i = 0; i < code.length; i++) {
    const x = 20 + i * 25 + Math.random() * 5
    const y = 20 + Math.random() * 5
    const rotate = (Math.random() - 0.5) * 0.3
    
    ctx.save()
    ctx.translate(x, y)
    ctx.rotate(rotate)
    ctx.fillText(code[i], 0, 0)
    ctx.restore()
  }
  
  for (let i = 0; i < 3; i++) {
    ctx.beginPath()
    ctx.strokeStyle = `rgba(${Math.random() * 255},${Math.random() * 255},${Math.random() * 255},0.5)`
    ctx.moveTo(Math.random() * canvas.width, Math.random() * canvas.height)
    ctx.lineTo(Math.random() * canvas.width, Math.random() * canvas.height)
    ctx.stroke()
  }
  
  for (let i = 0; i < 50; i++) {
    ctx.fillStyle = `rgba(${Math.random() * 255},${Math.random() * 255},${Math.random() * 255},0.5)`
    ctx.beginPath()
    ctx.arc(Math.random() * canvas.width, Math.random() * canvas.height, 1, 0, 2 * Math.PI)
    ctx.fill()
  }
}

const refreshCaptcha = () => {
  generateCaptcha()
  loginForm.captcha = ''
}

onMounted(() => {
  generateCaptcha()
})

// 添加 WAF 检测函数
const checkSQLInjection = (str) => {
  const sqlPatterns = [
    'and\\s+', 'or\\s+', 'from\\s+', 'execute\\s+', 'update\\s+', 
    'count\\s*\\(', 'chr\\s*\\(', 'mid\\s*\\(', 'char\\s*\\(', 
    'union\\s+', 'select\\s+', 'delete\\s+', 'insert\\s+', 
    'limit\\s+', 'concat\\s*\\(', '\\\\', '&&', '\\|\\|',
    "'", '%', '_'
  ]
  
  const pattern = new RegExp(sqlPatterns.join('|'), 'i')
  return pattern.test(str)
}

const checkXSSAttack = (str) => {
  const xssPatterns = [
    'script', 'alert', 'onclick', 'function', 'get', 'post',
    'document', 'cookie', '<', '>', 'onerror', 'onload',
    'javascript:', 'vbscript:', 'data:', 'eval\\s*\\('
  ]
  
  const pattern = new RegExp(xssPatterns.join('|'), 'i')
  return pattern.test(str)
}

const sanitizeInput = (str) => {
  // SQL注入防护
  let sanitized = str
    .replace(/and\s+/gi, 'sqlwaf')
    .replace(/or\s+/gi, 'sqlwaf')
    .replace(/from\s+/gi, 'sqlwaf')
    .replace(/execute\s+/gi, 'sqlwaf')
    .replace(/update\s+/gi, 'sqlwaf')
    .replace(/count/gi, 'sqlwaf')
    .replace(/chr/gi, 'sqlwaf')
    .replace(/mid/gi, 'sqlwaf')
    .replace(/char/gi, 'sqlwaf')
    .replace(/union\s+/gi, 'sqlwaf')
    .replace(/select\s+/gi, 'sqlwaf')
    .replace(/delete\s+/gi, 'sqlwaf')
    .replace(/insert\s+/gi, 'sqlwaf')
    .replace(/limit\s+/gi, 'sqlwaf')
    .replace(/concat/gi, 'sqlwaf')
    .replace(/\\/g, '\\\\')
    .replace(/&&/g, '')
    .replace(/\|\|/g, '')
    .replace(/'/g, '')
    .replace(/%/g, '\\%')
    .replace(/_/g, '\\_')

  // XSS防护
  sanitized = sanitized
    .replace(/script/gi, 'jswaf')
    .replace(/alert/gi, 'jswaf')
    .replace(/onclick/gi, 'jswaf')
    .replace(/function/gi, 'jswaf')
    .replace(/get/gi, 'jswaf')
    .replace(/post/gi, 'jswaf')
    .replace(/document/gi, 'jswaf')
    .replace(/cookie/gi, 'jswaf')
    .replace(/</g, 'jswaf')
    .replace(/>/g, 'jswaf')

  return sanitized
}

const handleLogin = () => {
  if (!loginFormRef.value) return
  
  loginFormRef.value.validate((valid) => {
    if (valid) {
      // 检查SQL注入
      if (checkSQLInjection(loginForm.username) || checkSQLInjection(loginForm.password)) {
        ElNotification({
          title: '安全警告',
          message: '检测到疑似SQL注入攻击，已进行防护处理',
          type: 'warning',
          duration: 3000
        })
        loginForm.username = sanitizeInput(loginForm.username)
        loginForm.password = sanitizeInput(loginForm.password)
        return
      }

      // 检查XSS攻击
      if (checkXSSAttack(loginForm.username) || checkXSSAttack(loginForm.password)) {
        ElNotification({
          title: '安全警告',
          message: '检测到疑似XSS攻击，已进行防护处理',
          type: 'warning',
          duration: 3000
        })
        loginForm.username = sanitizeInput(loginForm.username)
        loginForm.password = sanitizeInput(loginForm.password)
        return
      }

      loading.value = true
      
      if (loginForm.captcha.toLowerCase() !== captchaText.value.toLowerCase()) {
        ElMessage.error('验证码错误')
        refreshCaptcha()
        loading.value = false
        return
      }
      
      if (loginForm.username === 'admin' && loginForm.password === '123456') {
        localStorage.setItem('token', 'demo-token')
        router.push('/')
        ElMessage.success('登录成功')
      } else {
        ElMessage.error('用户名或密码错误')
        refreshCaptcha()
      }
      loading.value = false
    }
  })
}
</script>

<style scoped>
.login-container {
  height: 100vh;
  width: 100vw;
  position: fixed;
  top: 0;
  left: 0;
  display: flex;
  justify-content: center;
  align-items: center;
  background: linear-gradient(135deg, #1e2a4a 0%, #293952 100%);
  overflow: hidden;
}

.login-box {
  width: 400px;
  padding: 40px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
  animation: fadeIn 0.5s ease;
}

.login-header {
  text-align: center;
  margin-bottom: 40px;
}

.login-header h2 {
  color: #303133;
  font-size: 28px;
  margin: 0;
  margin-bottom: 8px;
}

.subtitle {
  color: #909399;
  font-size: 16px;
  margin: 0;
}

.login-form {
  margin-bottom: 20px;
}

.remember-me {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.login-button {
  width: 100%;
  padding: 12px 0;
}

:deep(.el-input__wrapper) {
  padding: 12px;
}

:deep(.el-input__prefix) {
  margin-right: 8px;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(-20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* 悬浮效果 */
.login-box:hover {
  transform: translateY(-5px);
  transition: all 0.3s ease;
}

:deep(.el-button:hover) {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(64, 158, 255, 0.2);
  transition: all 0.3s ease;
}

.captcha-container {
  display: flex;
  gap: 10px;
  align-items: center;
}

.captcha-image {
  cursor: pointer;
  border-radius: 4px;
  overflow: hidden;
  flex-shrink: 0;
}

.captcha-image canvas {
  display: block;
}

:deep(.el-input) {
  flex: 1;
}
</style> 