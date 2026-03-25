<template>
  <div class="detail-container">
    <!-- 顶部导航栏（和首页一样） -->
    <el-header class="header">
      <div class="header-content">
        <div class="logo" @click="router.push('/')">📚 我的博文</div>
        <div class="header-right">
          <template v-if="userStore.token">
            <el-button type="primary" @click="router.push('/write')">
              <el-icon><Edit /></el-icon> 写博文
            </el-button>
            <span class="user-name">{{ userStore.username }}</span>
          </template>
          <template v-else>
            <el-button @click="router.push('/login')">登录</el-button>
          </template>
        </div>
      </div>
    </el-header>

    <!-- 博文内容区 -->
    <el-main class="main-content">
      <div v-if="loading" class="loading">
        <el-skeleton :rows="10" animated />
      </div>

      <template v-else-if="post">
        <el-card class="detail-card">
          <!-- 标题 -->
          <h1 class="detail-title">{{ post.title }}</h1>
          
          <!-- 元信息 -->
          <div class="detail-meta">
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
            <!-- 如果是作者本人，显示编辑和删除按钮 -->
            <template v-if="userStore.userId && post.author_id == userStore.userId">
              <el-button type="primary" size="small" link @click="goToEdit">
                <el-icon><Edit /></el-icon> 编辑
              </el-button>
              <el-button type="danger" size="small" link @click="handleDelete">
                <el-icon><Delete /></el-icon> 删除
              </el-button>
            </template>
          </div>

          <!-- 标签 -->
          <div class="detail-tags" v-if="post.tags">
            <el-tag v-for="tag in post.tags.split(',')" :key="tag" size="small" type="success" effect="plain" style="margin-right: 8px">
              {{ tag }}
            </el-tag>
          </div>

          <el-divider />

          <!-- Markdown内容渲染 -->
          <div class="markdown-body" v-html="renderedContent"></div>
        </el-card>
      </template>
    </el-main>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Edit, User, View, Clock, Delete } from '@element-plus/icons-vue'
import { marked } from 'marked'
import hljs from 'highlight.js'
import 'highlight.js/styles/github.css'
import { getPostDetail, deletePost } from '@/api'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

// 配置marked，支持代码高亮
marked.setOptions({
  highlight: function(code, lang) {
    const language = hljs.getLanguage(lang) ? lang : 'plaintext'
    return hljs.highlight(code, { language }).value
  },
  breaks: true
})

const loading = ref(true)
const post = ref(null)

// 渲染后的HTML内容
const renderedContent = computed(() => {
  if (!post.value) return ''
  return marked.parse(post.value.content)
})

// 获取博文详情
const fetchPostDetail = async () => {
  loading.value = true
  try {
    const res = await getPostDetail(route.params.id)
    post.value = res.data
  } catch (err) {
    ElMessage.error('获取博文详情失败')
  } finally {
    loading.value = false
  }
}

// 格式化日期
const formatDate = (dateStr) => {
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-CN', { year: 'numeric', month: 'long', day: 'numeric' })
}

// 跳转到编辑页
const goToEdit = () => {
  router.push(`/edit/${post.value.id}`)
}

// 删除博文
const handleDelete = async () => {
  try {
    await ElMessageBox.confirm('确定要删除这篇博文吗？删除后无法恢复！', '警告', {
      confirmButtonText: '确定删除',
      cancelButtonText: '取消',
      type: 'warning',
      confirmButtonClass: 'el-button--danger'
    })
    await deletePost(post.value.id)
    ElMessage.success('删除成功！')
    router.push('/')
  } catch (err) {
    if (err !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

onMounted(() => {
  fetchPostDetail()
})
</script>

<style scoped>
.detail-container {
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
  max-width: 900px;
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

.main-content {
  flex: 1;
  max-width: 900px;
  width: 100%;
  margin: 0 auto;
  padding: 30px 20px;
}

.detail-card {
  padding: 30px;
}

.detail-title {
  font-size: 28px;
  font-weight: bold;
  color: #303133;
  margin-bottom: 20px;
  line-height: 1.4;
}

.detail-meta {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 15px;
  margin-bottom: 15px;
  color: #909399;
  font-size: 14px;
}

.meta-text {
  display: flex;
  align-items: center;
  gap: 4px;
}

.detail-tags {
  margin-bottom: 20px;
}
</style>