# Swagger Documentation Filtering

This guide explains how to use the Swagger filtering tool to create customized API documentation by removing specific endpoints based on tags.

## Overview

The Swagger filtering tool allows you to:
- Remove endpoints with specific tags from Swagger documentation
- Create public API documentation by hiding internal endpoints
- Generate client-specific documentation
- Maintain multiple versions of API documentation

## Quick Start

### Generate Public Documentation
```bash
# Generate public docs (removes Access and Permission endpoints)
./scripts/generate-public-docs.sh
```

### Custom Filtering
```bash
# Remove specific tags
./scripts/filter-swagger.sh "Access,Permission,Example"

# Create filtered copy
./scripts/filter-swagger.sh "Internal,Admin" "docs/swagger.json" "docs/swagger-public.json"
```

## Tool Usage

### Basic Commands

#### Linux/macOS
```bash
# Remove default tags (Access, Example, Permission)
./scripts/filter-swagger.sh

# Remove specific tags
./scripts/filter-swagger.sh "Access,Permission"

# Specify input and output files
./scripts/filter-swagger.sh "Access,Example" "docs/swagger.json" "docs/swagger-filtered.json"
```

#### Windows
```cmd
REM Remove default tags
scripts\filter-swagger.bat

REM Remove specific tags
scripts\filter-swagger.bat "Access,Permission"

REM Specify files
scripts\filter-swagger.bat "Access,Example" "docs\swagger.json" "docs\swagger-filtered.json"
```

### Direct Tool Usage

```bash
# Build the tool
go build -o tools/swagger-filter/swagger-filter tools/swagger-filter/main.go

# Run with options
tools/swagger-filter/swagger-filter \
  -input="docs/swagger.json" \
  -output="docs/swagger-public.json" \
  -remove-tags="Access,Permission,Internal" \
  -verbose
```

## Command Line Options

| Option | Description | Default | Example |
|--------|-------------|---------|---------|
| `-input` | Input swagger JSON file | `docs/swagger.json` | `-input="api-docs.json"` |
| `-output` | Output file (empty = overwrite input) | Same as input | `-output="public-docs.json"` |
| `-remove-tags` | Comma-separated tags to remove | `Access,Example,Permission` | `-remove-tags="Admin,Internal"` |
| `-verbose` | Enable detailed logging | `false` | `-verbose` |

## Use Cases

### 1. Public API Documentation

Remove internal and admin endpoints from public documentation:

```bash
# Remove internal endpoints
./scripts/filter-swagger.sh "Access,Permission,Admin,Internal"
```

**Before filtering:**
- `/v1/profile` (Access)
- `/v1/permissions` (Permission)
- `/v1/examples` (Example)
- `/v1/audit-logs` (Audit)

**After filtering:**
- `/v1/examples` (Example) - removed
- `/v1/audit-logs` (Audit) - kept

### 2. Client-Specific Documentation

Create different documentation for different client types:

```bash
# For mobile clients (remove web-specific endpoints)
./scripts/filter-swagger.sh "WebOnly,Desktop" "docs/swagger.json" "docs/swagger-mobile.json"

# For partner APIs (remove internal endpoints)
./scripts/filter-swagger.sh "Internal,Admin" "docs/swagger.json" "docs/swagger-partner.json"
```

### 3. Environment-Specific Documentation

```bash
# Production docs (remove debug/test endpoints)
./scripts/filter-swagger.sh "Debug,Test,Development" "docs/swagger.json" "docs/swagger-prod.json"

# Development docs (keep all endpoints)
cp docs/swagger.json docs/swagger-dev.json
```

### 4. Security-Focused Filtering

```bash
# Remove sensitive endpoints
./scripts/filter-swagger.sh "Sensitive,Admin,Internal,Debug"
```

## Integration with Build Process

### Add to Build Scripts

#### In `scripts/build.sh`:
```bash
# Generate documentation
echo "Generating Swagger documentation..."
swag init -g cmd/api/main.go -o docs

# Create public version
echo "Creating public documentation..."
./scripts/filter-swagger.sh "Access,Permission" "docs/swagger.json" "docs/swagger-public.json"
```

#### In `scripts/generate-public-docs.sh`:
```bash
# Complete workflow
swag init -g cmd/api/main.go -o docs
cp docs/swagger.json docs/swagger-full.json
./scripts/filter-swagger.sh "Access,Permission"
cp docs/swagger.json docs/swagger-public.json
```

### CI/CD Integration

#### GitHub Actions
```yaml
- name: Generate and Filter Documentation
  run: |
    swag init -g cmd/api/main.go -o docs
    ./scripts/filter-swagger.sh "Access,Permission,Internal"
    
- name: Upload Public Documentation
  uses: actions/upload-artifact@v3
  with:
    name: api-documentation
    path: docs/swagger.json
```

