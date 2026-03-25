import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useUserStore = defineStore('user', () => {
  // 状态：Token、用户ID、用户名
  const token = ref(localStorage.getItem('token') || '')
  const userId = ref(localStorage.getItem('userId') || '')
  const username = ref(localStorage.getItem('username') || '')

  // 动作：登录成功后保存信息
  function setUserInfo(data) {
    token.value = data.token
    userId.value = data.user_id
    username.value = data.username
    // 同时存到localStorage里，刷新页面后信息不会丢失
    localStorage.setItem('token', data.token)
    localStorage.setItem('userId', data.user_id)
    localStorage.setItem('username', data.username)
  }

  // 动作：退出登录，清空信息
  function logout() {
    token.value = ''
    userId.value = ''
    username.value = ''
    localStorage.removeItem('token')
    localStorage.removeItem('userId')
    localStorage.removeItem('username')
  }

  return { token, userId, username, setUserInfo, logout }
})