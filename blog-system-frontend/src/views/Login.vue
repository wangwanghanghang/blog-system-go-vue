<template>
  <div class="login-container">
    <el-card class="login-card">
      <template #header>
        <div class="card-header">
          <span>📚 博文管理系统</span>
        </div>
      </template>

      <!-- Tab切换：登录 / 注册 -->
      <el-tabs v-model="activeTab" type="border-card">
        <!-- 登录面板 -->
        <el-tab-pane label="登录" name="login">
          <el-form :model="loginForm" :rules="loginRules" ref="loginFormRef" label-width="60px">
            <el-form-item label="用户名" prop="username">
              <el-input v-model="loginForm.username" placeholder="请输入用户名" />
            </el-form-item>
            <el-form-item label="密码" prop="password">
              <el-input v-model="loginForm.password" type="password" placeholder="请输入密码" show-password />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleLogin" style="width: 100%" :loading="loginLoading">
                登录
              </el-button>
            </el-form-item>
          </el-form>
          <div class="demo-tip">
            <el-text type="info" size="small">测试账号：demo / 密码：123456</el-text>
          </div>
        </el-tab-pane>

        <!-- 注册面板 -->
        <el-tab-pane label="注册" name="register">
          <el-form :model="registerForm" :rules="registerRules" ref="registerFormRef" label-width="60px">
            <el-form-item label="用户名" prop="username">
              <el-input v-model="registerForm.username" placeholder="请设置用户名" />
            </el-form-item>
            <el-form-item label="密码" prop="password">
              <el-input v-model="registerForm.password" type="password" placeholder="请设置密码" show-password />
            </el-form-item>
            <el-form-item label="昵称" prop="nickname">
              <el-input v-model="registerForm.nickname" placeholder="请输入昵称（选填）" />
            </el-form-item>
            <el-form-item>
              <el-button type="success" @click="handleRegister" style="width: 100%" :loading="registerLoading">
                注册
              </el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>
      </el-tabs>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { login, register } from '@/api'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const userStore = useUserStore()

// --------------------------
// 登录相关
// --------------------------
const activeTab = ref('login')
const loginFormRef = ref(null)
const loginLoading = ref(false)
const loginForm = reactive({
  username: '',
  password: ''
})
const loginRules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
}

// 登录按钮点击
const handleLogin = async () => {
  await loginFormRef.value.validate() // 先校验表单
  loginLoading.value = true
  try {
    const res = await login(loginForm)
    userStore.setUserInfo(res.data) // 保存用户信息
    ElMessage.success('登录成功！')
    router.push('/') // 跳转到首页
  } catch (err) {
    console.error(err)
  } finally {
    loginLoading.value = false
  }
}

// --------------------------
// 注册相关
// --------------------------
const registerFormRef = ref(null)
const registerLoading = ref(false)
const registerForm = reactive({
  username: '',
  password: '',
  nickname: ''
})
const registerRules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur', min: 6 }]
}

// 注册按钮点击
const handleRegister = async () => {
  await registerFormRef.value.validate()
  registerLoading.value = true
  try {
    await register(registerForm)
    ElMessage.success('注册成功！请登录')
    activeTab.value = 'login' // 切换到登录面板
    registerForm.username = ''
    registerForm.password = ''
    registerForm.nickname = ''
  } catch (err) {
    console.error(err)
  } finally {
    registerLoading.value = false
  }
}
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.login-card {
  width: 450px;
}

.card-header {
  text-align: center;
  font-size: 20px;
  font-weight: bold;
}

.demo-tip {
  text-align: center;
  margin-top: 10px;
}
</style>