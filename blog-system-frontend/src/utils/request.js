import axios from 'axios'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/user'

// 创建axios实例
const request = axios.create({
  baseURL: 'http://127.0.0.1:8080/api', // 后端API的基础地址
  timeout: 10000 // 请求超时时间10秒
})

// 请求拦截器：在发送请求前自动添加Token
request.interceptors.request.use(
  config => {
    const userStore = useUserStore()
    // 如果有Token，就加到请求头里
    if (userStore.token) {
      config.headers.Authorization = `Bearer ${userStore.token}`
    }
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

// 响应拦截器：统一处理后端返回的结果
request.interceptors.response.use(
  response => {
    const res = response.data
    // 如果后端返回的code不是200，说明操作失败，弹出错误提示
    if (res.code !== 200) {
      ElMessage.error(res.msg || '请求失败')
      return Promise.reject(new Error(res.msg || '请求失败'))
    }
    return res // 成功的话，直接返回数据
  },
  error => {
    // 处理网络错误、401未授权等
    if (error.response && error.response.status === 401) {
      ElMessage.error('登录已过期，请重新登录')
      const userStore = useUserStore()
      userStore.logout()
      window.location.href = '/login'
    } else {
      ElMessage.error(error.message || '网络错误')
    }
    return Promise.reject(error)
  }
)

export default request