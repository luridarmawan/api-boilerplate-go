# Custom Endpoint & API Key Guide

Panduan penggunaan custom endpoint dan API key untuk AI chat completion di module example.

## 🎯 Fitur Custom Endpoint

Endpoint chat completion sekarang mendukung penggunaan custom endpoint dan API key, memungkinkan Anda untuk:
- Menggunakan provider AI yang berbeda per request
- Beralih antara OpenAI, OpenRouter, Anthropic, atau provider lainnya
- Menggunakan API key yang berbeda untuk setiap request

## 📋 Parameter Baru

### Request Body
```json
{
  "message": "Your message here",
  "model": "gpt-3.5-turbo",
  "max_tokens": 500,
  "temperature": 0.7,
  "system_prompt": "Custom system prompt",
  "custom_endpoint": "https://openrouter.ai/api/v1",
  "custom_api_key": "sk-or-v1-your-openrouter-key"
}
```

### Parameter Details

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `custom_endpoint` | string | ❌ | Custom API endpoint URL |
| `custom_api_key` | string | ❌ | Custom API key for the endpoint |

**Important**: Jika salah satu dari `custom_endpoint` atau `custom_api_key` disediakan, maka keduanya harus disediakan.

## 🚀 Contoh Penggunaan

### 1. OpenRouter.ai
```bash
curl -X POST http://localhost:3000/v1/examples/chat/completion \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-jwt-token" \
  -d '{
    "message": "Explain quantum computing",
    "model": "anthropic/claude-3-haiku",
    "max_tokens": 300,
    "temperature": 0.7,
    "custom_endpoint": "https://openrouter.ai/api/v1",
    "custom_api_key": "sk-or-v1-your-openrouter-key"
  }'
```

### 2. Anthropic Claude
```bash
curl -X POST http://localhost:3000/v1/examples/chat/completion \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-jwt-token" \
  -d '{
    "message": "Write a poem about AI",
    "model": "claude-3-haiku-20240307",
    "max_tokens": 200,
    "temperature": 0.8,
    "custom_endpoint": "https://api.anthropic.com/v1",
    "custom_api_key": "sk-ant-your-anthropic-key"
  }'
```

### 3. Local Ollama
```bash
curl -X POST http://localhost:3000/v1/examples/chat/completion \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-jwt-token" \
  -d '{
    "message": "Explain machine learning",
    "model": "llama2",
    "max_tokens": 400,
    "temperature": 0.6,
    "custom_endpoint": "http://localhost:11434/v1",
    "custom_api_key": ""
  }'
```

### 4. Azure OpenAI
```bash
curl -X POST http://localhost:3000/v1/examples/chat/completion \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-jwt-token" \
  -d '{
    "message": "Summarize this text",
    "model": "gpt-35-turbo",
    "max_tokens": 150,
    "temperature": 0.5,
    "custom_endpoint": "https://your-resource.openai.azure.com/openai/deployments/your-deployment",
    "custom_api_key": "your-azure-openai-key"
  }'
```

### 5. Default (Environment Config)
```bash
curl -X POST http://localhost:3000/v1/examples/chat/completion \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-jwt-token" \
  -d '{
    "message": "Hello, world!",
    "model": "gpt-3.5-turbo",
    "max_tokens": 100
  }'
```

## 🎨 JavaScript Examples

### Regular Chat with Custom Endpoint
```javascript
async function chatWithCustomEndpoint(message, endpoint, apiKey) {
  const response = await fetch('/v1/examples/chat/completion', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': 'Bearer ' + localStorage.getItem('token')
    },
    body: JSON.stringify({
      message: message,
      model: 'anthropic/claude-3-haiku',
      max_tokens: 300,
      temperature: 0.7,
      custom_endpoint: endpoint,
      custom_api_key: apiKey
    })
  });

  const data = await response.json();
  
  if (data.status === 'success') {
    console.log('AI Response:', data.data.response);
    console.log('Model used:', data.data.model);
    console.log('Processing time:', data.data.processing_time);
  } else {
    console.error('Error:', data.message);
  }
}

// Usage
chatWithCustomEndpoint(
  'Explain blockchain technology',
  'https://openrouter.ai/api/v1',
  'sk-or-v1-your-key'
);
```

### Streaming Chat with Custom Endpoint
```javascript
async function streamChatWithCustomEndpoint(message, endpoint, apiKey) {
  const response = await fetch('/v1/examples/chat/completion/stream', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': 'Bearer ' + localStorage.getItem('token')
    },
    body: JSON.stringify({
      message: message,
      model: 'gpt-4',
      max_tokens: 500,
      temperature: 0.8,
      custom_endpoint: endpoint,
      custom_api_key: apiKey
    })
  });

  const reader = response.body.getReader();
  const decoder = new TextDecoder();

  while (true) {
    const { done, value } = await reader.read();
    
    if (done) break;
    
    const chunk = decoder.decode(value);
    const lines = chunk.split('\n');
    
    for (const line of lines) {
      if (line.startsWith('data: ')) {
        const data = line.slice(6);
        
        if (data === '[DONE]') {
          console.log('\nStream completed');
          return;
        }
        
        try {
          const parsed = JSON.parse(data);
          
          if (parsed.type === 'content') {
            process.stdout.write(parsed.content);
          } else if (parsed.type === 'error') {
            console.error('\nStream error:', parsed.message);
          }
        } catch (e) {
          // Skip invalid JSON
        }
      }
    }
  }
}

// Usage
streamChatWithCustomEndpoint(
  'Tell me a story about AI',
  'https://api.openai.com/v1',
  'sk-your-openai-key'
);
```

## 🔧 Provider-Specific Configurations