#### GitLab CI/CD
```yaml
generate-docs:
  script:
    - swag init -g cmd/api/main.go -o docs
    - ./scripts/filter-swagger.sh "Access,Permission,Internal"
  artifacts:
    paths:
      - docs/swagger.json
    expire_in: 1 week
```

## Advanced Usage

### Multiple Output Files

Create different versions for different audiences:

```bash
# Complete documentation
swag init -g cmd/api/main.go -o docs

# Public API (remove internal endpoints)
./scripts/filter-swagger.sh "Access,Permission,Admin" "docs/swagger.json" "docs/swagger-public.json"

# Partner API (remove examples and debug)
./scripts/filter-swagger.sh "Example,Debug,Test" "docs/swagger.json" "docs/swagger-partner.json"

# Mobile API (remove web-specific endpoints)
./scripts/filter-swagger.sh "WebOnly,Desktop,Admin" "docs/swagger.json" "docs/swagger-mobile.json"
```

### Batch Processing

```bash
#!/bin/bash
# Generate multiple filtered versions

TAGS_CONFIG=(
    "Access,Permission,Admin:public"
    "Example,Debug:partner"
    "WebOnly,Desktop:mobile"
    "Internal,Sensitive:external"
)

for config in "${TAGS_CONFIG[@]}"; do
    IFS=':' read -r tags suffix <<< "$config"
    echo "Creating $suffix documentation..."
    ./scripts/filter-swagger.sh "$tags" "docs/swagger.json" "docs/swagger-$suffix.json"
done
```

### Custom Tag Strategy

Organize your endpoints with strategic tagging:

```go
// @Tags Public
// @Summary Get public data
func GetPublicData(c *fiber.Ctx) error { ... }

// @Tags Internal
// @Summary Internal admin function
func AdminFunction(c *fiber.Ctx) error { ... }

// @Tags Partner
// @Summary Partner-specific endpoint
func PartnerEndpoint(c *fiber.Ctx) error { ... }
```

Then filter accordingly:
```bash
# Public documentation
./scripts/filter-swagger.sh "Internal,Admin,Debug"

# Partner documentation  
./scripts/filter-swagger.sh "Internal,Admin"

# Internal documentation (keep all)
cp docs/swagger.json docs/swagger-internal.json
```

## File Structure

After running the filtering tools, you'll have:

```
docs/
├── swagger.json          # Filtered documentation (used by /docs endpoint)
├── swagger-full.json     # Complete documentation (backup)
├── swagger-public.json   # Public API documentation
├── swagger-partner.json  # Partner API documentation
├── swagger-mobile.json   # Mobile-specific documentation
└── swagger.yaml          # YAML version (generated by swag)
```

## Troubleshooting

### Common Issues

1. **Tool not found**
   ```bash
   # Build the tool first
   go build -o tools/swagger-filter/swagger-filter tools/swagger-filter/main.go
   ```

2. **Permission denied**
   ```bash
   # Make scripts executable
   chmod +x scripts/filter-swagger.sh
   chmod +x scripts/generate-public-docs.sh
   ```

3. **Invalid JSON**
   ```bash
   # Regenerate swagger documentation
   swag init -g cmd/api/main.go -o docs
   ```

4. **No endpoints removed**
   ```bash
   # Check if tags exist in your swagger annotations
   grep -r "@Tags" internal/modules/
   ```

### Verbose Output

Use verbose mode to see what's happening:

```bash
./scripts/filter-swagger.sh "Access,Permission" "" "" -verbose
```

Output:
```
🗑️  Removing: GET /v1/profile
🗑️  Removing: PUT /v1/access/{id}/expired-date
✅ Keeping: GET /v1/examples
✅ Keeping: POST /v1/examples
📊 Summary: 4 endpoints removed, 8 endpoints kept
```

## Best Practices

1. **Backup Original**: Always keep a backup of the complete documentation
2. **Consistent Tagging**: Use consistent tag names across your API
3. **Environment-Specific**: Create different filtered versions for different environments
4. **Automation**: Integrate filtering into your build/deployment process
5. **Documentation**: Document which endpoints are available in each filtered version
6. **Testing**: Test filtered documentation to ensure it's valid
7. **Version Control**: Consider versioning your filtered documentation

## Security Considerations

1. **Sensitive Data**: Never expose internal/admin endpoints in public documentation
2. **API Keys**: Remove endpoints that expose API key management
3. **Debug Information**: Filter out debug/development endpoints in production
4. **User Data**: Be careful with endpoints that expose user information
5. **Rate Limiting**: Document rate limits appropriately for each audience

## Performance

- Fast processing even for large Swagger files (1000+ endpoints)
- Memory efficient JSON parsing
- Preserves original JSON formatting
- No external dependencies beyond Go standard library

---

For more information, see:
- [Tool Documentation](../tools/swagger-filter/README.md)
- [API Documentation](../README.md)
- [Deployment Guide](./DEPLOYMENT.md)