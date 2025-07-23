# Module Generator Tool

This tool automatically generates a new module with all the required files following the project's standard structure.

## Usage

### Windows

```bash
tools\generate-module.bat <module-name>
```

Example:
```bash
tools\generate-module.bat product
```

### Linux/Mac

```bash
./tools/generate-module.sh <module-name>
```

Example:
```bash
./tools/generate-module.sh product
```

### Direct Go Command

```bash
go run tools/module-generator/main.go <module-name>
```

Example:
```bash
go run tools/module-generator/main.go product
```

## What It Does

The tool will:

1. Create a new folder at `internal/modules/<module-name>/`
2. Generate the following files:
   - `<module-name>_model.go`
   - `<module-name>_repository.go`
   - `<module-name>_handler.go`
   - `<module-name>_route.go`
3. Provide instructions for updating `main.go` to register the new module

## Generated Files

### Model File

Contains:
- The main data structure with UUIDv7 ID
- Request structure for creating/updating
- Table name definition
- BeforeCreate hook for UUIDv7 generation

### Repository File

Contains:
- Repository interface with CRUD operations
- Implementation of all repository methods
- Database interactions

### Handler File

Contains:
- Handler struct and constructor
- API endpoints with Swagger documentation
- CRUD operations (Create, Read, Update, Delete)
- Soft delete and restore functionality

### Route File

Contains:
- Route registration function
- All API endpoints with middleware
- Permission checks for each operation

## After Generation

After generating the module, you need to:

1. Add the module to the AutoMigrate function in `main.go`:
   ```go
   err := db.AutoMigrate(&access.User{}, &example.Example{}, &<module-name>.<ModuleName>{}, ...)
   ```

2. Initialize and register the module in `main.go`:
   ```go
   // Initialize module
   <module-name>Repo := <module-name>.NewRepository(db)
   <module-name>Handler := <module-name>.NewHandler(<module-name>Repo)
   <module-name>.Register<ModuleName>Routes(app, <module-name>Handler, authMiddleware, rateLimitMiddleware, middleware.RequirePermission)
   ```

## Customization

If you need to customize the generated module:

1. Add additional fields to the model structure
2. Add new repository methods for specific queries
3. Create new handler methods for additional endpoints
4. Register new routes in the route file

## Example

```bash
go run tools/module-generator/main.go product
```

This will generate:
- `internal/modules/product/product_model.go`
- `internal/modules/product/product_repository.go`
- `internal/modules/product/product_handler.go`
- `internal/modules/product/product_route.go`

And provide instructions for updating `main.go`.