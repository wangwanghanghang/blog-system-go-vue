package models

import (
	"blog-system/config"
	"blog-system/utils"
	"log"
)

// AutoMigrate 自动迁移数据库表：项目启动时调用，自动创建/更新表结构
func AutoMigrate() {
	// 自动创建 User 和 Post 表
	err := config.DB.AutoMigrate(&User{}, &Post{})
	if err != nil {
		log.Fatal("数据库表迁移失败：", err)
	}
	log.Println("数据库表创建/更新成功！")
}

// InitTestData 初始化测试数据：如果数据库是空的，就自动插入测试用户和测试博文
func InitTestData() {
	// --------------------------
	// 1. 先创建一个测试用户
	// --------------------------
	var testUser User
	// 检查用户名为 "demo" 的用户是否存在，不存在才创建
	config.DB.Where("username = ?", "demo").First(&testUser)
	if testUser.ID == 0 {
		// 加密测试密码（密码是 123456）
		hashedPassword, _ := utils.HashPassword("123456")

		testUser = User{
			Username: "demo",
			Password: hashedPassword,
			Nickname: "博文测试员",
		}
		// 插入用户
		if err := config.DB.Create(&testUser).Error; err != nil {
			log.Println("创建测试用户失败：", err)
		} else {
			log.Println("✅ 测试用户创建成功！账号：demo，密码：123456")
		}
	}

	// --------------------------
	// 2. 再创建几篇测试博文
	// --------------------------
	var postCount int64
	config.DB.Model(&Post{}).Count(&postCount)
	if postCount < 5 { // 如果博文表是空的，才插入
		testPosts := []Post{
			{
				Title:    "欢迎使用博文管理系统",
				Content:  "# 欢迎来到我的小站 🎉\n\n这是一个用 **Go + Vue3** 开发的个人博文管理系统。\n\n## 主要功能\n\n- ✅ 用户注册登录\n- ✅ 博文的增删改查\n- ✅ Markdown 内容渲染\n- ✅ 分类和标签管理\n\n希望你喜欢这个系统！",
				AuthorID: testUser.ID,
				Category: "公告",
				Tags:     "系统介绍,欢迎",
				Views:    128,
			},
			{
				Title:    "Go语言入门：第一个程序 Hello World",
				Content:  "# Go 语言入门：Hello World\n\n作为一名 Gopher，我们的第一个程序当然是经典的 **Hello World**。\n\n## 代码示例\n\n```go\npackage main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"Hello, World! 你好，Go！\")\n}\n```\n\n## 运行方式\n\n在终端输入：\n```bash\ngo run main.go\n```\n\n你就会看到输出啦！",
				AuthorID: testUser.ID,
				Category: "技术",
				Tags:     "Go,入门,编程",
				Views:    256,
			},
			{
				Title:    "Vue3 组合式 API 使用心得",
				Content:  "# Vue3 组合式 API 心得\n\nVue3 的 `<script setup>` 真的太香了！\n\n## 对比 Options API\n\n以前写 Vue2，数据、方法、计算属性要分开写，代码长了就很乱。\n\n现在用组合式 API，可以把相关的逻辑放在一起，代码更清晰，复用也更方便。\n\n## 一个简单的例子\n\n```vue\n<script setup>\nimport { ref, computed } from 'vue'\n\nconst count = ref(0)\nconst doubleCount = computed(() => count.value * 2)\n\nfunction increment() {\n  count.value++\n}\n</script>\n```\n\n简洁又强大！",
				AuthorID: testUser.ID,
				Category: "技术",
				Tags:     "Vue3,前端,JavaScript",
				Views:    189,
			},
			{
				Title:    "春日读书随笔：《代码整洁之道》",
				Content:  "# 春日读书随笔\n\n最近天气变暖，周末在阳台上读完了《代码整洁之道》。\n\n## 书中印象最深的一句话\n\n> \"代码应该像写文章一样，有起承转合，让人读起来舒服。\"\n\n## 我的收获\n\n以前写代码只追求能跑通，现在开始慢慢注意变量命名、函数拆分这些细节了。\n\n写整洁的代码，不仅是为了别人，也是为了未来的自己。",
				AuthorID: testUser.ID,
				Category: "生活",
				Tags:     "读书,随笔,代码整洁",
				Views:    95,
			},
		}

		// 批量插入博文
		if err := config.DB.Create(&testPosts).Error; err != nil {
			log.Println("创建测试博文失败：", err)
		} else {
			log.Printf("✅ 测试博文创建成功！共插入 %d 篇\n", len(testPosts))
		}
	} else {
		log.Println("数据库中已有数据，跳过初始化测试数据。")
	}
}
