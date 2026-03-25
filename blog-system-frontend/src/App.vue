<template>
  <router-view />
</template>

<script setup>
</script>

<style>
/* 全局样式重置 */
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: 'Helvetica Neue', Helvetica, 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', Arial, sans-serif;
  background-color: #f5f7fa;
}

/* Markdown内容样式（后面详情页会用到） */
.markdown-body {
  line-height: 1.8;
  color: #333;
}

.markdown-body h1, .markdown-body h2, .markdown-body h3 {
  margin-top: 24px;
  margin-bottom: 16px;
  font-weight: 600;
  line-height: 1.25;
}

.markdown-body h1 {
  font-size: 2em;
  padding-bottom: 0.3em;
  border-bottom: 1px solid #eaecef;
}

.markdown-body h2 {
  font-size: 1.5em;
  padding-bottom: 0.3em;
  border-bottom: 1px solid #eaecef;
}

.markdown-body p {
  margin-top: 0;
  margin-bottom: 16px;
}

.markdown-body pre {
  padding: 16px;
  overflow: auto;
  font-size: 85%;
  line-height: 1.45;
  background-color: #f6f8fa;
  border-radius: 6px;
  margin-bottom: 16px;
}

.markdown-body code {
  padding: 0.2em 0.4em;
  margin: 0;
  font-size: 85%;
  background-color: rgba(175, 184, 193, 0.2);
  border-radius: 3px;
}

.markdown-body pre code {
  padding: 0;
  margin: 0;
  font-size: 100%;
  background-color: transparent;
}

/* ⭐ 新增：解决图片过大的核心样式 */
.markdown-body img {
  max-width: 100%; /* 图片最大宽度不超过父容器 */
  height: auto; /* 高度按比例自适应，不会变形 */
  border-radius: 4px; /* 给图片加个圆角，更美观 */
  margin: 10px 0; /* 给图片上下加间距，和文字分开 */
  display: block; /* 让图片居中 */
  margin-left: auto;
  margin-right: auto;
}
</style>
