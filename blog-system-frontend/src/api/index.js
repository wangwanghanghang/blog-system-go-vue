import request from '@/utils/request'

// --------------------------
// 用户相关接口
// --------------------------
export const register = (data) => request.post('/register', data)
export const login = (data) => request.post('/login', data)
export const getUserInfo = () => request.get('/user/info')

// --------------------------
// 博文相关接口
// --------------------------
export const getPostList = (params) => request.get('/posts', { params })
export const getPostDetail = (id) => request.get(`/posts/${id}`)
export const createPost = (data) => request.post('/posts', data)
export const updatePost = (id, data) => request.put(`/posts/${id}`, data)
export const deletePost = (id) => request.delete(`/posts/${id}`)

// --------------------------
// 评论相关接口
// --------------------------
export const getComments = (postId, params) => request.get(`/posts/${postId}/comments`, { params })
export const createComment = (data) => request.post('/comments', data)
export const deleteComment = (id) => request.delete(`/comments/${id}`)