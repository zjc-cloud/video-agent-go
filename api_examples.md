# Video Agent API Examples

本文档提供了 video-agent-go 服务的详细 API 使用示例。

## 🎯 系统架构模式

我们的系统支持三种不同的处理模式：

1. **固定工作流** - 预定义4步处理流程
2. **LLM 智能编排** - LLM 选择和协调预定义 Agent
3. **🔧 Tool-based 编排** - LLM 动态选择和组合工具

## 🔧 Tool-based 编排 (推荐)

### 1. 创建视频 - Tool-based

使用 LLM 动态选择工具的智能视频生成：

```bash
curl -X POST http://localhost:8080/api/v1/video/generate-tools \
  -H "Content-Type: application/json" \
  -d '{
    "text": "制作一个解释人工智能发展历程的教育视频",
    "style": "专业教育风格"
  }'
```

**响应示例：**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "task_id": "task_abc123",
    "status": "initializing",
    "progress": 0,
    "current_stage": "tool_selection",
    "processing_steps": [
      "Analyzing user requirements with LLM",
      "LLM selecting appropriate tools",
      "Executing tools dynamically",
      "LLM orchestrating workflow",
      "Quality validation with tools"
    ]
  }
}
```

### 2. 查看可用工具

```bash
curl http://localhost:8080/api/v1/tools/list
```

**响应示例：**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "total_tools": 6,
    "tools": [
      {
        "name": "analyze_content",
        "description": "Analyze user input to understand content requirements and suggest optimal processing strategy",
        "parameters": {
          "type": "object",
          "properties": {
            "user_text": {
              "type": "string",
              "description": "User's original request text"
            }
          },
          "required": ["user_text"]
        }
      },
      {
        "name": "generate_script",
        "description": "Generate video script based on user requirements and content analysis",
        "parameters": {
          "type": "object",
          "properties": {
            "content_type": {
              "type": "string",
              "description": "Type of content (educational, commercial, entertainment)",
              "enum": ["educational", "commercial", "entertainment", "news"]
            },
            "target_audience": {
              "type": "string",
              "description": "Target audience for the video"
            }
          },
          "required": ["content_type", "target_audience"]
        }
      },
      {
        "name": "generate_images",
        "description": "Generate images for video scenes using AI image generation",
        "parameters": {
          "type": "object",
          "properties": {
            "prompts": {
              "type": "array",
              "description": "Array of image generation prompts"
            }
          },
          "required": ["prompts"]
        }
      },
      {
        "name": "generate_voice",
        "description": "Generate voice narration for video content using text-to-speech",
        "parameters": {
          "type": "object",
          "properties": {
            "text": {
              "type": "string",
              "description": "Text content to convert to speech"
            }
          },
          "required": ["text"]
        }
      },
      {
        "name": "check_quality",
        "description": "Analyze and validate the quality of generated content",
        "parameters": {
          "type": "object",
          "properties": {
            "content_type": {
              "type": "string",
              "description": "Type of content to check",
              "enum": ["script", "images", "audio", "video"]
            }
          },
          "required": ["content_type"]
        }
      },
      {
        "name": "render_video",
        "description": "Render final video from script, images, and audio components",
        "parameters": {
          "type": "object",
          "properties": {
            "script": {
              "type": "object",
              "description": "Video script with timing information"
            }
          },
          "required": ["script"]
        }
      }
    ]
  }
}
```

### 3. 查看工具执行日志

```bash
curl http://localhost:8080/api/v1/tools/execution/task_abc123
```

