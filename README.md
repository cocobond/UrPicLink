# UrPicBed - GitHub图床服务

UrPicBed是一个基于Golang开发的图片上传服务，支持将图片上传到GitHub仓库并返回外链链接。通过这种方式，你可以轻松地将本地图片转换为可分享的网络链接。

## 功能特性

- 🖼️ 支持上传图片文件（JPG、PNG、GIF、WebP、BMP）
- 📝 支持上传Base64编码的图片
- 🔄 自动负载均衡到多个GitHub仓库
- 🐳 Docker一键部署
- 🔧 灵活的配置文件
- 🌐 跨域支持
- 📊 健康检查接口

## 快速开始

### 使用Docker Compose（推荐）

1. 克隆项目并进入目录：
```bash
cd UrPicBed
```

2. 修改配置文件 `config/config.yaml`：
```yaml
server:
  port: 8080
  host: "0.0.0.0"

github:
  commit-url-list:
    - "https://api.github.com/repos/你的用户名/仓库名1/contents/"
    - "https://api.github.com/repos/你的用户名/仓库名2/contents/"
  
  authorization: "Bearer 你的GitHub_TOKEN"
  accept: "application/vnd.github+json"
  api-version: "2022-11-28"
```

3. 启动服务：
```bash
docker-compose up -d
```

4. 访问健康检查：
```bash
curl http://localhost:8080/health
```

### 手动构建

1. 安装Go 1.21+
2. 下载依赖：
```bash
go mod download
```

3. 构建项目：
```bash
go build -o urpicbed .
```

4. 运行服务：
```bash
./urpicbed
```

## API文档

### 健康检查

**GET** `/health`

检查服务状态

**响应示例：**
```json
{
  "success": true,
  "message": "服务运行正常"
}
```

### 上传Base64图片

**POST** `/api/v1/upload/base64`

**请求体：**
```json
{
  "base64": "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg==",
  "mime": "image/png"
}
```

**响应示例：**
```json
{
  "success": true,
  "message": "上传成功",
  "data": "https://raw.githubusercontent.com/用户名/仓库名/main/1703123456.png"
}
```

### 上传图片文件

**POST** `/api/v1/upload/file`

**请求体：** `multipart/form-data`

- `file`: 图片文件（支持JPG、PNG、GIF、WebP、BMP，最大10MB）

**响应示例：**
```json
{
  "success": true,
  "message": "上传成功",
  "data": "https://raw.githubusercontent.com/用户名/仓库名/main/1703123456.jpg"
}
```

## 配置说明

### GitHub配置

- `commit-url-list`: GitHub仓库的API URL列表，支持多个仓库实现负载均衡
- `authorization`: GitHub个人访问令牌（Personal Access Token）
- `accept`: GitHub API的Accept头
- `api-version`: GitHub API版本

### 服务器配置

- `port`: 服务端口
- `host`: 服务主机地址

## 获取GitHub Token

1. 访问 [GitHub Settings > Developer settings > Personal access tokens](https://github.com/settings/tokens)
2. 点击 "Generate new token (classic)"
3. 选择权限：
   - `repo` (完整的仓库访问权限)
4. 生成并复制Token
5. 在配置文件中使用：`Bearer YOUR_TOKEN`

## 创建GitHub仓库

1. 在GitHub上创建新的公开仓库
2. 仓库名称建议格式：`picbed-01`, `picbed-02` 等
3. 将仓库URL添加到配置文件的 `commit-url-list` 中

## 使用示例

### 使用curl上传Base64图片

```bash
curl -X POST http://localhost:8080/api/v1/upload/base64 \
  -H "Content-Type: application/json" \
  -d '{
    "base64": "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg==",
    "mime": "image/png"
  }'
```

### 使用curl上传图片文件

```bash
curl -X POST http://localhost:8080/api/v1/upload/file \
  -F "file=@/path/to/your/image.jpg"
```

### 使用JavaScript上传

```javascript
// 上传Base64图片
async function uploadBase64(base64Data, mimeType) {
  const response = await fetch('http://localhost:8080/api/v1/upload/base64', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      base64: base64Data,
      mime: mimeType
    })
  });
  
  return await response.json();
}

// 上传文件
async function uploadFile(file) {
  const formData = new FormData();
  formData.append('file', file);
  
  const response = await fetch('http://localhost:8080/api/v1/upload/file', {
    method: 'POST',
    body: formData
  });
  
  return await response.json();
}
```

## 注意事项

1. **GitHub Token安全**: 请妥善保管你的GitHub Token，不要提交到代码仓库
2. **仓库权限**: 确保GitHub Token有对应仓库的写入权限
3. **文件大小**: 单个文件最大支持10MB
4. **文件类型**: 只支持图片文件格式
5. **负载均衡**: 系统会随机选择仓库进行上传，实现负载均衡

## 故障排除

### 常见错误

1. **401 Unauthorized**: GitHub Token无效或过期
2. **403 Forbidden**: Token权限不足
3. **404 Not Found**: 仓库URL错误或仓库不存在
4. **413 Payload Too Large**: 文件超过10MB限制

### 日志查看

```bash
# 查看Docker容器日志
docker-compose logs -f urpicbed

# 查看应用日志
docker logs -f urpicbed
```

## 许可证

MIT License
