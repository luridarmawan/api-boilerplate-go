# Real-time Swagger Filtering

This guide explains how to use the `FilterSwagger` function for real-time filtering of Swagger documentation in your API endpoints.

## Overview

The `FilterSwagger` function allows you to dynamically filter Swagger JSON documentation by removing endpoints with specific tags. This is useful for:

- Creating public API documentation by hiding internal endpoints
- Providing different documentation for different user roles
- Customizing API documentation based on client requirements
- Real-time filtering without pre-generating multiple files

## Function Signatures

### Basic Usage

```go
func FilterSwagger(data []byte, tagsToRemove string, verbose ...bool) ([]byte, error)
```

### With Options

```go
func FilterSwaggerWithOptions(data []byte, options FilterSwaggerOptions) ([]byte, error)
```

### Pretty Formatted

```go
func FilterSwaggerPretty(data []byte, tagsToRemove string, verbose ...bool) ([]byte, error)
```

### Statistics

```go
func GetSwaggerFilterStats(data []byte, tagsToRemove string) (map[string]int, error)
```

## Implementation Examples

### 1. Basic Real-time Filtering

```go
package main

import (
    "apiserver/internal/utils"
    "os"
    "path/filepath"
    "github.com/gofiber/fiber/v2"
)

func handlePublicSwagger(c *fiber.Ctx) error {
    // Read swagger file
    baseDir, _ := filepath.Abs(".")
    jsonPath := filepath.Join(baseDir, "docs", "swagger.json")
    data, err := os.ReadFile(jsonPath)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "error": "Failed to read swagger file",
        })
    }

    // Filter out internal endpoints
    filteredData, err := utils.FilterSwagger(data, "Access,Example,Permission", false)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "error": "Failed to filter swagger",
        })
    }

    c.Set("Content-Type", "application/json")
    return c.Send(filteredData)
}
```

### 2. Dynamic Filtering with Query Parameters

```go
func handleDynamicSwagger(c *fiber.Ctx) error {
    // Read swagger file
    baseDir, _ := filepath.Abs(".")
    jsonPath := filepath.Join(baseDir, "docs", "swagger.json")
    data, err := os.ReadFile(jsonPath)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "error": "Failed to read swagger file",
        })
    }

    // Get filtering options from query parameters
    tagsToRemove := c.Query("remove", "Access,Permission")
    verbose := c.Query("verbose", "false") == "true"
    pretty := c.Query("pretty", "false") == "true"

    var filteredData []byte
    if pretty {
        filteredData, err = utils.FilterSwaggerPretty(data, tagsToRemove, verbose)
    } else {
        filteredData, err = utils.FilterSwagger(data, tagsToRemove, verbose)
    }

    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "error": "Failed to filter swagger",
        })
    }

    c.Set("Content-Type", "application/json")
    return c.Send(filteredData)
}
```

### 3. Role-based Documentation

```go
func handleRoleBasedSwagger(c *fiber.Ctx) error {
    // Get user role from context or JWT token
    userRole := getUserRole(c) // Your implementation

    // Read swagger file
    baseDir, _ := filepath.Abs(".")
    jsonPath := filepath.Join(baseDir, "docs", "swagger.json")
    data, err := os.ReadFile(jsonPath)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "error": "Failed to read swagger file",
        })
    }

    // Define tags to remove based on user role
    var tagsToRemove string
    switch userRole {
    case "public":
        tagsToRemove = "Access,Permission,Admin,Internal"
    case "partner":
        tagsToRemove = "Admin,Internal"
    case "admin":
        tagsToRemove = "" // Show all endpoints
    default:
        tagsToRemove = "Access,Permission,Admin,Internal,Example"
    }

    // Filter if needed
    if tagsToRemove != "" {
        filteredData, err := utils.FilterSwagger(data, tagsToRemove, false)
        if err != nil {
            return c.Status(500).JSON(fiber.Map{
                "error": "Failed to filter swagger",
            })
        }
        c.Set("Content-Type", "application/json")
        return c.Send(filteredData)
    }

    // Return unfiltered for admin
    c.Set("Content-Type", "application/json")
    return c.Send(data)
}
```

### 4. Using FilterSwaggerOptions