**响应示例：**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "task_id": "task_abc123",
    "tool_calls": [
      {
        "tool_call_id": "call_1",
        "tool_name": "analyze_content",
        "arguments": {
          "user_text": "制作一个解释人工智能发展历程的教育视频"
        },
        "result": {
          "success": true,
          "content_type": "educational",
          "complexity": "medium",
          "target_audience": "general",
          "estimated_duration": 90,
          "key_topics": ["ai_origins", "machine_learning", "deep_learning", "current_applications"]
        },
        "timestamp": 1640995200,
        "duration_ms": 500
      },
      {
        "tool_call_id": "call_2",
        "tool_name": "generate_script",
        "arguments": {
          "content_type": "educational",
          "target_audience": "general",
          "style": "professional",
          "duration": 90,
          "key_points": ["ai_origins", "machine_learning", "deep_learning", "current_applications"]
        },
        "result": {
          "success": true,
          "title": "AI Development History",
          "structure": ["introduction", "origins", "evolution", "current_state", "conclusion"],
          "estimated_duration": 95
        },
        "timestamp": 1640995205,
        "duration_ms": 2500
      },
      {
        "tool_call_id": "call_3",
        "tool_name": "generate_images",
        "arguments": {
          "prompts": [
            "Historical timeline of AI development",
            "Early computer systems and algorithms",
            "Modern neural networks visualization",
            "Current AI applications in daily life"
          ],
          "style": "educational_infographic",
          "resolution": "1920x1080"
        },
        "result": {
          "success": true,
          "images": [
            {"id": "img_1", "url": "/uploads/images/ai_timeline.jpg"},
            {"id": "img_2", "url": "/uploads/images/early_computers.jpg"},
            {"id": "img_3", "url": "/uploads/images/neural_networks.jpg"},
            {"id": "img_4", "url": "/uploads/images/ai_applications.jpg"}
          ],
          "count": 4
        },
        "timestamp": 1640995210,
        "duration_ms": 8000
      },
      {
        "tool_call_id": "call_4",
        "tool_name": "generate_voice",
        "arguments": {
          "text": "人工智能的发展历程是一个充满创新与突破的故事...",
          "voice_type": "professional",
          "language": "zh-CN",
          "emotion": "educational"
        },
        "result": {
          "success": true,
          "audio_file": "/uploads/audio/ai_history_narration.mp3",
          "duration": 95,
          "voice_type": "professional"
        },
        "timestamp": 1640995220,
        "duration_ms": 3100
      },
      {
        "tool_call_id": "call_5",
        "tool_name": "check_quality",
        "arguments": {
          "content_type": "video",
          "content_data": {
            "script": "...",
            "images": "...",
            "audio": "..."
          }
        },
        "result": {
          "success": true,
          "quality_scores": {
            "overall": 0.92,
            "accuracy": 0.95,
            "clarity": 0.88,
            "engagement": 0.90
          },
          "passed": true,
          "recommendations": ["Consider adding more visual transitions"]
        },
        "timestamp": 1640995225,
        "duration_ms": 1200
      },
      {
        "tool_call_id": "call_6",
        "tool_name": "render_video",
        "arguments": {
          "script": "...",
          "images": ["img_1", "img_2", "img_3", "img_4"],
          "audio": "/uploads/audio/ai_history_narration.mp3",
          "effects": ["fade_transitions", "text_overlay"]
        },
        "result": {
          "success": true,
          "video_file": "/uploads/videos/ai_history_educational.mp4",
          "duration": 95,
          "resolution": "1920x1080",
          "file_size": "45.2MB"
        },
        "timestamp": 1640995230,
        "duration_ms": 12300
      }
    ],
    "total_calls": 6,
    "execution_model": "LLM-driven tool selection",
    "total_duration": "27.6s"
  }
}
```

## 🧠 LLM 智能编排

### 创建智能视频

```bash
curl -X POST http://localhost:8080/api/v1/video/generate-smart \
  -H "Content-Type: application/json" \
  -d '{
    "text": "制作一个介绍我们新产品的宣传视频",
    "style": "现代商务风格",
    "images": ["product1.jpg", "product2.jpg"],
    "custom_scripts": [
      {
        "language": "python",
        "code": "# 产品特效处理\nenhance_product_visuals()",
        "stage": "post_processing"
      }
    ]
  }'
