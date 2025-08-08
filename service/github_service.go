package service

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

	"urpicbed/config"
	"urpicbed/models"
)

type GithubService struct {
	config *config.GithubConfig
}

func NewGithubService(cfg *config.GithubConfig) *GithubService {
	return &GithubService{
		config: cfg,
	}
}

// UploadBase64Image 上传base64图片到GitHub
func (s *GithubService) UploadBase64Image(base64Data, mimeType string) (string, error) {
	// 解码base64数据
	decodedData, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return "", fmt.Errorf("base64解码失败: %v", err)
	}

	// 重新编码为base64（GitHub API要求）
	encodedData := base64.StdEncoding.EncodeToString(decodedData)

	// 生成文件名
	filename := s.generateFilename(mimeType)

	// 随机选择一个仓库
	repoURL := s.getRandomRepoURL()

	// 构建完整的API URL
	apiURL := repoURL + filename

	// 创建GitHub API请求
	githubReq := models.GithubAPIRequest{
		Message: "Upload image via UrPicBed",
		Content: encodedData,
	}

	// 序列化请求
	reqBody, err := json.Marshal(githubReq)
	if err != nil {
		return "", fmt.Errorf("请求序列化失败: %v", err)
	}

	// 创建HTTP请求
	req, err := http.NewRequest("PUT", apiURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Accept", s.config.Accept)
	req.Header.Set("Authorization", s.config.Authorization)
	req.Header.Set("X-GitHub-Api-Version", s.config.APIVersion)
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %v", err)
	}

	// 检查响应状态
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API错误: %s, 响应: %s", resp.Status, string(respBody))
	}

	// 解析响应
	var githubResp models.GithubAPIResponse
	if err := json.Unmarshal(respBody, &githubResp); err != nil {
		return "", fmt.Errorf("解析响应失败: %v", err)
	}

	return githubResp.Content.DownloadURL, nil
}

// UploadFileImage 上传文件图片到GitHub
func (s *GithubService) UploadFileImage(fileData []byte, mimeType string) (string, error) {
	// 编码为base64
	encodedData := base64.StdEncoding.EncodeToString(fileData)

	// 生成文件名
	filename := s.generateFilename(mimeType)

	// 随机选择一个仓库
	repoURL := s.getRandomRepoURL()

	// 构建完整的API URL
	apiURL := repoURL + filename

	// 创建GitHub API请求
	githubReq := models.GithubAPIRequest{
		Message: "Upload image via UrPicBed",
		Content: encodedData,
	}

	// 序列化请求
	reqBody, err := json.Marshal(githubReq)
	if err != nil {
		return "", fmt.Errorf("请求序列化失败: %v", err)
	}

	// 创建HTTP请求
	req, err := http.NewRequest("PUT", apiURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Accept", s.config.Accept)
	req.Header.Set("Authorization", s.config.Authorization)
	req.Header.Set("X-GitHub-Api-Version", s.config.APIVersion)
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %v", err)
	}

	// 检查响应状态
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API错误: %s, 响应: %s", resp.Status, string(respBody))
	}

	// 解析响应
	var githubResp models.GithubAPIResponse
	if err := json.Unmarshal(respBody, &githubResp); err != nil {
		return "", fmt.Errorf("解析响应失败: %v", err)
	}

	return githubResp.Content.DownloadURL, nil
}

// generateFilename 生成文件名
func (s *GithubService) generateFilename(mimeType string) string {
	timestamp := time.Now().Unix()
	extension := s.getExtensionFromMime(mimeType)
	return fmt.Sprintf("%d%s", timestamp, extension)
}

// getExtensionFromMime 根据MIME类型获取文件扩展名
func (s *GithubService) getExtensionFromMime(mimeType string) string {
	switch mimeType {
	case "image/jpeg", "image/jpg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/gif":
		return ".gif"
	case "image/webp":
		return ".webp"
	case "image/bmp":
		return ".bmp"
	default:
		return ".jpg"
	}
}

// getRandomRepoURL 随机选择一个仓库URL
func (s *GithubService) getRandomRepoURL() string {
	if len(s.config.CommitURLList) == 0 {
		return ""
	}
	rand.Seed(time.Now().UnixNano())
	return s.config.CommitURLList[rand.Intn(len(s.config.CommitURLList))]
}
