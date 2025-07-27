# OpenRouter.ai Integration Example

Contoh praktis penggunaan endpoint chat completion dengan OpenRouter.ai sebagai custom provider.

## 🎯 Tentang OpenRouter.ai

OpenRouter.ai adalah platform yang menyediakan akses ke berbagai model AI melalui satu API endpoint, termasuk:
- OpenAI GPT models
- Anthropic Claude models
- Meta Llama models
- Google PaLM models
- Dan banyak lagi

## 🔧 Setup OpenRouter

### 1. Daftar di OpenRouter.ai
1. Kunjungi [https://openrouter.ai](https://openrouter.ai)
2. Buat akun dan login
3. Dapatkan API key dari dashboard
4. Top up credit untuk menggunakan models

### 2. API Key Format
```
sk-or-v1-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

## 📋 Model yang Tersedia

### Popular Models di OpenRouter:

| Model | Provider | Cost | Use Case |
|-------|----------|------|----------|
| `anthropic/claude-3-haiku` | Anthropic | $0.25/1M tokens | Fast, cheap |
| `anthropic/claude-3-sonnet` | Anthropic | $3/1M tokens | Balanced |
| `openai/gpt-3.5-turbo` | OpenAI | $0.5/1M tokens | General purpose |
| `openai/gpt-4` | OpenAI | $30/1M tokens | High quality |
| `meta-llama/llama-2-70b-chat` | Meta | $0.7/1M tokens | Open source |
| `google/palm-2-chat-bison` | Google | $0.5/1M tokens | Google's model |

## 🚀 Contoh Penggunaan

### 1. Basic Chat dengan Claude Haiku
```bash
curl -X POST http://localhost:3000/v1/examples/chat/completion \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-jwt-token" \
  -d '{
    "message": "Jelaskan apa itu REST API dalam bahasa Indonesia",
    "model": "anthropic/claude-3-haiku",
    "max_tokens": 300,
    "temperature": 0.7,
    "system_prompt": "Kamu adalah asisten AI yang membantu menjelaskan konsep teknologi dalam bahasa Indonesia.",
    "custom_endpoint": "https://openrouter.ai/api/v1",
    "custom_api_key": "sk-or-v1-your-openrouter-key"
  }'
```

### 2. Code Generation dengan GPT-4
```bash
curl -X POST http://localhost:3000/v1/examples/chat/completion \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-jwt-token" \
  -d '{
    "message": "Buatkan function Go untuk validasi email dengan regex",
    "model": "openai/gpt-4",
    "max_tokens": 500,
    "temperature": 0.3,
    "system_prompt": "Kamu adalah expert Go programmer. Berikan code yang clean dan well-documented.",
    "custom_endpoint": "https://openrouter.ai/api/v1",
    "custom_api_key": "sk-or-v1-your-openrouter-key"
  }'
```

### 3. Creative Writing dengan Llama 2
```bash
curl -X POST http://localhost:3000/v1/examples/chat/completion \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-jwt-token" \
  -d '{
    "message": "Tulis cerita pendek tentang robot yang belajar memasak",
    "model": "meta-llama/llama-2-70b-chat",
    "max_tokens": 800,
    "temperature": 0.9,
    "system_prompt": "Kamu adalah penulis kreatif yang ahli dalam menulis cerita menarik.",
    "custom_endpoint": "https://openrouter.ai/api/v1",
    "custom_api_key": "sk-or-v1-your-openrouter-key"
  }'
```

### 4. Streaming Chat dengan Claude Sonnet
```bash
curl -X POST http://localhost:3000/v1/examples/chat/completion/stream \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-jwt-token" \
  -d '{
    "message": "Jelaskan konsep machine learning dari dasar hingga advanced",
    "model": "anthropic/claude-3-sonnet",
    "max_tokens": 1000,
    "temperature": 0.6,
    "system_prompt": "Kamu adalah profesor AI yang menjelaskan konsep kompleks dengan mudah dipahami.",
    "custom_endpoint": "https://openrouter.ai/api/v1",
    "custom_api_key": "sk-or-v1-your-openrouter-key"
  }'