```

### 查看可用 Agent

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
        "name": "VoiceGenerator",
        "description": "Generates voiceovers and narration using TTS",
        "capabilities": ["voice_synthesis", "narration", "multilingual_tts"]
      },
      {
        "name": "VideoRender",
        "description": "Renders and combines media into final video",
        "capabilities": ["video_rendering", "media_composition", "format_conversion"]
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
      },
      {
        "name": "Optimization",
        "description": "Optimizes and improves content quality",
        "capabilities": ["content_optimization", "performance_enhancement", "quality_improvement"]
      }
    ]
  }
}
```

## 📊 传统固定流程

### 基础视频生成

```bash
curl -X POST http://localhost:8080/api/v1/video/generate \
  -H "Content-Type: application/json" \
  -d '{
    "text": "制作一个关于环保的短视频",
    "style": "清新自然风格"
  }'
```

## 📈 任务状态查询

```bash
curl http://localhost:8080/api/v1/video/status/task_abc123
```

**响应示例：**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "task_id": "task_abc123",
    "status": "completed",
    "result": {
      "task_id": "task_abc123",
      "title": "AI Development History",
      "status": "completed",
      "final": "/uploads/videos/ai_history_educational.mp4",
      "duration": 95,
      "quality_score": 0.92
    }
  }
}
```

## 📋 任务列表查询

```bash
curl http://localhost:8080/api/v1/video/list
```

## 🏥 健康检查

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
    "modes": "tool-based"
  }
}
```

## 🔧 JavaScript 集成示例

### Tool-based 模式

```javascript
class VideoAgentClient {
  constructor(baseURL = 'http://localhost:8080') {
    this.baseURL = baseURL;
  }

  // 使用 Tool-based 智能编排生成视频
  async generateVideoWithTools(request) {
    const response = await fetch(`${this.baseURL}/api/v1/video/generate-tools`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(request),
    });
    return response.json();
  }

  // 查看可用工具
  async getAvailableTools() {
    const response = await fetch(`${this.baseURL}/api/v1/tools/list`);
    return response.json();
  }

  // 查看工具执行日志
  async getToolExecutionLog(taskId) {
    const response = await fetch(`${this.baseURL}/api/v1/tools/execution/${taskId}`);
    return response.json();
  }

  // 轮询任务状态直到完成
  async waitForCompletion(taskId, pollInterval = 2000) {
    while (true) {
      const statusResponse = await this.getTaskStatus(taskId);
      const status = statusResponse.data;
      
      console.log(`Task ${taskId} status: ${status.status}`);
      
      if (status.status === 'completed') {
        return status.result;
      } else if (status.status === 'failed') {
        throw new Error('Task failed');
      }
      
      await new Promise(resolve => setTimeout(resolve, pollInterval));
    }
  }

  async getTaskStatus(taskId) {
    const response = await fetch(`${this.baseURL}/api/v1/video/status/${taskId}`);
    return response.json();
  }
}

// 使用示例
async function createEducationalVideo() {
  const client = new VideoAgentClient();
  
  try {
    // 1. 先查看可用工具
    const toolsResponse = await client.getAvailableTools();
    console.log('Available tools:', toolsResponse.data.tools.map(t => t.name));
    
    // 2. 创建视频任务
    const response = await client.generateVideoWithTools({
      text: "制作一个关于量子计算基础概念的教育视频",
      style: "学术教育风格"
    });
    
    const taskId = response.data.task_id;
    console.log(`Task created: ${taskId}`);
    
    // 3. 等待完成
    const result = await client.waitForCompletion(taskId);
    console.log('Video completed:', result.final);
    
    // 4. 查看详细的工具执行日志
    const logResponse = await client.getToolExecutionLog(taskId);
    console.log('Tool execution summary:');
    logResponse.data.tool_calls.forEach(call => {
      console.log(`- ${call.tool_name}: ${call.duration_ms}ms`);
    });
    
  } catch (error) {
    console.error('Error:', error);
  }
}

// 执行
createEducationalVideo();
```

