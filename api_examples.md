# API 使用示例

## 基础视频生成

### 1. 简单文本生成视频

```bash
curl -X POST http://localhost:8080/api/v1/video/generate \
  -H "Content-Type: application/json" \
  -d '{
    "text": "介绍人工智能的发展历程，从图灵测试到现代大语言模型",
    "style": "科技风格"
  }'
```

**响应示例：**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "task_id": "550e8400-e29b-41d4-a716-446655440000",
    "status": "processing"
  }
}
```

### 2. 查询生成状态

```bash
curl http://localhost:8080/api/v1/video/status/550e8400-e29b-41d4-a716-446655440000
```

**响应示例：**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "task_id": "550e8400-e29b-41d4-a716-446655440000",
    "status": "completed",
    "result": {
      "title": "人工智能发展历程",
      "style": "科技风格",
      "shots": [
        {
          "scene": "图灵测试的提出",
          "image_prompt": "A vintage computer room from 1950s with Alan Turing working on early computers",
          "voiceover": "1950年，阿兰·图灵提出了著名的图灵测试",
          "duration": 5,
          "clip_path": "uploads/images/image_1234567890.png",
          "voice_path": "uploads/audio/voice_1234567890.mp3",
          "subtitle": "图灵测试的诞生"
        }
      ],
      "bgm": "科技感背景音乐",
      "final": "uploads/videos/AI_Development_1234567890.mp4",
      "task_id": "550e8400-e29b-41d4-a716-446655440000",
      "status": "completed"
    }
  }
}
```

## 高级功能示例

### 1. 包含参考图片的视频生成

```bash
curl -X POST http://localhost:8080/api/v1/video/generate \
  -H "Content-Type: application/json" \
  -d '{
    "text": "展示我们公司的产品特性",
    "images": [
      "https://example.com/product1.jpg",
      "https://example.com/product2.jpg"
    ],
    "style": "商务风格",
    "audio": "https://example.com/background.mp3"
  }'
```

### 2. 获取所有任务列表

```bash
curl http://localhost:8080/api/v1/video/list
```

**响应示例：**
```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "id": 1,
      "task_id": "550e8400-e29b-41d4-a716-446655440000",
      "input": "{\"text\":\"介绍人工智能\",\"style\":\"科技风格\"}",
      "output": "{\"title\":\"人工智能发展历程\"...}",
      "created_at": "2024-01-01T10:00:00Z"
    }
  ]
}
```

### 3. 健康检查

```bash
curl http://localhost:8080/api/v1/health
```

**响应示例：**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "status": "healthy",
    "service": "video-agent-go"
  }
}
```

## 错误处理示例

### 1. 无效请求

```bash
curl -X POST http://localhost:8080/api/v1/video/generate \
  -H "Content-Type: application/json" \
  -d '{}'
```

**错误响应：**
```json
{
  "code": 400,
  "message": "Invalid request body"
}
```

### 2. 任务不存在

```bash
curl http://localhost:8080/api/v1/video/status/invalid-task-id
```

**错误响应：**
```json
{
  "code": 404,
  "message": "Task not found"
}
```

## 集成示例

### JavaScript/Node.js

```javascript
const axios = require('axios');

async function generateVideo(text, style) {
  try {
    // 创建视频生成任务
    const response = await axios.post('http://localhost:8080/api/v1/video/generate', {
      text: text,
      style: style
    });
    
    const taskId = response.data.data.task_id;
    console.log('Task created:', taskId);
    
    // 轮询检查状态
    let status = 'processing';
    while (status === 'processing') {
      await new Promise(resolve => setTimeout(resolve, 5000)); // 等待5秒
      
      const statusResponse = await axios.get(`http://localhost:8080/api/v1/video/status/${taskId}`);
      status = statusResponse.data.data.status;
      
      console.log('Current status:', status);
    }
    
    if (status === 'completed') {
      console.log('Video generated successfully!');
      console.log('Result:', statusResponse.data.data.result);
    } else {
      console.log('Video generation failed');
    }
    
  } catch (error) {
    console.error('Error:', error.response?.data || error.message);
  }
}

// 使用示例
generateVideo('创建一个关于环保的短视频', '自然风格');
```

### Python

```python
import requests
import time
import json

def generate_video(text, style):
    try:
        # 创建视频生成任务
        response = requests.post('http://localhost:8080/api/v1/video/generate', 
                               json={'text': text, 'style': style})
        response.raise_for_status()
        
        task_id = response.json()['data']['task_id']
        print(f'Task created: {task_id}')
        
        # 轮询检查状态
        status = 'processing'
        while status == 'processing':
            time.sleep(5)  # 等待5秒
            
            status_response = requests.get(f'http://localhost:8080/api/v1/video/status/{task_id}')
            status_response.raise_for_status()
            
            status = status_response.json()['data']['status']
            print(f'Current status: {status}')
        
        if status == 'completed':
            print('Video generated successfully!')
            result = status_response.json()['data']['result']
            print(f'Video URL: {result["final"]}')
        else:
            print('Video generation failed')
            
    except requests.exceptions.RequestException as e:
        print(f'Error: {e}')

# 使用示例
generate_video('创建一个关于科技创新的短视频', '现代风格')
```

## 部署后测试

如果服务部署在远程服务器上，请替换 `localhost:8080` 为实际的服务地址：

```bash
export VIDEO_API_HOST="https://your-api-domain.com"

curl $VIDEO_API_HOST/api/v1/health
``` 