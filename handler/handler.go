package handler

import (
	"io"
	"log"
	"net/http"
	"strings"

	"urpicbed/models"
	"urpicbed/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	githubService *service.GithubService
}

func NewHandler(githubService *service.GithubService) *Handler {
	return &Handler{
		githubService: githubService,
	}
}

// SetupRoutes 设置路由
func (h *Handler) SetupRoutes(r *gin.Engine) {
	// 健康检查
	r.GET("/health", h.HealthCheck)

	// API路由组
	api := r.Group("/api/v1")
	{
		api.POST("/upload/base64", h.UploadBase64Image)
		api.POST("/upload/file", h.UploadFileImage)
	}
}

// HealthCheck 健康检查接口
func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "服务运行正常",
	})
}

// UploadBase64Image 上传base64图片
func (h *Handler) UploadBase64Image(c *gin.Context) {
	var req models.UploadImageRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	// 验证base64数据
	if req.Base64 == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "base64数据不能为空",
		})
		return
	}

	// 验证MIME类型
	if req.Mime == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "MIME类型不能为空",
		})
		return
	}

	// 上传到GitHub
	downloadURL, err := h.githubService.UploadBase64Image(req.Base64, req.Mime)
	if err != nil {
		log.Printf("上传base64图片失败: %v", err)
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "上传失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "上传成功",
		Data:    downloadURL,
	})
}

// UploadFileImage 上传文件图片
func (h *Handler) UploadFileImage(c *gin.Context) {
	// 获取上传的文件
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "获取文件失败: " + err.Error(),
		})
		return
	}
	defer file.Close()

	// 检查文件大小（限制为10MB）
	if header.Size > 10*1024*1024 {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "文件大小不能超过10MB",
		})
		return
	}

	// 检查文件类型
	contentType := header.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "只支持图片文件",
		})
		return
	}

	// 读取文件内容
	fileData, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "读取文件失败: " + err.Error(),
		})
		return
	}

	// 上传到GitHub
	downloadURL, err := h.githubService.UploadFileImage(fileData, contentType)
	if err != nil {
		log.Printf("上传文件图片失败: %v", err)
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "上传失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "上传成功",
		Data:    downloadURL,
	})
}