## 🐍 Python 集成示例

### Tool-based 模式

```python
import requests
import time
import json

class VideoAgentClient:
    def __init__(self, base_url="http://localhost:8080"):
        self.base_url = base_url
    
    def generate_video_with_tools(self, request_data):
        """使用 Tool-based 智能编排生成视频"""
        response = requests.post(
            f"{self.base_url}/api/v1/video/generate-tools",
            json=request_data
        )
        return response.json()
    
    def get_available_tools(self):
        """获取可用工具列表"""
        response = requests.get(f"{self.base_url}/api/v1/tools/list")
        return response.json()
    
    def get_tool_execution_log(self, task_id):
        """获取工具执行日志"""
        response = requests.get(f"{self.base_url}/api/v1/tools/execution/{task_id}")
        return response.json()
    
    def get_task_status(self, task_id):
        """获取任务状态"""
        response = requests.get(f"{self.base_url}/api/v1/video/status/{task_id}")
        return response.json()
    
    def wait_for_completion(self, task_id, poll_interval=2):
        """等待任务完成"""
        while True:
            status_response = self.get_task_status(task_id)
            status = status_response['data']['status']
            
            print(f"Task {task_id} status: {status}")
            
            if status == 'completed':
                return status_response['data']['result']
            elif status == 'failed':
                raise Exception('Task failed')
            
            time.sleep(poll_interval)

def create_commercial_video():
    """创建商业宣传视频示例"""
    client = VideoAgentClient()
    
    try:
        # 1. 查看可用工具
        tools_response = client.get_available_tools()
        tool_names = [tool['name'] for tool in tools_response['data']['tools']]
        print(f"Available tools: {tool_names}")
        
        # 2. 创建视频任务
        response = client.generate_video_with_tools({
            "text": "为我们的新款智能手表制作产品发布视频",
            "style": "现代科技风格"
        })
        
        task_id = response['data']['task_id']
        print(f"Task created: {task_id}")
        
        # 3. 等待完成
        result = client.wait_for_completion(task_id)
        print(f"Video completed: {result['final']}")
        
        # 4. 分析工具使用情况
        log_response = client.get_tool_execution_log(task_id)
        tool_calls = log_response['data']['tool_calls']
        
        print("\n🔧 Tool Execution Analysis:")
        total_duration = 0
        for call in tool_calls:
            duration = call['duration_ms']
            total_duration += duration
            print(f"  • {call['tool_name']}: {duration}ms")
            
        print(f"  📊 Total processing time: {total_duration}ms")
        print(f"  🛠️  Tools used: {len(tool_calls)}")
        print(f"  🎯 Execution model: {log_response['data']['execution_model']}")
        
    except Exception as error:
        print(f"Error: {error}")

if __name__ == "__main__":
    create_commercial_video()
```

## 🔍 工具执行流程分析

基于 Tool-based 模式，LLM 的典型工具选择流程：

### 教育内容
```
analyze_content → generate_script → generate_images → generate_voice → check_quality → render_video
```

### 商业内容  
```
analyze_content → generate_images → generate_script → generate_voice → render_video
```

### 质量问题自动修复
```
check_quality → [发现问题] → generate_script + generate_images → check_quality → render_video
```

## 🎯 最佳实践

1. **优先使用 Tool-based 模式** - 提供最大灵活性和最佳结果
2. **监控工具执行日志** - 了解 LLM 的决策过程
3. **根据内容类型选择合适模式** - 简单任务可用固定流程
4. **关注质量分数** - Tool-based 模式通常能达到更高质量

---

通过这些 API，您可以充分利用我们的 **Tool-based AI Agent** 系统，实现真正智能化的视频生成！🔧🚀 