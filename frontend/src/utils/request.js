import axios from 'axios'

const request = axios.create({
  baseURL: 'http://localhost:8889/api', 
})

// 添加请求拦截器，用于调试
request.interceptors.request.use(
  config => {
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
    console.error('Response Error:', error.response?.data || error)
    return Promise.reject(error)
  }
)

export default request 