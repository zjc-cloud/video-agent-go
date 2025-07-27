# API 使用示例

## 🚀 新功能：LLM 驱动的智能视频生成

### 1. 智能视频生成（推荐）

使用 LLM 自动分析需求并选择最优的处理流程：

```bash
curl -X POST http://localhost:8080/api/v1/video/generate-smart \
  -H "Content-Type: application/json" \
  -d '{
    "text": "制作一个介绍人工智能发展历程的教育视频，从图灵测试到现代大语言模型",
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
    "status": "analyzing",
    "progress": 0,
    "current_stage": "task_analysis",
    "processing_steps": [
      "Analyzing user requirements",
      "Generating execution plan", 
      "Selecting optimal agents",
      "Dynamic execution"
    ]
  }
}
```

### 2. 查看可用的智能体

```bash
curl http://localhost:8080/api/v1/agents/list
```

**响应示例：**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "total_agents": 7,
    "agents": [
      {
        "name": "ScriptGenerator",
        "description": "Generates video scripts and storyboards using AI",
        "capabilities": ["script_generation", "storyboard_creation", "narrative_structure"]
      },
      {
        "name": "ImageGenerator", 
        "description": "Generates images using DALL-E based on prompts",
        "capabilities": ["image_generation", "visual_creativity", "scene_creation"]
      },
      {
        "name": "Analysis",
        "description": "Analyzes content and provides insights for optimization",
        "capabilities": ["content_analysis", "sentiment_analysis", "engagement_prediction"]
      },
      {
        "name": "QualityCheck",
        "description": "Validates output quality and identifies issues", 
        "capabilities": ["quality_validation", "error_detection", "compliance_check"]
      }
    ]
  }
}
```

### 3. 查看执行日志

查看 LLM 如何智能选择和执行各个子 agent：

```bash
curl http://localhost:8080/api/v1/execution/log/550e8400-e29b-41d4-a716-446655440000
```

**响应示例：**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "task_id": "550e8400-e29b-41d4-a716-446655440000",
    "execution_log": [
      {
        "step_id": "step_1",
        "agent_name": "Analysis",
        "action": "analyze_content",
        "timestamp": 1640995200,
        "duration_ms": 1500,
        "success": true,
        "message": "Content analysis completed: educational content, intermediate complexity"
      },
      {
        "step_id": "step_2", 
        "agent_name": "ScriptGenerator",
        "action": "generate_script",
        "timestamp": 1640995205,
        "duration_ms": 3000,
        "success": true,
        "message": "Generated script with 5 shots focusing on AI milestones"
      },
      {
        "step_id": "step_3",
        "agent_name": "ImageGenerator",
        "action": "generate_images",
        "timestamp": 1640995210,
        "duration_ms": 8000, 
        "success": true,
        "message": "Generated 5 high-quality images for each AI milestone"
      },
      {
        "step_id": "step_4",
        "agent_name": "QualityCheck",
        "action": "validate_quality",
        "timestamp": 1640995220,
        "duration_ms": 2000,
        "success": true,
        "message": "Quality score: 0.92 - Excellent quality, suggested adding subtitles"
      },
      {
        "step_id": "step_5",
        "agent_name": "Optimization", 
        "action": "enhance_content",
        "timestamp": 1640995225,
        "duration_ms": 4000,
        "success": true,
        "message": "Applied color correction and enhanced transitions"
      }
    ],
    "total_steps": 5
  }
}
```

## 🆚 对比：固定流程 vs LLM 驱动

### 固定流程（原有方式）

```bash
curl -X POST http://localhost:8080/api/v1/video/generate \
  -H "Content-Type: application/json" \
  -d '{
    "text": "介绍人工智能的发展历程",
    "style": "科技风格"
  }'
```

**执行流程（固定）：**
1. ✅ 生成脚本
2. ✅ 生成图像  
3. ✅ 生成语音
4. ✅ 渲染视频
5. ✅ 完成

### LLM 驱动（智能方式）

```bash
curl -X POST http://localhost:8080/api/v1/video/generate-smart \
  -H "Content-Type: application/json" \
  -d '{
    "text": "介绍人工智能的发展历程",
    "style": "科技风格"
  }'
```

**执行流程（动态）：**
1. 🧠 **LLM 分析任务** - "这是教育内容，需要专业性和易懂性"
2. 🔍 **选择分析代理** - 深入理解用户需求
3. 📝 **智能脚本生成** - 基于分析结果优化脚本结构
4. 🎨 **图像生成** - 针对AI发展历程选择合适的视觉风格
5. 🎙️ **语音合成** - 选择适合教育内容的语音风格
6. ✅ **质量检查** - 自动验证内容质量
7. ⚡ **智能优化** - 根据检查结果进行优化
8. 🎬 **最终渲染** - 生成高质量视频

## 🌟 智能特性示例

### 1. 自适应工作流

不同类型的视频请求，LLM 会选择不同的处理流程：

#### 新闻类视频
```json
{
  "text": "制作一个关于今日科技新闻的视频",
  "style": "新闻播报"
}
```
**LLM 选择的流程：**
- Analysis → ScriptGenerator → VoiceGenerator → VideoRender
- *跳过图像生成，重点关注语音质量*

