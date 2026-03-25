<template>
  <div class="home-container">
    <!-- 顶部导航栏 -->
    <el-header class="header">
      <div class="header-content">
        <div class="logo" @click="router.push('/')">📚 我的博文</div>
        <div class="header-right">
          <!-- 如果已登录，显示写博文按钮和用户信息 -->
          <template v-if="userStore.token">
            <el-button type="primary" @click="router.push('/write')">
              <el-icon><Edit /></el-icon> 写博文
            </el-button>
            <el-dropdown @command="handleCommand">
              <span class="user-name">
                <el-avatar :size="30" style="margin-right: 8px">{{ userStore.username.charAt(0) }}</el-avatar>
                {{ userStore.username }}
              </span>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="logout">退出登录</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>
          <!-- 如果未登录，显示登录按钮 -->
          <template v-else>
            <el-button @click="router.push('/login')">登录</el-button>
          </template>
        </div>
      </div>
    </el-header>

    <!-- 主内容区 -->
    <el-main class="main-content">
      <div class="post-list">
        <!-- 加载中 -->
        <div v-if="loading" class="loading">
          <el-skeleton :rows="5" animated />
        </div>

        <!-- 博文列表 -->
        <template v-else>
          <el-card v-for="post in postList" :key="post.id" class="post-card" shadow="hover" @click="goToDetail(post.id)">
            <div class="post-title">{{ post.title }}</div>
            <div class="post-meta">
              <el-tag size="small" type="info" v-if="post.category">{{ post.category }}</el-tag>
              <span class="meta-text">
                <el-icon><User /></el-icon> {{ post.author?.nickname || post.author?.username }}
              </span>
              <span class="meta-text">
                <el-icon><View /></el-icon> {{ post.views }} 阅读
              </span>
              <span class="meta-text">
                <el-icon><Clock /></el-icon> {{ formatDate(post.created_at) }}
              </span>
            </div>
            <div class="post-summary">{{ post.content.substring(0, 150).replace(/[#*`]/g, '') }}...</div>
            <div class="post-tags" v-if="post.tags">
              <el-tag v-for="tag in post.tags.split(',')" :key="tag" size="small" type="success" effect="plain" style="margin-right: 5px">
                {{ tag }}
              </el-tag>
            </div>
          </el-card>

          <!-- 分页 -->
          <div class="pagination">
            <el-pagination
              v-model:current-page="currentPage"
              v-model:page-size="pageSize"
              :total="total"
              :page-sizes="[10, 20, 50]"
              layout="total, sizes, prev, pager, next, jumper"
              @size-change="fetchPostList"
              @current-change="fetchPostList"
            />
          </div>
        </template>
      </div>
    </el-main>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Edit, User, View, Clock } from '@element-plus/icons-vue'
import { getPostList } from '@/api'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const userStore = useUserStore()

// 状态
const loading = ref(true)
const postList = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

// 获取博文列表
const fetchPostList = async () => {
  loading.value = true
  try {
    const res = await getPostList({ page: currentPage.value, page_size: pageSize.value })
    postList.value = res.data.list
    total.value = res.data.total
  } catch (err) {
    ElMessage.error('获取博文列表失败')
  } finally {
    loading.value = false
  }
}

// 跳转到详情页
const goToDetail = (id) => {
  router.push(`/post/${id}`)
}

// 格式化日期
const formatDate = (dateStr) => {
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-CN', { year: 'numeric', month: 'long', day: 'numeric' })
}

// 下拉菜单操作
const handleCommand = async (command) => {
  if (command === 'logout') {
    try {
      await ElMessageBox.confirm('确定要退出登录吗？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      })
      userStore.logout()
      ElMessage.success('已退出登录')
      router.push('/')
    } catch {
      // 用户取消
    }
  }
}

// 页面加载时获取数据
onMounted(() => {
  fetchPostList()
})
</script>

<style scoped>
.home-container {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.header {
  background-color: #fff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  padding: 0;
  position: sticky;
  top: 0;
  z-index: 100;
}

.header-content {
  max-width: 1000px;
  margin: 0 auto;
  height: 60px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 20px;
}

.logo {
  font-size: 20px;
  font-weight: bold;
  cursor: pointer;
  color: #409eff;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 15px;
}

.user-name {
  display: flex;
  align-items: center;
  cursor: pointer;
}

.main-content {
  flex: 1;
  max-width: 1000px;
  width: 100%;
  margin: 0 auto;
  padding: 30px 20px;
}

.post-card {
  margin-bottom: 20px;
  cursor: pointer;
  transition: transform 0.2s;
}

.post-card:hover {
  transform: translateY(-2px);
}

.post-title {
  font-size: 20px;
  font-weight: bold;
  color: #303133;
  margin-bottom: 12px;
}

.post-meta {
  display: flex;
  align-items: center;
  gap: 15px;
  margin-bottom: 12px;
  color: #909399;
  font-size: 13px;
}

.meta-text {
  display: flex;
  align-items: center;
  gap: 4px;
}

.post-summary {
  color: #606266;
  line-height: 1.6;
  margin-bottom: 12px;
}

.pagination {
  display: flex;
  justify-content: center;
  margin-top: 30px;
}
</style>