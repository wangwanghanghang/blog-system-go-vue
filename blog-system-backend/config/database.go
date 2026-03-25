package config

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DB 全局数据库连接对象：其他层通过这个变量操作数据库
var DB *gorm.DB

// InitDB 初始化数据库连接：项目启动时调用，连接MySQL
func InitDB() {
	// 数据库连接信息：请把下面的 root:yourpassword 换成你自己的MySQL账号和密码
	// 格式：用户名:密码@tcp(127.0.0.1:3306)/数据库名?charset=utf8mb4&parseTime=True&loc=Local
	dsn := "root:Wh051107@tcp(127.0.0.1:3306)/blog_db?charset=utf8mb4&parseTime=True&loc=Local"

	var err error
	// 用GORM连接MySQL
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("数据库连接失败！请检查账号密码和数据库是否创建：", err) // 如果连接失败，直接终止程序
	}

	fmt.Println("数据库连接成功！")
}