#### 产品展示视频  
```json
{
  "text": "展示我们的新款智能手机特性",
  "style": "商务风格",
  "images": ["product1.jpg", "product2.jpg"]
}
```  
**LLM 选择的流程：**
- Analysis → ImageGenerator → ScriptGenerator → VoiceGenerator → VideoRender → QualityCheck
- *重点处理图像，确保产品展示效果*

#### 教育内容
```json
{
  "text": "解释量子计算的基本原理", 
  "style": "教育风格"
}
```
**LLM 选择的流程：**
- Analysis → ScriptGenerator → ImageGenerator → VoiceGenerator → QualityCheck → Optimization → VideoRender
- *全流程处理，确保教育内容的准确性和易懂性*

### 2. 动态质量优化

LLM 会根据中间结果动态调整后续步骤：

```bash
# 如果质量检查发现问题
{
  "step_id": "quality_check_1",
  "agent_name": "QualityCheck", 
  "result": {
    "quality_score": 0.65,  // 低于阈值
    "issues": ["audio_quality_low", "image_blur"]
  }
}

# LLM 自动决定重新处理
{
  "step_id": "adaptive_reprocess",
  "agent_name": "VoiceGenerator",
  "action": "regenerate_with_higher_quality"
}
```

### 3. 智能错误恢复

```bash
# 如果某个 agent 失败
{
  "step_id": "image_gen_1",
  "agent_name": "ImageGenerator",
  "success": false,
  "error": "DALL-E API rate limit exceeded"
}

# LLM 自动选择备用方案
{
  "step_id": "fallback_strategy",
  "agent_name": "ImageGenerator", 
  "action": "use_stock_images_with_custom_prompts"
}
```

## 📊 性能对比

| 特性 | 固定流程 | LLM 驱动 |
|------|----------|----------|
| 执行步骤 | 4步（固定） | 3-8步（动态） |
| 质量检查 | ❌ 无 | ✅ 自动 |
| 错误恢复 | ❌ 失败即停止 | ✅ 智能重试 |
| 个性化 | ❌ 统一处理 | ✅ 内容自适应 |
| 优化能力 | ❌ 无 | ✅ 自动优化 |
| 可观测性 | ❌ 简单日志 | ✅ 详细执行日志 |

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
    "service": "video-agent-go",
    "mode": "llm-orchestrated"
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

async function generateVideoSmart(text, style) {
  try {
    // 使用智能生成接口
    const response = await axios.post('http://localhost:8080/api/v1/video/generate-smart', {
      text: text,
      style: style
    });
    
    const taskId = response.data.data.task_id;
    console.log('Smart task created:', taskId);
    
    // 轮询检查状态
    let status = 'analyzing';
    while (status === 'analyzing' || status === 'processing') {
      await new Promise(resolve => setTimeout(resolve, 5000));
      
      const statusResponse = await axios.get(`http://localhost:8080/api/v1/video/status/${taskId}`);
      status = statusResponse.data.data.status;
      
      console.log('Current status:', status);
      
      // 查看执行日志
      if (status === 'processing') {
        const logResponse = await axios.get(`http://localhost:8080/api/v1/execution/log/${taskId}`);
        console.log('Execution steps:', logResponse.data.data.execution_log.length);
      }
    }
    
    if (status === 'completed') {
      console.log('Video generated successfully with LLM orchestration!');
    }
    
  } catch (error) {
    console.error('Error:', error.response?.data || error.message);
  }
}

// 使用示例
generateVideoSmart('创建一个关于环保的短视频', '自然风格');
```

### Python

```python
import requests
import time
import json

def generate_video_smart(text, style):
    try:
        # 使用智能生成接口
        response = requests.post('http://localhost:8080/api/v1/video/generate-smart', 
                               json={'text': text, 'style': style})
        response.raise_for_status()
        
        task_id = response.json()['data']['task_id']
        print(f'Smart task created: {task_id}')
        
        # 轮询检查状态
        status = 'analyzing'
        while status in ['analyzing', 'processing']:
            time.sleep(5)
            
            status_response = requests.get(f'http://localhost:8080/api/v1/video/status/{task_id}')
            status_response.raise_for_status()
            
            status = status_response.json()['data']['status']
            print(f'Current status: {status}')
            
            # 查看智能体执行情况
            if status == 'processing':
                log_response = requests.get(f'http://localhost:8080/api/v1/execution/log/{task_id}')
                if log_response.status_code == 200:
                    steps = log_response.json()['data']['execution_log']
                    print(f'Executed {len(steps)} intelligent steps')
                    for step in steps[-3:]:  # 显示最近3步
                        print(f"  - {step['agent_name']}: {step['message']}")
        
        if status == 'completed':
            print('Video generated successfully using LLM orchestration!')
            
    except requests.exceptions.RequestException as e:
        print(f'Error: {e}')

# 使用示例
generate_video_smart('创建一个关于科技创新的短视频', '现代风格')
```

## 部署后测试

如果服务部署在远程服务器上，请替换 `localhost:8080` 为实际的服务地址：

```bash
export VIDEO_API_HOST="https://your-api-domain.com"

# 测试智能生成
curl $VIDEO_API_HOST/api/v1/video/generate-smart \
  -H "Content-Type: application/json" \
  -d '{"text": "测试智能视频生成", "style": "现代风格"}'

# 查看可用智能体
curl $VIDEO_API_HOST/api/v1/agents/list
``` 