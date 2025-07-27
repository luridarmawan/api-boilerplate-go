# AI Chat Completion Example

Contoh penggunaan endpoint AI Chat Completion yang telah ditambahkan ke module example.

## 📋 Endpoints yang Tersedia

### 1. Chat Completion (Regular)
```
POST /v1/examples/chat/completion
```

### 2. Chat Completion (Streaming)
```
POST /v1/examples/chat/completion/stream
```

## 🔧 Konfigurasi

Pastikan environment variables AI sudah dikonfigurasi di file `.env`:

```env
AI_BASE_URL=https://api.openai.com/v1
AI_API_KEY=your-openai-api-key-here
AI_TIMEOUT=30
```

## 📖 Penggunaan

### 1. Regular Chat Completion

#### Request
```bash
curl -X POST http://localhost:3000/v1/examples/chat/completion \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-jwt-token" \
  -d '{
    "message": "Jelaskan apa itu REST API dalam bahasa Indonesia",
    "model": "gpt-3.5-turbo",
    "max_tokens": 300,
    "temperature": 0.7,
    "system_prompt": "Kamu adalah asisten AI yang membantu menjelaskan konsep pemrograman dalam bahasa Indonesia."
  }'
```

#### Response
```json
{
  "status": "success",
  "data": {
    "response": "REST API (Representational State Transfer Application Programming Interface) adalah...",
    "model": "gpt-3.5-turbo",
    "usage": {
      "prompt_tokens": 45,
      "completion_tokens": 120,
      "total_tokens": 165
    },
    "processing_time": "2.345s"
  }
}
```

### 2. Streaming Chat Completion

#### Request
```bash
curl -X POST http://localhost:3000/v1/examples/chat/completion/stream \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-jwt-token" \
  -d '{
    "message": "Ceritakan tentang sejarah pemrograman Go",
    "model": "gpt-3.5-turbo",
    "max_tokens": 500,
    "temperature": 0.8
  }'
```

#### Response (Server-Sent Events)
```
data: {"type":"start","message":"Starting AI response..."}

data: {"type":"content","content":"Go adalah","model":"gpt-3.5-turbo"}

data: {"type":"content","content":" bahasa pemrograman","model":"gpt-3.5-turbo"}

data: {"type":"content","content":" yang dikembangkan","model":"gpt-3.5-turbo"}

...

data: {"type":"end","message":"Stream completed"}

data: [DONE]
```

## 🎯 Contoh Penggunaan dalam JavaScript

### Regular Chat Completion
```javascript
async function chatCompletion(message) {
  const response = await fetch('/v1/examples/chat/completion', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': 'Bearer ' + localStorage.getItem('token')
    },
    body: JSON.stringify({
      message: message,
      model: 'gpt-3.5-turbo',
      max_tokens: 300,
      temperature: 0.7,
      system_prompt: 'Kamu adalah asisten yang membantu.'
    })
  });

  const data = await response.json();
  
  if (data.status === 'success') {
    console.log('AI Response:', data.data.response);
    console.log('Tokens used:', data.data.usage.total_tokens);
    console.log('Processing time:', data.data.processing_time);
  } else {
    console.error('Error:', data.message);
  }
}

// Penggunaan
chatCompletion('Apa itu artificial intelligence?');
```

### Streaming Chat Completion
```javascript
async function streamingChat(message) {
  const response = await fetch('/v1/examples/chat/completion/stream', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': 'Bearer ' + localStorage.getItem('token')
    },
    body: JSON.stringify({
      message: message,
      model: 'gpt-3.5-turbo',
      max_tokens: 500,
      temperature: 0.7
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
          console.log('Stream completed');
          return;
        }
        
        try {
          const parsed = JSON.parse(data);
          
          if (parsed.type === 'content') {
            process.stdout.write(parsed.content); // Print content as it arrives
          } else if (parsed.type === 'error') {
            console.error('Stream error:', parsed.message);
          }
        } catch (e) {
          // Skip invalid JSON
        }
      }
    }
  }
}

// Penggunaan
streamingChat('Jelaskan konsep machine learning');
```

## 🎨 Contoh Frontend Integration

