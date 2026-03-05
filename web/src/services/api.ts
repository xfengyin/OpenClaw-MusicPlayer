import axios from 'axios'

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || '/api',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器
api.interceptors.request.use(
  config => {
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

// 响应拦截器
api.interceptors.response.use(
  response => {
    if (response.data.code !== 200) {
      return Promise.reject(new Error(response.data.msg || '请求失败'))
    }
    return response.data.data
  },
  error => {
    return Promise.reject(error)
  }
)

export default api
