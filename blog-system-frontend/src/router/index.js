import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/stores/user'

// 定义路由：每个路由对应一个页面
const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
    meta: { requiresAuth: false } // 不需要登录就能访问
  },
  {
    path: '/',
    name: 'Home',
    component: () => import('@/views/Home.vue'),
    meta: { requiresAuth: false } // 首页公开
  },
  {
    path: '/post/:id',
    name: 'PostDetail',
    component: () => import('@/views/PostDetail.vue'),
    meta: { requiresAuth: false } // 详情页公开
  },
  {
    path: '/write',
    name: 'WritePost',
    component: () => import('@/views/WritePost.vue'),
    meta: { requiresAuth: true } // 写博文需要登录
  },
  {
    path: '/edit/:id',
    name: 'EditPost',
    component: () => import('@/views/WritePost.vue'),
    meta: { requiresAuth: true } // 编辑博文需要登录
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫：在跳转前检查是否需要登录
router.beforeEach((to, from) => {
  const userStore = useUserStore()
  
  if (to.meta.requiresAuth && !userStore.token) {
    // 如果页面需要登录但用户没登录，跳转到登录页
    return '/login'
  } else {
    return true
  }
})

export default router