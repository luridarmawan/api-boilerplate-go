# Permission Manager CLI Tool

A command line tool to add permissions to access IDs by managing their associated group permissions in the API boilerplate system.

## Overview

This tool allows you to add permissions to users (access IDs) by adding the permissions to the group that the user belongs to. It follows the existing permission system architecture:

```
Access (User) → Group → Permissions
```

## Installation

### Build from Source

1. **Using build scripts:**
   ```bash
   # On Linux/macOS
   ./scripts/build-permission-manager.sh
   
   # On Windows
   scripts\build-permission-manager.bat
   ```

2. **Manual build:**
   ```bash
   # Build for current platform
   go build -o bin/permission-manager ./cmd/permission-manager
   
   # Build for Windows
   GOOS=windows GOARCH=amd64 go build -o bin/permission-manager.exe ./cmd/permission-manager
   ```

### Pre-built Binaries

After building, you'll find binaries in the `bin/` directory:
- `permission-manager` (current platform)
- `permission-manager-windows-amd64.exe`
- `permission-manager-linux-amd64`
- `permission-manager-darwin-amd64`
- `permission-manager-darwin-arm64`

## Usage

```bash
permission-manager [access_id] [resource] [action]
```

### Arguments

- **access_id**: UUID of the access/user to add permission to
- **resource**: Resource name (e.g., 'configurations', 'examples')
- **action**: Action name (e.g., 'create', 'read', 'update', 'delete')

### Examples

```bash
# Add 'create' permission for 'configurations' resource
./bin/permission-manager 019847a9-4efb-72c1-92fb-2c5eab3335d1 configurations create

# Add 'read' permission for 'examples' resource
./bin/permission-manager 019847a9-4efb-72c1-92fb-2c5eab3335d1 examples read

# Add 'update' permission for 'configurations' resource
./bin/permission-manager 019847a9-4efb-72c1-92fb-2c5eab3335d1 configurations update
```

### Help

```bash
./bin/permission-manager --help
./bin/permission-manager -h
```

## Expected Behavior

### Success Cases

```bash
$ ./bin/permission-manager 019847a9-4efb-72c1-92fb-2c5eab3335d1 configurations create
✓ Permission 'configurations:create' successfully added to group 'Editor' for access ID '019847a9-4efb-72c1-92fb-2c5eab3335d1'
```

### Warning Cases

```bash
$ ./bin/permission-manager 019847a9-4efb-72c1-92fb-2c5eab3335d1 configurations create
⚠ Warning: Permission 'configurations:create' already exists in group 'Editor'
```

### Error Cases

```bash
# Invalid UUID format
$ ./bin/permission-manager invalid-id configurations create
✗ Error: access_id must be a valid UUID format

# Non-existent access ID
$ ./bin/permission-manager 00000000-0000-0000-0000-000000000000 configurations create
✗ Error: access ID '00000000-0000-0000-0000-000000000000' not found

# Non-existent permission
$ ./bin/permission-manager 019847a9-4efb-72c1-92fb-2c5eab3335d1 invalid-resource create
✗ Error: permission 'invalid-resource:create' not found in database

# Access ID without group
$ ./bin/permission-manager 019847a9-4efb-72c1-92fb-2c5eab3335d1 configurations create
✗ Error: access ID '019847a9-4efb-72c1-92fb-2c5eab3335d1' is not assigned to any group
```

## Prerequisites

1. **Database Connection**: The tool uses the same database configuration as the main API server. Ensure your `.env` file is properly configured.

2. **Existing Permissions**: The tool only works with existing permissions in the database. It will not create new permissions.

3. **Group Assignment**: The access ID must be assigned to a group. Users without groups cannot have permissions added.

## Configuration

The tool uses the same configuration system as the main API server:

- Reads from `.env` file in the project root
- Uses the same database connection settings
- Supports all environment variables used by the main API

## Database Schema

The tool works with these database tables:

```sql
-- Access table (users)
access (id, name, email, api_key, group_id, ...)

-- Groups table
groups (id, name, description, ...)

-- Permissions table
permissions (id, name, resource, action, ...)

-- Many-to-many relationship
group_permissions (group_id, permission_id)
```

## Exit Codes

- `0`: Success
- `1`: Error (invalid arguments, database error, etc.)

## Troubleshooting

### Common Issues

1. **Database Connection Failed**
   - Check your `.env` file configuration
   - Ensure the database server is running
   - Verify database credentials

2. **Permission Not Found**
   - Check if the permission exists: `SELECT * FROM permissions WHERE resource = 'your-resource' AND action = 'your-action';`
   - Ensure the permission is active (`status_id = 0`)

3. **Access ID Not Found**
   - Verify the UUID format is correct
   - Check if the access ID exists: `SELECT * FROM access WHERE id = 'your-access-id';`
   - Ensure the access is active (`status_id = 0`)

4. **No Group Assigned**
   - Check if the user has a group: `SELECT * FROM access WHERE id = 'your-access-id';`
   - Assign a group to the user if needed

### Debug Mode

The tool shows GORM debug output when database operations occur, which can help with troubleshooting.

## Integration

This tool integrates seamlessly with the existing API boilerplate:

- Uses the same configuration system (`configs/config.go`)
- Leverages existing models (`access`, `permission`, `group`)
- Follows the same database connection pattern
- Respects the same status and soft-delete conventions

## Development

### Project Structure

```
cmd/permission-manager/
├── main.go              # Main CLI application
└── README.md           # This file

scripts/
├── build-permission-manager.sh   # Linux/macOS build script
└── build-permission-manager.bat  # Windows build script

test/
└── permission-manager-test.http  # Test scenarios

docs/
└── PERMISSION-MANAGER-CLI.md     # Architecture documentation
```

### Testing

See `test/permission-manager-test.http` for comprehensive test scenarios and manual testing commands.

### Contributing

1. Follow the existing code style and patterns
2. Add tests for new functionality
3. Update documentation as needed
4. Ensure cross-platform compatibility