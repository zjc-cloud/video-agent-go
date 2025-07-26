# Video Agent Go

多模态视频生成服务，使用 Golang + Hertz 框架构建的高性能视频处理 API。

## ✨ 功能特性

- 🎬 **多模态输入支持** - 文本、图片、音频、视频
- 🤖 **AI 驱动** - 基于 OpenAI GPT-4 和 DALL-E 3
- 🎙️ **语音合成** - 自动生成高质量旁白
- 🎨 **图像生成** - AI 生成视频场景图像
- 📱 **RESTful API** - 简洁易用的 HTTP 接口
- 🐳 **容器化部署** - Docker & Docker Compose 支持
- 💾 **多种存储** - 本地存储 + 云存储支持
- 📊 **任务监控** - 实时进度跟踪

## 🏗️ 项目架构

```
video-agent-go/
├── cmd/           # 程序入口
├── config/        # 配置管理
├── handler/       # HTTP 处理器
├── model/         # 数据模型
├── agent/         # 核心业务逻辑
│   ├── script.go      # 脚本生成
│   ├── image2video.go # 图片生成
│   ├── narration.go   # 语音合成
│   ├── render.go      # 视频渲染
│   ├── subtitle.go    # 字幕生成
│   └── observer.go    # 任务监控
├── storage/       # 存储抽象层
├── uploads/       # 文件上传目录
└── temp/          # 临时文件目录
```

## 🚀 快速开始

### 1. 环境准备

```bash
# 克隆项目
git clone <repository>
cd video-agent-go

# 设置环境
make setup
```

### 2. 配置环境变量

```bash
# 复制并编辑环境变量
cp .env.example .env
# 编辑 .env 文件，填入你的 OpenAI API Key
```

### 3. 启动服务

#### 方式一：Docker Compose（推荐）
```bash
make docker-run
```

#### 方式二：本地开发
```bash
# 启动依赖服务
docker-compose up -d mysql redis

# 启动应用
make dev
```

### 4. 测试 API

```bash
# 健康检查
curl http://localhost:8080/api/v1/health

# 生成视频
curl -X POST http://localhost:8080/api/v1/video/generate \
  -H "Content-Type: application/json" \
  -d '{
    "text": "创建一个关于人工智能的短视频",
    "style": "现代科技风"
  }'
```

## 📖 API 文档

### 生成视频
```http
POST /api/v1/video/generate
Content-Type: application/json

{
  "text": "视频描述文本",
  "images": ["image_url1", "image_url2"],
  "style": "视频风格",
  "audio": "背景音频URL"
}
```

### 查询任务状态
```http
GET /api/v1/video/status/{taskId}
```

### 获取所有任务
```http
GET /api/v1/video/list
```

## 🛠️ 开发指南

### 项目依赖

- Go 1.22+
- MySQL 8.0+
- Redis (可选)
- FFmpeg (视频处理)
- OpenAI API Key

### 本地开发

```bash
# 安装开发工具
make dev-tools

# 热重载开发
make dev-watch

# 代码格式化
make fmt

# 代码检查
make lint

# 运行测试
make test
```

### 构建部署

```bash
# 构建应用
make build

# 构建 Docker 镜像
make docker-build

# 生产构建
make prod-build
```

## 🔧 配置说明

### 环境变量

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| `DB_HOST` | 数据库主机 | localhost |
| `DB_PORT` | 数据库端口 | 3306 |
| `DB_USER` | 数据库用户 | root |
| `DB_PASSWORD` | 数据库密码 | - |
| `DB_NAME` | 数据库名称 | video_agent |
| `OPENAI_API_KEY` | OpenAI API 密钥 | **必填** |
| `SERVER_PORT` | 服务端口 | 8080 |
| `STORAGE_TYPE` | 存储类型 | local |

### 存储配置

支持本地存储和云存储（AWS S3）：

- `local`: 文件存储在本地 `uploads/` 目录
- `cloud`: 上传到云存储服务

## 📦 Docker 部署

### Docker Compose

```yaml
# 查看 docker-compose.yml 文件
docker-compose up -d
```

### Kubernetes

```bash
# TODO: 添加 K8s 部署配置
```

## 🧪 测试

```bash
# 运行所有测试
make test

# 运行特定包的测试
go test -v ./agent/

# 基准测试
go test -bench=. ./...
```

## 📝 TODO

- [ ] 添加更多视频风格模板
- [ ] 支持更多AI模型（Claude, Gemini等）
- [ ] 实现视频编辑功能
- [ ] 添加批量处理支持
- [ ] 性能优化和缓存
- [ ] 完善监控和日志
- [ ] 添加单元测试覆盖
- [ ] 支持更多云存储provider

## 🤝 贡献指南

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 开启 Pull Request

## 📄 许可证

该项目基于 MIT 许可证开源。查看 [LICENSE](LICENSE) 文件了解详情。

## 🆘 常见问题

### Q: OpenAI API 调用失败？
A: 请检查 API Key 是否正确，并确保账户有足够余额。

### Q: FFmpeg 相关错误？
A: 请确保系统已安装 FFmpeg，或使用 Docker 部署。

### Q: 数据库连接失败？
A: 检查数据库配置和网络连接，确保 MySQL 服务正常运行。

## 📞 联系我们

- Issue: [GitHub Issues](https://github.com/your-repo/issues)
- Email: your-email@example.com

---

⭐ 如果这个项目对你有帮助，请给个 Star！