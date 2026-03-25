package routes

import (
	"blog-system/controller"
	"blog-system/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupRoutes 配置所有路由
func SetupRoutes(r *gin.Engine) {
	// 1. 使用跨域中间件
	r.Use(middleware.Cors())

	// ⭐ 新增：配置静态文件访问路由，让前端能访问上传的图片
	// 访问 http://127.0.0.1:8080/uploads/xxx.jpg 时，去 ./uploads 文件夹找文件
	r.StaticFS("/uploads", http.Dir("./uploads"))

	// 2. 公开接口（不需要登录）
	public := r.Group("/api")
	{
		// 用户相关
		public.POST("/register", controller.Register)
		public.POST("/login", controller.Login)

		// 博文相关（公开查看）
		public.GET("/posts", controller.GetPostList)
		public.GET("/posts/:id", controller.GetPostDetail)
		// 评论相关（公开查看）
		public.GET("/posts/:id/comments", controller.GetPostComments)
	}

	// 3. 私有接口（需要登录）
	private := r.Group("/api")
	private.Use(middleware.AuthRequired())
	{
		// 用户相关
		private.GET("/user/info", controller.GetUserInfo)

		// 博文相关
		private.POST("/posts", controller.CreatePost)
		private.PUT("/posts/:id", controller.UpdatePost)
		private.DELETE("/posts/:id", controller.DeletePost)

		// 评论相关（发布和删除）
		private.POST("/comments", controller.CreateComment)
		private.DELETE("/comments/:id", controller.DeleteComment)

		// ⭐ 新增：图片上传接口（需要登录才能上传）
		private.POST("/upload/image", controller.UploadImage)
	}
}
