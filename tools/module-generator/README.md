# Module Generator Tool

This tool automatically generates a new module with all the required files following the project's standard structure. It can also optionally generate permission scripts and test files.

## Usage

### Windows

```bash
tools\generate-module.bat <module-name> [--with-permissions]
```

Examples:
```bash
tools\generate-module.bat product
tools\generate-module.bat product --with-permissions
```

### Linux/Mac

```bash
./tools/generate-module.sh <module-name> [--with-permissions]
```

Examples:
```bash
./tools/generate-module.sh product
./tools/generate-module.sh product --with-permissions
```

### Direct Go Command

```bash
go run tools/module-generator/main.go <module-name> [--with-permissions]
```

Examples:
```bash
go run tools/module-generator/main.go product
go run tools/module-generator/main.go product --with-permissions
```

## What It Does

### Basic Module Generation

The tool will:

1. Create a new folder at `internal/modules/<module-name>/`
2. Generate the following files:
   - `<module-name>_model.go`
   - `<module-name>_repository.go`
   - `<module-name>_handler.go`
   - `<module-name>_route.go`
3. Provide instructions for updating `main.go` to register the new module

### With Permissions Flag (`--with-permissions`)

When using the `--with-permissions` flag, the tool will additionally:

1. Create permission scripts:
   - `scripts/add-<module-name>-permissions.go` (Go script)
   - `scripts/add-<module-name>-permissions.sql` (SQL script)
2. Create test file:
   - `test/<module-name>-api-test.http` (HTTP test requests)
3. Generate permissions for CRUD operations (create, read, update, delete)
4. Assign permissions to appropriate groups:
   - **Admin**: All permissions (create, read, update, delete)
   - **Editor**: create, read, update permissions
   - **Viewer**: read permission only
   - **General client**: read permission only

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

### If Generated With Permissions

If you used the `--with-permissions` flag:

3. Run the permission script to add permissions to the database:
   ```bash
   go run scripts/add-<module-name>-permissions.go
   ```
   
   Or use the SQL script:
   ```bash
   psql -d your_database_name -f scripts/add-<module-name>-permissions.sql
   ```

4. Test the API endpoints using the generated test file:
   - Open `test/<module-name>-api-test.http` in VS Code with REST Client extension
   - Or use the requests as examples for your preferred HTTP client

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