```

## 🎨 JavaScript Implementation

### Smart Model Selection
```javascript
class OpenRouterClient {
  constructor(apiKey) {
    this.apiKey = apiKey;
    this.endpoint = 'https://openrouter.ai/api/v1';
  }

  // Select model based on task type and budget
  selectModel(taskType, budget = 'medium') {
    const models = {
      'coding': {
        'low': 'anthropic/claude-3-haiku',
        'medium': 'openai/gpt-3.5-turbo',
        'high': 'openai/gpt-4'
      },
      'creative': {
        'low': 'meta-llama/llama-2-70b-chat',
        'medium': 'anthropic/claude-3-haiku',
        'high': 'anthropic/claude-3-sonnet'
      },
      'analysis': {
        'low': 'anthropic/claude-3-haiku',
        'medium': 'anthropic/claude-3-sonnet',
        'high': 'openai/gpt-4'
      },
      'general': {
        'low': 'anthropic/claude-3-haiku',
        'medium': 'openai/gpt-3.5-turbo',
        'high': 'anthropic/claude-3-sonnet'
      }
    };

    return models[taskType]?.[budget] || models['general'][budget];
  }

  async chat(message, taskType = 'general', budget = 'medium') {
    const model = this.selectModel(taskType, budget);
    
    const response = await fetch('/v1/examples/chat/completion', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': 'Bearer ' + localStorage.getItem('token')
      },
      body: JSON.stringify({
        message: message,
        model: model,
        max_tokens: 500,
        temperature: 0.7,
        custom_endpoint: this.endpoint,
        custom_api_key: this.apiKey
      })
    });

    const data = await response.json();
    
    if (data.status === 'success') {
      return {
        response: data.data.response,
        model: data.data.model,
        tokens: data.data.usage.total_tokens,
        processingTime: data.data.processing_time
      };
    } else {
      throw new Error(data.message);
    }
  }

  async streamChat(message, taskType = 'general', budget = 'medium') {
    const model = this.selectModel(taskType, budget);
    
    const response = await fetch('/v1/examples/chat/completion/stream', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': 'Bearer ' + localStorage.getItem('token')
      },
      body: JSON.stringify({
        message: message,
        model: model,
        max_tokens: 800,
        temperature: 0.7,
        custom_endpoint: this.endpoint,
        custom_api_key: this.apiKey
      })
    });

    return response.body.getReader();
  }
}

// Usage
const openRouter = new OpenRouterClient('sk-or-v1-your-key');

// Different use cases
async function examples() {
  // Coding task with high quality
  const codeResult = await openRouter.chat(
    'Buatkan REST API handler untuk user authentication',
    'coding',
    'high'
  );
  console.log('Code:', codeResult.response);

  // Creative task with medium budget
  const storyResult = await openRouter.chat(
    'Tulis puisi tentang teknologi',
    'creative',
    'medium'
  );
  console.log('Story:', storyResult.response);

  // Analysis task with low budget
  const analysisResult = await openRouter.chat(
    'Analisis tren teknologi 2024',
    'analysis',
    'low'
  );
  console.log('Analysis:', analysisResult.response);
}
```

### Cost Tracking
```javascript
class CostTracker {
  constructor() {
    this.usage = [];
  }

  // Estimate cost based on model and tokens
  estimateCost(model, tokens) {
    const pricing = {
      'anthropic/claude-3-haiku': 0.25 / 1000000,
      'anthropic/claude-3-sonnet': 3 / 1000000,
      'openai/gpt-3.5-turbo': 0.5 / 1000000,
      'openai/gpt-4': 30 / 1000000,
      'meta-llama/llama-2-70b-chat': 0.7 / 1000000
    };

    return (pricing[model] || 0) * tokens;
  }

  trackUsage(model, tokens, response) {
    const cost = this.estimateCost(model, tokens);
    
    this.usage.push({
      timestamp: new Date().toISOString(),
      model: model,
      tokens: tokens,
      cost: cost,
      response_length: response.length
    });

    console.log(`Used ${tokens} tokens with ${model}, estimated cost: $${cost.toFixed(6)}`);
  }

