package main

import (
	"fmt"
	"log"

	"urpicbed/config"
	"urpicbed/handler"
	"urpicbed/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	if err := config.LoadConfig("config/config.yaml"); err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 创建GitHub服务
	githubService := service.NewGithubService(&config.AppConfig.Github)

	// 创建处理器
	handler := handler.NewHandler(githubService)

	// 设置Gin模式
	gin.SetMode(gin.ReleaseMode)

	// 创建Gin引擎
	r := gin.Default()

	// 添加CORS中间件
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// 设置路由
	handler.SetupRoutes(r)

	// 启动服务器
	addr := fmt.Sprintf("%s:%d", config.AppConfig.Server.Host, config.AppConfig.Server.Port)
	log.Printf("UrPicBed服务启动在: %s", addr)
	log.Printf("健康检查: http://%s/health", addr)
	log.Printf("Base64上传: POST http://%s/api/v1/upload/base64", addr)
	log.Printf("文件上传: POST http://%s/api/v1/upload/file", addr)

	if err := r.Run(addr); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}
