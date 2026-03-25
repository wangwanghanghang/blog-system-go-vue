<template>
  <div class="write-container">
    <!-- 顶部导航栏 -->
    <el-header class="header">
      <div class="header-content">
        <div class="back-btn" @click="router.back()">
          <el-icon><ArrowLeft /></el-icon> 返回
        </div>
        <div class="header-title">{{ isEdit ? '编辑博文' : '写新博文' }}</div>
        <el-button 
          type="primary" 
          @click="handleSubmit" 
          :loading="submitLoading"
          :disabled="loading"
        >
          <el-icon><Check /></el-icon> {{ isEdit ? '保存修改' : '发布博文' }}
        </el-button>
      </div>
    </el-header>

    <!-- 编辑区 -->
    <el-main class="main-content">
      <el-card class="write-card" v-loading="loading" element-loading-text="正在加载内容，请稍候...">
        <el-form :model="postForm" label-width="80px">
          <!-- 标题 -->
          <el-form-item label="博文标题">
            <el-input 
              v-model="postForm.title" 
              placeholder="请输入一个吸引人的标题" 
              size="large"
              :disabled="loading"
            />
          </el-form-item>

          <!-- 分类和标签 -->
          <el-row :gutter="20">
            <el-col :span="12">
              <el-form-item label="分类">
                <el-input 
                  v-model="postForm.category" 
                  placeholder="如：技术、生活、读书"
                  :disabled="loading"
                />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="标签">
                <el-input 
                  v-model="postForm.tags" 
                  placeholder="多个标签用逗号分隔，如：Go,Vue"
                  :disabled="loading"
                />
              </el-form-item>
            </el-col>
          </el-row>

          <!-- ByteMD 编辑器 -->
          <el-form-item label="博文内容">
            <div class="editor-wrapper">
              <!-- 上传图片的按钮 -->
              <div style="margin-bottom: 10px;">
                <el-upload
                  class="image-uploader"
                  :show-file-list="false"
                  :before-upload="beforeImgUpload"
                  :http-request="handleImgUpload"
                  accept="image/*"
                  :disabled="loading"
                >
                  <el-button size="small" type="primary" :loading="uploadLoading" :disabled="loading">
                    <el-icon><Upload /></el-icon> 插入本地图片
                  </el-button>
                  <span style="margin-left: 10px; color: #909399; font-size: 13px;">支持 jpg、png、gif，最大 5MB</span>
                </el-upload>
              </div>
              
              <!-- ByteMD 编辑器主体（已修复绑定） -->
              <Editor
                :value="postForm.content"
                @change="handleEditorChange"
                :plugins="plugins"
                placeholder="在这里用 Markdown 写你的博文..."
                :disabled="loading"
              />
            </div>
          </el-form-item>
        </el-form>
      </el-card>
    </el-main>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ArrowLeft, Check, Upload } from '@element-plus/icons-vue'
import { Editor } from '@bytemd/vue-next'
import gfm from '@bytemd/plugin-gfm'
import highlight from '@bytemd/plugin-highlight'
import { getPostDetail, createPost, updatePost } from '@/api'
import request from '@/utils/request'

const router = useRouter()
const route = useRoute()

// ByteMD 插件配置（支持表格、代码高亮等）
const plugins = [gfm(), highlight()]

// 判断是编辑还是新建
const isEdit = computed(() => !!route.params.id)
const loading = ref(false) // 全局加载状态
const submitLoading = ref(false)
const uploadLoading = ref(false)

// 表单数据
const postForm = reactive({
  title: '',
  content: '',
  category: '',
  tags: ''
})

// 重置表单
const resetForm = () => {
  postForm.title = ''
  postForm.content = ''
  postForm.category = ''
  postForm.tags = ''
}

// 图片上传前的校验
const beforeImgUpload = (file) => {
  const isImage = file.type.startsWith('image/')
  const isLt5M = file.size / 1024 / 1024 < 5

  if (!isImage) {
    ElMessage.error('只能上传图片文件!')
    return false
  }
  if (!isLt5M) {
    ElMessage.error('图片大小不能超过 5MB!')
    return false
  }
  return true
}

// 处理图片上传
const handleImgUpload = async (options) => {
  uploadLoading.value = true
  const file = options.file
  const formData = new FormData()
  formData.append('image', file)

  try {
    const res = await request.post('/upload/image', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
    
    // 上传成功后，把 Markdown 图片语法插入到内容末尾
    const imgMarkdown = `\n![图片](${res.data.url})\n`
    postForm.content += imgMarkdown
    
    ElMessage.success('图片插入成功！')
  } catch (err) {
    ElMessage.error('图片上传失败')
  } finally {
    uploadLoading.value = false
  }
}

// ⭐ 新增：强制同步编辑器内容到表单
const handleEditorChange = (val) => {
  postForm.content = val
}

// 如果是编辑，先获取原有数据
const fetchPostForEdit = async () => {
  if (!isEdit.value) return
  loading.value = true
  try {
    const res = await getPostDetail(route.params.id)
    const post = res.data
    // 强制赋值，确保所有字段都更新
    postForm.title = post.title
    postForm.content = post.content
    postForm.category = post.category || ''
    postForm.tags = post.tags || ''
  } catch (err) {
    ElMessage.error('获取博文信息失败')
    router.back()
  } finally {
    loading.value = false
  }
}

// 提交发布/保存
const handleSubmit = async () => {
  if (!postForm.title.trim()) {
    ElMessage.warning('请输入博文标题')
    return
  }
  if (!postForm.content.trim()) {
    ElMessage.warning('请输入博文内容')
    return
  }

  // 编辑模式下，二次确认，避免误覆盖
  if (isEdit.value) {
    try {
      await ElMessageBox.confirm('确定要保存修改吗？修改后会覆盖原有内容', '确认保存', {
        confirmButtonText: '确定保存',
        cancelButtonText: '取消',
        type: 'warning'
      })
    } catch {
      return // 用户取消，不执行保存
    }
  }

  submitLoading.value = true
  try {
    if (isEdit.value) {
      await updatePost(route.params.id, postForm)
      ElMessage.success('保存成功！')
      router.back()
    } else {
      const res = await createPost(postForm)
      ElMessage.success('发布成功！')
      resetForm() // 发布成功后重置表单
      const postId = res.data.id || res.data.ID
      router.push(`/post/${postId}`)
    }
  } catch (err) {
    console.error(err)
  } finally {
    submitLoading.value = false
  }
}

// 页面加载时执行
onMounted(() => {
  if (isEdit.value) {
    fetchPostForEdit()
  } else {
    resetForm()
  }
})
</script>

<style scoped>
.write-container {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  background-color: #f5f7fa;
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

.back-btn {
  display: flex;
  align-items: center;
  gap: 5px;
  cursor: pointer;
  color: #606266;
  font-size: 15px;
}

.back-btn:hover {
  color: #409eff;
}

.header-title {
  font-size: 18px;
  font-weight: bold;
}

.main-content {
  flex: 1;
  max-width: 1000px;
  width: 100%;
  margin: 0 auto;
  padding: 20px;
}

.editor-wrapper {
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  padding: 15px;
  background-color: #fff;
}

/* 让 ByteMD 编辑器占满宽度 */
:deep(.bytemd) {
  height: 500px;
}
</style>