```go
func handleAdvancedSwagger(c *fiber.Ctx) error {
    // Read swagger file
    baseDir, _ := filepath.Abs(".")
    jsonPath := filepath.Join(baseDir, "docs", "swagger.json")
    data, err := os.ReadFile(jsonPath)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "error": "Failed to read swagger file",
        })
    }

    // Configure filtering options
    options := utils.FilterSwaggerOptions{
        RemoveTags: []string{"Access", "Permission", "Internal"},
        Verbose:    c.Query("debug") == "true",
    }

    filteredData, err := utils.FilterSwaggerWithOptions(data, options)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "error": "Failed to filter swagger",
        })
    }

    c.Set("Content-Type", "application/json")
    return c.Send(filteredData)
}
```

### 5. Swagger Statistics Endpoint

```go
func handleSwaggerStats(c *fiber.Ctx) error {
    // Read swagger file
    baseDir, _ := filepath.Abs(".")
    jsonPath := filepath.Join(baseDir, "docs", "swagger.json")
    data, err := os.ReadFile(jsonPath)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "error": "Failed to read swagger file",
        })
    }

    // Get filtering tags from query
    tagsToRemove := c.Query("remove", "Access,Permission")

    // Get statistics
    stats, err := utils.GetSwaggerFilterStats(data, tagsToRemove)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "error": "Failed to get swagger stats",
        })
    }

    return c.JSON(fiber.Map{
        "status": "success",
        "data": stats,
        "filter": tagsToRemove,
    })
}
```

## Route Registration

Add these endpoints to your Fiber app:

```go
func main() {
    app := fiber.New()

    // Public swagger (filtered)
    app.Get("/swagger-public/doc.json", handlePublicSwagger)

    // Dynamic swagger with query parameters
    app.Get("/swagger-filtered/doc.json", handleDynamicSwagger)

    // Role-based swagger
    app.Get("/swagger-role/doc.json", authMiddleware, handleRoleBasedSwagger)

    // Swagger statistics
    app.Get("/swagger-stats", handleSwaggerStats)

    // Advanced swagger with options
    app.Get("/swagger-advanced/doc.json", handleAdvancedSwagger)

    app.Listen(":3000")
}
```

## API Endpoints

### Available Endpoints

| Endpoint | Description | Query Parameters |
|----------|-------------|------------------|
| `/swagger-public/doc.json` | Public API docs (removes Access, Example, Permission) | None |
| `/swagger-filtered/doc.json` | Dynamic filtering | `remove`, `verbose`, `pretty` |
| `/swagger-stats` | Filtering statistics | `remove` |
| `/swagger-role/doc.json` | Role-based filtering | None (uses auth context) |

### Query Parameters

#### `/swagger-filtered/doc.json`

- `remove` - Comma-separated tags to remove (default: "Access,Permission")
- `verbose` - Enable verbose logging (true/false, default: false)
- `pretty` - Pretty format JSON (true/false, default: false)

**Examples:**
```bash
# Remove specific tags
GET /swagger-filtered/doc.json?remove=Access,Permission,Internal

# Pretty formatted output
GET /swagger-filtered/doc.json?pretty=true

# Verbose logging
GET /swagger-filtered/doc.json?verbose=true&remove=Access

# Combined options
GET /swagger-filtered/doc.json?remove=Access,Permission&pretty=true&verbose=true
```

#### `/swagger-stats`

- `remove` - Tags to analyze for removal (default: "Access,Permission")

**Example:**
```bash
GET /swagger-stats?remove=Access,Permission,Internal
```

**Response:**
```json
{
  "status": "success",
  "data": {
    "total_paths": 10,
    "total_endpoints": 25,
    "removed_endpoints": 8,
    "kept_endpoints": 17
  },
  "filter": "Access,Permission,Internal"
}
```

## Performance Considerations

### Caching

For production use, consider caching filtered results:

