# Swagger Filter Tool

A command-line tool to filter Swagger JSON documentation by removing endpoints with specific tags.

## Features

- ✅ Remove endpoints by tags (e.g., "Access", "Example", "Permission")
- ✅ Preserve JSON formatting and structure
- ✅ Support for multiple tags filtering
- ✅ Verbose logging option
- ✅ Backup and overwrite options
- ✅ Cross-platform support (Windows, Linux, macOS)

## Usage

### Quick Usage with Scripts

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
REM Remove default tags (Access, Example, Permission)
scripts\filter-swagger.bat

REM Remove specific tags
scripts\filter-swagger.bat "Access,Permission"

REM Specify input and output files
scripts\filter-swagger.bat "Access,Example" "docs\swagger.json" "docs\swagger-filtered.json"
```

### Direct Tool Usage

#### Build the tool
```bash
cd tools/swagger-filter
go build -o swagger-filter main.go
```

#### Run the tool
```bash
# Basic usage - remove default tags
./swagger-filter -input="../../docs/swagger.json"

# Remove specific tags
./swagger-filter -input="../../docs/swagger.json" -remove-tags="Access,Permission,Example"

# Save to different file
./swagger-filter -input="../../docs/swagger.json" -output="../../docs/swagger-filtered.json"

# Enable verbose logging
./swagger-filter -input="../../docs/swagger.json" -verbose
```

## Command Line Options

| Option | Description | Default |
|--------|-------------|---------|
| `-input` | Input swagger JSON file path | `docs/swagger.json` |
| `-output` | Output swagger JSON file path | Same as input (overwrite) |
| `-remove-tags` | Comma-separated tags to remove | `Access,Example,Permission` |
| `-verbose` | Enable verbose logging | `false` |

## Examples

### Example 1: Remove Access and Permission endpoints
```bash
./swagger-filter -input="docs/swagger.json" -remove-tags="Access,Permission" -verbose
```

Output:
```
Input file: docs/swagger.json
Output file: docs/swagger.json
Tags to remove: [Access Permission]
🗑️  Removing: GET /v1/profile
🗑️  Removing: PUT /v1/access/{id}/expired-date
🗑️  Removing: DELETE /v1/access/{id}/expired-date
🗑️  Removing: PUT /v1/access/{id}/rate-limit
✅ Keeping: GET /v1/examples
✅ Keeping: POST /v1/examples
✅ Keeping: GET /v1/examples/{id}
📊 Summary: 4 endpoints removed, 8 endpoints kept
✅ Successfully filtered swagger.json
📁 Output: docs/swagger.json
```

### Example 2: Create filtered copy
```bash
./swagger-filter -input="docs/swagger.json" -output="docs/swagger-public.json" -remove-tags="Access,Permission,Audit"
```

### Example 3: Remove only Example endpoints
```bash
./swagger-filter -input="docs/swagger.json" -remove-tags="Example"
```

## Integration with Build Process

### Add to build script
```bash
# In scripts/build.sh
echo "Generating Swagger documentation..."
swag init -g cmd/api/main.go -o docs

echo "Filtering Swagger documentation..."
./scripts/filter-swagger.sh "Access,Permission"
```

### Add to CI/CD Pipeline

#### GitHub Actions
```yaml
- name: Generate and Filter Swagger
  run: |
    swag init -g cmd/api/main.go -o docs
    ./scripts/filter-swagger.sh "Access,Permission"
```

#### GitLab CI/CD
```yaml
build:
  script:
    - swag init -g cmd/api/main.go -o docs
    - ./scripts/filter-swagger.sh "Access,Permission"
```

## Use Cases

1. **Public API Documentation**: Remove internal/admin endpoints from public docs
2. **Client-Specific Documentation**: Create different documentation for different client types
3. **Security**: Hide sensitive endpoints from public documentation
4. **API Versioning**: Create different versions of documentation
5. **Development vs Production**: Different endpoint visibility per environment

## JSON Structure

The tool works with standard Swagger/OpenAPI 2.0 JSON format:

```json
{
  "swagger": "2.0",
  "info": {...},
  "host": "api.example.com",
  "basePath": "/",
  "paths": {
    "/v1/examples": {
      "get": {
        "tags": ["Example"],
        "summary": "Get all examples",
        ...
      }
    }
  }
}
```

## Error Handling

The tool handles various error scenarios:

- ✅ Missing input file
- ✅ Invalid JSON format
- ✅ Missing tags in endpoints
- ✅ File permission issues
- ✅ Invalid command line arguments

## Performance

- Fast processing even for large Swagger files
- Memory efficient JSON parsing
- Preserves original JSON formatting
- No external dependencies

## Troubleshooting

### Common Issues

1. **File not found**
   ```
   Error reading file docs/swagger.json: no such file or directory
   ```
   Solution: Ensure the swagger.json file exists and path is correct

2. **Permission denied**
   ```
   Error writing file docs/swagger.json: permission denied
   ```
   Solution: Check file permissions or run with appropriate privileges

3. **Invalid JSON**
   ```
   Error parsing JSON: invalid character '}' looking for beginning of value
   ```
   Solution: Ensure the input file is valid JSON format

### Debug Mode

Use `-verbose` flag to see detailed processing information:
```bash
./swagger-filter -input="docs/swagger.json" -verbose
```

## Contributing

To extend the tool:

1. Add new filtering criteria (by path, method, etc.)
2. Add support for OpenAPI 3.0
3. Add configuration file support
4. Add batch processing for multiple files

## License

This tool is part of the API boilerplate project and follows the same MIT license.