package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

const baseURL = "http://localhost:8080"

// 测试用的1x1像素PNG图片的base64编码
const testImageBase64 = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg=="

func main() {
	fmt.Println("开始测试UrPicBed API...")

	// 测试健康检查
	fmt.Println("\n1. 测试健康检查...")
	if err := testHealthCheck(); err != nil {
		fmt.Printf("健康检查失败: %v\n", err)
		return
	}
	fmt.Println("✓ 健康检查通过")

	// 测试Base64上传
	fmt.Println("\n2. 测试Base64图片上传...")
	if err := testBase64Upload(); err != nil {
		fmt.Printf("Base64上传失败: %v\n", err)
		return
	}
	fmt.Println("✓ Base64上传测试通过")

	// 测试文件上传
	fmt.Println("\n3. 测试文件上传...")
	if err := testFileUpload(); err != nil {
		fmt.Printf("文件上传失败: %v\n", err)
		return
	}
	fmt.Println("✓ 文件上传测试通过")

	fmt.Println("\n🎉 所有测试通过！")
}

func testHealthCheck() error {
	resp, err := http.Get(baseURL + "/health")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("期望状态码200，实际得到%d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	if success, ok := result["success"].(bool); !ok || !success {
		return fmt.Errorf("健康检查返回失败状态")
	}

	return nil
}

func testBase64Upload() error {
	reqBody := map[string]string{
		"base64": testImageBase64,
		"mime":   "image/png",
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	resp, err := http.Post(baseURL+"/api/v1/upload/base64", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("期望状态码200，实际得到%d，响应: %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	if success, ok := result["success"].(bool); !ok || !success {
		return fmt.Errorf("上传返回失败状态: %v", result)
	}

	if data, ok := result["data"].(string); ok && data != "" {
		fmt.Printf("✓ 上传成功，获得链接: %s\n", data)
	}

	return nil
}

func testFileUpload() error {
	// 创建测试图片文件
	testImageData, err := base64.StdEncoding.DecodeString(testImageBase64)
	if err != nil {
		return err
	}

	// 创建临时文件
	tmpFile, err := os.CreateTemp("", "test_image_*.png")
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// 写入测试图片数据
	if _, err := tmpFile.Write(testImageData); err != nil {
		return err
	}

	// 重新打开文件用于上传
	file, err := os.Open(tmpFile.Name())
	if err != nil {
		return err
	}
	defer file.Close()

	// 创建multipart请求
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	part, err := writer.CreateFormFile("file", "test.png")
	if err != nil {
		return err
	}

	if _, err := io.Copy(part, file); err != nil {
		return err
	}

	writer.Close()

	// 发送请求
	resp, err := http.Post(baseURL+"/api/v1/upload/file", writer.FormDataContentType(), &buf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("期望状态码200，实际得到%d，响应: %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	if success, ok := result["success"].(bool); !ok || !success {
		return fmt.Errorf("上传返回失败状态: %v", result)
	}

	if data, ok := result["data"].(string); ok && data != "" {
		fmt.Printf("✓ 上传成功，获得链接: %s\n", data)
	}

	return nil
}