  getDailyUsage() {
    const today = new Date().toDateString();
    const todayUsage = this.usage.filter(u => 
      new Date(u.timestamp).toDateString() === today
    );

    const totalTokens = todayUsage.reduce((sum, u) => sum + u.tokens, 0);
    const totalCost = todayUsage.reduce((sum, u) => sum + u.cost, 0);

    return {
      requests: todayUsage.length,
      tokens: totalTokens,
      cost: totalCost
    };
  }
}

// Usage with cost tracking
const costTracker = new CostTracker();
const openRouter = new OpenRouterClient('sk-or-v1-your-key');

async function chatWithTracking(message, taskType, budget) {
  const result = await openRouter.chat(message, taskType, budget);
  
  costTracker.trackUsage(
    result.model,
    result.tokens,
    result.response
  );

  return result;
}
```

## 🔄 Model Comparison

### Performance Test Script
```javascript
async function compareModels(prompt) {
  const models = [
    'anthropic/claude-3-haiku',
    'anthropic/claude-3-sonnet',
    'openai/gpt-3.5-turbo',
    'openai/gpt-4'
  ];

  const results = [];

  for (const model of models) {
    const startTime = Date.now();
    
    try {
      const response = await fetch('/v1/examples/chat/completion', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'Bearer ' + localStorage.getItem('token')
        },
        body: JSON.stringify({
          message: prompt,
          model: model,
          max_tokens: 200,
          temperature: 0.7,
          custom_endpoint: 'https://openrouter.ai/api/v1',
          custom_api_key: 'sk-or-v1-your-key'
        })
      });

      const data = await response.json();
      const endTime = Date.now();

      if (data.status === 'success') {
        results.push({
          model: model,
          response: data.data.response,
          tokens: data.data.usage.total_tokens,
          responseTime: endTime - startTime,
          processingTime: data.data.processing_time
        });
      }
    } catch (error) {
      console.error(`Error with ${model}:`, error);
    }
  }

  return results;
}

// Usage
compareModels('Explain quantum computing in simple terms').then(results => {
  console.table(results);
});
```

## 💡 Best Practices untuk OpenRouter

### 1. Model Selection Strategy
```javascript
function selectOptimalModel(prompt, requirements) {
  const promptLength = prompt.length;
  const complexity = requirements.complexity || 'medium';
  const budget = requirements.budget || 'medium';
  const speed = requirements.speed || 'medium';

  // Fast and cheap for simple tasks
  if (promptLength < 100 && complexity === 'low' && budget === 'low') {
    return 'anthropic/claude-3-haiku';
  }

  // High quality for complex tasks
  if (complexity === 'high' && budget === 'high') {
    return 'openai/gpt-4';
  }

  // Balanced option
  if (speed === 'high' && budget === 'medium') {
    return 'openai/gpt-3.5-turbo';
  }

  // Creative tasks
  if (requirements.creative && budget === 'low') {
    return 'meta-llama/llama-2-70b-chat';
  }

  // Default balanced option
  return 'anthropic/claude-3-sonnet';
}
```

### 2. Error Handling & Fallback
```javascript
async function robustChat(message, options = {}) {
  const fallbackModels = [
    'anthropic/claude-3-haiku',
    'openai/gpt-3.5-turbo',
    'meta-llama/llama-2-70b-chat'
  ];

  for (const model of fallbackModels) {
    try {
      const response = await fetch('/v1/examples/chat/completion', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'Bearer ' + localStorage.getItem('token')
        },
        body: JSON.stringify({
          message: message,
          model: model,
          max_tokens: options.maxTokens || 500,
          temperature: options.temperature || 0.7,
          custom_endpoint: 'https://openrouter.ai/api/v1',
          custom_api_key: 'sk-or-v1-your-key'
        })
      });

      const data = await response.json();
      
      if (data.status === 'success') {
        return {
          success: true,
          model: model,
          response: data.data.response,
          tokens: data.data.usage.total_tokens
        };
      }
    } catch (error) {
      console.log(`Model ${model} failed, trying next...`);
    }
  }

  return {
    success: false,
    error: 'All models failed'
  };
}
```

### 3. Rate Limiting & Queuing
```javascript
class RequestQueue {
  constructor(maxConcurrent = 3, delayMs = 1000) {
    this.queue = [];
    this.running = 0;
    this.maxConcurrent = maxConcurrent;
    this.delayMs = delayMs;
  }