### OpenRouter.ai
```json
{
  "custom_endpoint": "https://openrouter.ai/api/v1",
  "custom_api_key": "sk-or-v1-your-key",
  "model": "anthropic/claude-3-haiku"
}
```

**Available Models:**
- `anthropic/claude-3-haiku`
- `anthropic/claude-3-sonnet`
- `openai/gpt-4`
- `openai/gpt-3.5-turbo`
- `meta-llama/llama-2-70b-chat`

### Anthropic Direct
```json
{
  "custom_endpoint": "https://api.anthropic.com/v1",
  "custom_api_key": "sk-ant-your-key",
  "model": "claude-3-haiku-20240307"
}
```

### Azure OpenAI
```json
{
  "custom_endpoint": "https://your-resource.openai.azure.com/openai/deployments/your-deployment",
  "custom_api_key": "your-azure-key",
  "model": "gpt-35-turbo"
}
```

### Local Ollama
```json
{
  "custom_endpoint": "http://localhost:11434/v1",
  "custom_api_key": "",
  "model": "llama2"
}
```

## 🛡️ Security Considerations

### 1. API Key Protection
- Jangan hardcode API keys dalam frontend code
- Gunakan environment variables atau secure storage
- Implementasikan rate limiting per API key

### 2. Endpoint Validation
- Validasi URL endpoint untuk mencegah SSRF attacks
- Whitelist allowed endpoints jika diperlukan
- Monitor usage untuk detect abuse

### 3. Cost Management
- Track usage per custom endpoint
- Implement spending limits
- Monitor token consumption

## ⚡ Performance Tips

### 1. Connection Pooling
Custom clients dibuat per request. Untuk production, pertimbangkan:
- Client caching berdasarkan endpoint+key
- Connection pooling
- Request batching

### 2. Timeout Management
- Default timeout: 30s untuk regular, 60s untuk streaming
- Sesuaikan timeout berdasarkan provider
- Implement retry logic untuk network errors

### 3. Model Selection
- Pilih model yang sesuai dengan use case
- Pertimbangkan cost vs quality trade-off
- Test different models untuk optimal results

## 🔍 Error Handling

### Common Errors

#### 1. Missing Endpoint or API Key
```json
{
  "status": "error",
  "message": "Both custom_endpoint and custom_api_key must be provided together"
}
```

#### 2. Invalid API Key
```json
{
  "status": "error",
  "message": "AI API Error: Incorrect API key provided",
  "type": "invalid_request_error"
}
```

#### 3. Invalid Endpoint
```json
{
  "status": "error",
  "message": "Failed to process AI request: connection refused"
}
```

#### 4. Model Not Available
```json
{
  "status": "error",
  "message": "AI API Error: The model 'invalid-model' does not exist",
  "type": "invalid_request_error"
}
```

## 📊 Response Format

Response format tetap sama, dengan informasi tambahan tentang model yang digunakan:

```json
{
  "status": "success",
  "data": {
    "response": "AI response text",
    "model": "anthropic/claude-3-haiku",
    "usage": {
      "prompt_tokens": 45,
      "completion_tokens": 120,
      "total_tokens": 165
    },
    "processing_time": "2.345s"
  }
}
```

## 🧪 Testing

### Test dengan Different Providers
```bash
# Test OpenRouter
curl -X POST http://localhost:3000/v1/examples/chat/completion \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer test-token" \
  -d '{
    "message": "Test message",
    "custom_endpoint": "https://openrouter.ai/api/v1",
    "custom_api_key": "test-key"
  }'

# Test Default (Environment)
curl -X POST http://localhost:3000/v1/examples/chat/completion \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer test-token" \
  -d '{
    "message": "Test message"
  }'
```

## 📝 Best Practices

### 1. Fallback Strategy
```javascript
async function chatWithFallback(message) {
  const providers = [
    {
      endpoint: 'https://openrouter.ai/api/v1',
      key: 'sk-or-v1-key',
      model: 'anthropic/claude-3-haiku'
    },
    {
      endpoint: 'https://api.openai.com/v1',
      key: 'sk-openai-key',
      model: 'gpt-3.5-turbo'
    }
  ];

  for (const provider of providers) {
    try {
      const response = await chatCompletion(message, provider);
      return response;
    } catch (error) {
      console.log(`Provider ${provider.endpoint} failed, trying next...`);
    }
  }
  
  throw new Error('All providers failed');
}
```

### 2. Cost Optimization
```javascript
function selectOptimalProvider(messageLength, complexity) {
  if (messageLength < 100 && complexity === 'low') {
    return {
      endpoint: 'http://localhost:11434/v1',
      key: '',
      model: 'llama2'
    };
  } else if (complexity === 'high') {
    return {
      endpoint: 'https://api.openai.com/v1',
      key: 'sk-openai-key',
      model: 'gpt-4'
    };
  } else {
    return {
      endpoint: 'https://openrouter.ai/api/v1',
      key: 'sk-or-key',
      model: 'anthropic/claude-3-haiku'
    };
  }
}
```

### 3. Usage Tracking
```javascript
function trackUsage(provider, tokens, cost) {
  // Log usage for billing and monitoring
  console.log({
    provider: provider.endpoint,
    model: provider.model,
    tokens: tokens,
    estimated_cost: cost,
    timestamp: new Date().toISOString()
  });
}
```

## 🎉 Kesimpulan

Fitur custom endpoint dan API key memberikan fleksibilitas untuk:
- ✅ Menggunakan multiple AI providers dalam satu aplikasi
- ✅ Beralih provider berdasarkan kebutuhan per request
- ✅ Mengoptimalkan cost dengan memilih provider yang tepat
- ✅ Implementasi fallback strategy untuk reliability
- ✅ Testing dengan different models dan providers

**Status: READY TO USE** 🚀