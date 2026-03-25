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

          <el-divider />

          <!-- 评论区 -->
          <div class="comments-section">
            <h3 class="section-title">
              <el-icon><ChatLineRound /></el-icon> 评论 ({{ comments.length }})
            </h3>
            
            <!-- 发表评论 -->
            <div class="comment-form">
              <el-input
                v-model="commentContent"
                type="textarea"
                :rows="3"
                placeholder="写下你的评论..."
                maxlength="500"
                show-word-limit
              />
              <div class="form-actions">
                <el-button type="primary" @click="submitComment" :loading="submittingComment">发表评论</el-button>
              </div>
            </div>

            <!-- 评论列表 -->
            <div class="comment-list">
              <div v-if="comments.length === 0" class="no-comments">暂无评论，快来抢沙发吧~</div>
              <div v-else class="comment-item" v-for="comment in comments" :key="comment.id">
                <el-avatar :size="40" class="comment-avatar">{{ comment.user?.username?.charAt(0).toUpperCase() }}</el-avatar>
                <div class="comment-content">
                  <div class="comment-header">
                    <span class="comment-author">{{ comment.user?.nickname || comment.user?.username }}</span>
                    <span class="comment-time">{{ formatDate(comment.created_at) }}</span>
                  </div>
                  <div class="comment-text">{{ comment.content }}</div>
                  
                  <div class="comment-actions" v-if="userStore.userId === comment.user_id || userStore.isAdmin">
                    <el-button type="danger" link size="small" @click="handleDeleteComment(comment.id)">删除</el-button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </el-card>
      </template>
    </el-main>
  </div>
</template>

<style scoped>
/* 原有样式保留... */

.comments-section {
  margin-top: 40px;
}

.section-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 20px;
  margin-bottom: 20px;
  color: #333;
}

.comment-form {
  margin-bottom: 30px;
}

.form-actions {
  margin-top: 10px;
  text-align: right;
}

.comment-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.no-comments {
  text-align: center;
  color: #999;
  padding: 40px 0;
}

.comment-item {
  display: flex;
  gap: 16px;
  padding-bottom: 20px;
  border-bottom: 1px solid #eee;
}

.comment-content {
  flex: 1;
}

.comment-header {
  display: flex;
  align-items: center;
  margin-bottom: 8px;
}

.comment-author {
  font-weight: 600;
  color: #333;
  margin-right: 12px;
}

.comment-time {
  font-size: 12px;
  color: #999;
}

.comment-text {
  font-size: 14px;
  color: #666;
  line-height: 1.6;
}

.comment-actions {
  margin-top: 8px;
}
</style>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Edit, User, View, Clock, Delete, ChatLineRound } from '@element-plus/icons-vue'
import { marked } from 'marked'
import hljs from 'highlight.js'
import 'highlight.js/styles/github.css'
import { getPostDetail, deletePost, getComments, createComment, deleteComment } from '@/api'
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

// 评论相关状态
const comments = ref([])
const commentContent = ref('')
const submittingComment = ref(false)

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
    // 获取评论
    fetchComments()
  } catch (err) {
    ElMessage.error('获取博文详情失败')
  } finally {
    loading.value = false
  }
}

// 获取评论列表
const fetchComments = async () => {
  try {
    const res = await getComments(route.params.id)
    comments.value = res.data.comments
  } catch (err) {
    console.error(err)
  }
}

// 提交评论
const submitComment = async () => {
  if (!userStore.token) {
    ElMessage.warning('请先登录后再评论')
    router.push('/login')
    return
  }
  
  if (!commentContent.value.trim()) {
    ElMessage.warning('请输入评论内容')
    return
  }

  submittingComment.value = true
  try {
    await createComment({
      post_id: post.value.id,
      content: commentContent.value
    })
    ElMessage.success('评论成功')
    commentContent.value = ''
    fetchComments()
  } catch (err) {
    ElMessage.error('评论失败')
  } finally {
    submittingComment.value = false
  }
}

// 删除评论
const handleDeleteComment = async (id) => {
  try {
    await ElMessageBox.confirm('确定要删除这条评论吗？', '提示', {
      type: 'warning'
    })
    await deleteComment(id)
    ElMessage.success('删除成功')
    fetchComments()
  } catch (err) {
    // cancelled
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