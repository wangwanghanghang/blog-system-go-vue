package main

import (
	"blog-system/config"
	"blog-system/models"
	"blog-system/routes"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. 初始化数据库连接
	config.InitDB()

	// 2. 自动创建数据库表
	models.AutoMigrate()

	// 3. ⭐ 新增：自动初始化测试数据（如果数据库为空）
	models.InitTestData()

	// 4. 创建Gin引擎
	r := gin.Default()

	// 5. 配置路由
	routes.SetupRoutes(r)

	// 6. 启动服务，监听8080端口
	fmt.Println("服务启动成功！访问地址：http://127.0.0.1:8080")
	log.Fatal(r.Run(":8080"))
}