```go
import (
    "sync"
    "time"
)

type SwaggerCache struct {
    data      []byte
    timestamp time.Time
    mutex     sync.RWMutex
}

var swaggerCache = &SwaggerCache{}

func handleCachedSwagger(c *fiber.Ctx) error {
    tagsToRemove := c.Query("remove", "Access,Permission")
    cacheKey := tagsToRemove // In production, use a proper cache key

    swaggerCache.mutex.RLock()
    // Check if cache is valid (e.g., less than 5 minutes old)
    if time.Since(swaggerCache.timestamp) < 5*time.Minute {
        data := swaggerCache.data
        swaggerCache.mutex.RUnlock()
        
        if data != nil {
            c.Set("Content-Type", "application/json")
            return c.Send(data)
        }
    }
    swaggerCache.mutex.RUnlock()

    // Cache miss or expired, generate new filtered data
    baseDir, _ := filepath.Abs(".")
    jsonPath := filepath.Join(baseDir, "docs", "swagger.json")
    data, err := os.ReadFile(jsonPath)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Failed to read swagger file"})
    }

    filteredData, err := utils.FilterSwagger(data, tagsToRemove, false)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Failed to filter swagger"})
    }

    // Update cache
    swaggerCache.mutex.Lock()
    swaggerCache.data = filteredData
    swaggerCache.timestamp = time.Now()
    swaggerCache.mutex.Unlock()

    c.Set("Content-Type", "application/json")
    return c.Send(filteredData)
}
```

### Memory Usage

The filtering functions are memory-efficient and process JSON in-memory without creating temporary files.

### Benchmarks

Run benchmarks to measure performance:

```bash
go test -bench=BenchmarkFilterSwagger ./internal/utils/
```

## Error Handling

### Common Errors

1. **File not found**
   ```go
   if err != nil {
       return c.Status(404).JSON(fiber.Map{
           "error": "Swagger file not found",
           "details": err.Error(),
       })
   }
   ```

2. **Invalid JSON**
   ```go
   if err != nil {
       return c.Status(400).JSON(fiber.Map{
           "error": "Invalid swagger JSON format",
           "details": err.Error(),
       })
   }
   ```

3. **Filtering errors**
   ```go
   if err != nil {
       return c.Status(500).JSON(fiber.Map{
           "error": "Failed to filter swagger documentation",
           "details": err.Error(),
       })
   }
   ```

### Validation

Validate swagger JSON before filtering:

```go
func handleValidatedSwagger(c *fiber.Ctx) error {
    data, err := os.ReadFile("docs/swagger.json")
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Failed to read swagger file"})
    }

    // Validate swagger JSON
    if err := utils.ValidateSwaggerJSON(data); err != nil {
        return c.Status(400).JSON(fiber.Map{
            "error": "Invalid swagger JSON",
            "details": err.Error(),
        })
    }

    // Filter and return
    filteredData, err := utils.FilterSwagger(data, "Access,Permission", false)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Failed to filter swagger"})
    }

    c.Set("Content-Type", "application/json")
    return c.Send(filteredData)
}
```

## Testing

### Unit Tests

Test your filtering endpoints:

```go
func TestSwaggerFiltering(t *testing.T) {
    app := fiber.New()
    app.Get("/swagger-filtered/doc.json", handleDynamicSwagger)

    req := httptest.NewRequest("GET", "/swagger-filtered/doc.json?remove=Access", nil)
    resp, err := app.Test(req)
    
    assert.NoError(t, err)
    assert.Equal(t, 200, resp.StatusCode)
    assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
}
```

### Integration Tests

Test with real swagger files:

```go
func TestRealSwaggerFiltering(t *testing.T) {
    // Read actual swagger file
    data, err := os.ReadFile("../../docs/swagger.json")
    require.NoError(t, err)

    // Test filtering
    filtered, err := utils.FilterSwagger(data, "Access,Permission", false)
    require.NoError(t, err)
    require.NotEmpty(t, filtered)

    // Validate result is valid JSON
    var result map[string]interface{}
    err = json.Unmarshal(filtered, &result)
    require.NoError(t, err)
}
```

## Best Practices

1. **Error Handling**: Always handle errors gracefully
2. **Caching**: Cache filtered results for better performance
3. **Validation**: Validate input swagger JSON
4. **Logging**: Log filtering operations for debugging
5. **Security**: Don't expose sensitive endpoints in public docs
6. **Testing**: Test all filtering scenarios
7. **Documentation**: Document available filtering options

## Security Considerations

1. **Sensitive Data**: Never expose internal/admin endpoints in public documentation
2. **Input Validation**: Validate query parameters to prevent injection
3. **Rate Limiting**: Apply rate limiting to filtering endpoints
4. **Authentication**: Require authentication for sensitive documentation
5. **Audit Logging**: Log access to filtered documentation

---

For more information, see:
- [Swagger Filtering Guide](./SWAGGER-FILTERING.md)
- [Utils Package Documentation](../internal/utils/)
- [API Documentation](../README.md)