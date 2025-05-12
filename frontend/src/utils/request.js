import axios from 'axios'

const request = axios.create({
  baseURL: 'http://localhost:18889/api', 
})

// 添加请求拦截器，设置认证令牌
request.interceptors.request.use(
  config => {
    // 从localStorage获取令牌
    const token = localStorage.getItem('token')
    
    // 如果令牌存在，添加到请求头
    if (token) {
      config.headers['Authorization'] = `Bearer ${token}`
    }
    
    console.log('Request:', config)
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

// 响应拦截器
request.interceptors.response.use(
  response => response,
  error => {
    // 处理认证错误
    if (error.response && error.response.status === 401) {
      // 清除本地存储的令牌
      localStorage.removeItem('token')
      
      // 如果不是登录页面，则重定向到登录页
      if (window.location.pathname !== '/login') {
        window.location.href = '/login'
      }
    }
    
    console.error('Response Error:', error.response?.data || error)
    return Promise.reject(error)
  }
)

export default request 