### HTML + JavaScript
```html
<!DOCTYPE html>
<html>
<head>
    <title>AI Chat Example</title>
</head>
<body>
    <div id="chat-container">
        <div id="messages"></div>
        <input type="text" id="message-input" placeholder="Ketik pesan Anda...">
        <button onclick="sendMessage()">Kirim</button>
        <button onclick="sendStreamingMessage()">Kirim (Streaming)</button>
    </div>

    <script>
        const messagesDiv = document.getElementById('messages');
        const messageInput = document.getElementById('message-input');

        function addMessage(sender, content) {
            const messageDiv = document.createElement('div');
            messageDiv.innerHTML = `<strong>${sender}:</strong> ${content}`;
            messagesDiv.appendChild(messageDiv);
            messagesDiv.scrollTop = messagesDiv.scrollHeight;
        }

        async function sendMessage() {
            const message = messageInput.value.trim();
            if (!message) return;

            addMessage('You', message);
            messageInput.value = '';

            try {
                const response = await fetch('/v1/examples/chat/completion', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                        'Authorization': 'Bearer ' + localStorage.getItem('token')
                    },
                    body: JSON.stringify({
                        message: message,
                        system_prompt: 'Kamu adalah asisten AI yang membantu menjawab pertanyaan dalam bahasa Indonesia.'
                    })
                });

                const data = await response.json();
                
                if (data.status === 'success') {
                    addMessage('AI', data.data.response);
                } else {
                    addMessage('Error', data.message);
                }
            } catch (error) {
                addMessage('Error', 'Failed to send message: ' + error.message);
            }
        }

        async function sendStreamingMessage() {
            const message = messageInput.value.trim();
            if (!message) return;

            addMessage('You', message);
            messageInput.value = '';

            const aiMessageDiv = document.createElement('div');
            aiMessageDiv.innerHTML = '<strong>AI:</strong> ';
            const contentSpan = document.createElement('span');
            aiMessageDiv.appendChild(contentSpan);
            messagesDiv.appendChild(aiMessageDiv);

            try {
                const response = await fetch('/v1/examples/chat/completion/stream', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                        'Authorization': 'Bearer ' + localStorage.getItem('token')
                    },
                    body: JSON.stringify({
                        message: message,
                        system_prompt: 'Kamu adalah asisten AI yang membantu menjawab pertanyaan dalam bahasa Indonesia.'
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
                                return;
                            }
                            
                            try {
                                const parsed = JSON.parse(data);
                                
                                if (parsed.type === 'content') {
                                    contentSpan.textContent += parsed.content;
                                    messagesDiv.scrollTop = messagesDiv.scrollHeight;
                                } else if (parsed.type === 'error') {
                                    contentSpan.textContent += ' [Error: ' + parsed.message + ']';
                                }
                            } catch (e) {
                                // Skip invalid JSON
                            }
                        }
                    }
                }
            } catch (error) {
                contentSpan.textContent = 'Failed to send message: ' + error.message;
            }
        }

        // Enter key to send message
        messageInput.addEventListener('keypress', function(e) {
            if (e.key === 'Enter') {
                sendMessage();
            }
        });
    </script>
</body>
</html>
```

## 🔧 Parameter Request

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `message` | string | ✅ | - | Pesan yang akan dikirim ke AI |
| `model` | string | ❌ | `gpt-3.5-turbo` | Model AI yang digunakan |
| `max_tokens` | integer | ❌ | `500` | Maksimal token dalam response |
| `temperature` | float | ❌ | `0.7` | Kreativitas response (0.0-2.0) |
| `system_prompt` | string | ❌ | Default system prompt | Instruksi sistem untuk AI |

## 📊 Response Format

### Regular Response
```json
{
  "status": "success",
  "data": {
    "response": "AI response text",
    "model": "gpt-3.5-turbo",
    "usage": {
      "prompt_tokens": 45,
      "completion_tokens": 120,
      "total_tokens": 165
    },
    "processing_time": "2.345s"
  }
}
```

### Error Response
```json
{
  "status": "error",
  "message": "Error description",
  "type": "error_type"
}
```

### Streaming Events
```json
{"type":"start","message":"Starting AI response..."}
{"type":"content","content":"text chunk","model":"gpt-3.5-turbo"}
{"type":"end","message":"Stream completed"}
{"type":"error","message":"Error description"}
```

## 🔒 Authentication & Authorization

Endpoint ini memerlukan:
1. **Authentication**: Bearer token JWT
2. **Authorization**: Permission `examples:create`
3. **Rate Limiting**: Sesuai konfigurasi rate limit

## ⚡ Performance Tips

1. **Gunakan Streaming** untuk response yang panjang
2. **Set max_tokens** sesuai kebutuhan untuk menghemat cost
3. **Cache responses** untuk pertanyaan yang sering diulang
4. **Monitor usage** untuk mengontrol biaya API

## 🐛 Error Handling

### Common Errors

1. **Invalid API Key**
```json
{
  "status": "error",
  "message": "AI API Error: Incorrect API key provided",
  "type": "invalid_request_error"
}
```

2. **Rate Limit Exceeded**
```json
{
  "status": "error",
  "message": "AI API Error: Rate limit reached",
  "type": "rate_limit_exceeded"
}
```

3. **Model Not Found**
```json
{
  "status": "error",
  "message": "AI API Error: The model 'invalid-model' does not exist",
  "type": "invalid_request_error"
}
```

## 🧪 Testing

### Unit Test Example
```bash
# Test regular completion
curl -X POST http://localhost:3000/v1/examples/chat/completion \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer test-token" \
  -d '{"message": "Hello, world!"}'

# Test streaming completion
curl -X POST http://localhost:3000/v1/examples/chat/completion/stream \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer test-token" \
  -d '{"message": "Tell me a story"}'
```

## 📝 Notes

- Endpoint ini menggunakan AI client library yang telah dibuat di `internal/ai`
- Support untuk berbagai provider AI (OpenAI, Anthropic, Ollama, dll)
- Streaming menggunakan Server-Sent Events (SSE)
- Response time tracking untuk monitoring performance
- Error handling yang comprehensive untuk berbagai skenario