  async add(requestFn) {
    return new Promise((resolve, reject) => {
      this.queue.push({ requestFn, resolve, reject });
      this.process();
    });
  }

  async process() {
    if (this.running >= this.maxConcurrent || this.queue.length === 0) {
      return;
    }

    this.running++;
    const { requestFn, resolve, reject } = this.queue.shift();

    try {
      const result = await requestFn();
      resolve(result);
    } catch (error) {
      reject(error);
    } finally {
      this.running--;
      
      // Add delay between requests
      setTimeout(() => {
        this.process();
      }, this.delayMs);
    }
  }
}

// Usage
const requestQueue = new RequestQueue(2, 500); // Max 2 concurrent, 500ms delay

async function queuedChat(message, model) {
  return requestQueue.add(async () => {
    const response = await fetch('/v1/examples/chat/completion', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': 'Bearer ' + localStorage.getItem('token')
      },
      body: JSON.stringify({
        message: message,
        model: model,
        custom_endpoint: 'https://openrouter.ai/api/v1',
        custom_api_key: 'sk-or-v1-your-key'
      })
    });

    return response.json();
  });
}
```

## 📊 Monitoring & Analytics

### Usage Dashboard
```javascript
class OpenRouterDashboard {
  constructor() {
    this.stats = {
      totalRequests: 0,
      totalTokens: 0,
      totalCost: 0,
      modelUsage: {},
      errorRate: 0
    };
  }

  logRequest(model, tokens, cost, success) {
    this.stats.totalRequests++;
    
    if (success) {
      this.stats.totalTokens += tokens;
      this.stats.totalCost += cost;
      
      if (!this.stats.modelUsage[model]) {
        this.stats.modelUsage[model] = { requests: 0, tokens: 0, cost: 0 };
      }
      
      this.stats.modelUsage[model].requests++;
      this.stats.modelUsage[model].tokens += tokens;
      this.stats.modelUsage[model].cost += cost;
    } else {
      this.stats.errorRate = (this.stats.errorRate * (this.stats.totalRequests - 1) + 1) / this.stats.totalRequests;
    }
  }

  getReport() {
    return {
      summary: {
        totalRequests: this.stats.totalRequests,
        totalTokens: this.stats.totalTokens,
        totalCost: this.stats.totalCost.toFixed(6),
        errorRate: (this.stats.errorRate * 100).toFixed(2) + '%'
      },
      modelBreakdown: this.stats.modelUsage,
      recommendations: this.getRecommendations()
    };
  }

  getRecommendations() {
    const recommendations = [];
    
    // Cost optimization
    if (this.stats.totalCost > 10) {
      recommendations.push('Consider using cheaper models for simple tasks');
    }
    
    // Error rate
    if (this.stats.errorRate > 0.1) {
      recommendations.push('High error rate detected, check API keys and endpoints');
    }
    
    // Model usage
    const mostUsedModel = Object.keys(this.stats.modelUsage)
      .reduce((a, b) => this.stats.modelUsage[a].requests > this.stats.modelUsage[b].requests ? a : b);
    
    if (mostUsedModel === 'openai/gpt-4') {
      recommendations.push('Consider using GPT-3.5-turbo for cost savings');
    }
    
    return recommendations;
  }
}
```

## 🎉 Kesimpulan

OpenRouter.ai integration memberikan:
- ✅ Akses ke multiple AI models melalui satu endpoint
- ✅ Competitive pricing dan pay-per-use model
- ✅ Easy switching between different AI providers
- ✅ Cost optimization opportunities
- ✅ Fallback strategies untuk reliability

**Ready to use dengan endpoint `/v1/examples/chat/completion`!** 🚀