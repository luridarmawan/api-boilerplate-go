

## 🛠️ Permission Manager CLI Tool

A powerful command line tool to manage permissions for access IDs by adding them to the associated group.

### Quick Start

```bash
# Build the CLI tool
./scripts/build-permission-manager.sh  # Linux/macOS
# or
scripts\build-permission-manager.bat   # Windows

# Add permission to an access ID
./bin/permission-manager 019847a9-4efb-72c1-92fb-2c5eab3335d1 configurations create
```

### Usage Template

```bash
permission-manager [access_id] [resource] [action]
```

### Examples

```bash
# Add 'create' permission for 'configurations' resource
./bin/permission-manager 019847a9-4efb-72c1-92fb-2c5eab3335d1 configurations create

# Add 'read' permission for 'examples' resource
./bin/permission-manager 019847a9-4efb-72c1-92fb-2c5eab3335d1 examples read

# Add 'update' permission for 'configurations' resource
./bin/permission-manager 019847a9-4efb-72c1-92fb-2c5eab3335d1 configurations update
```

### Key Features

- ✅ **Group-based Permission Management**: Adds permissions to the group associated with the access ID
- ✅ **Validation**: Validates access ID format, existence, and permission availability
- ✅ **Error Handling**: Clear error messages for various failure scenarios
- ✅ **Duplicate Prevention**: Warns if permission already exists in the group
- ✅ **Cross-platform**: Builds for Windows, Linux, and macOS
- ✅ **Database Integration**: Uses the same database configuration as the main API

### Expected Output

**Success:**
```bash
✓ Permission 'configurations:create' successfully added to group 'Editor' for access ID '019847a9-4efb-72c1-92fb-2c5eab3335d1'
```

**Warning (already exists):**
```bash
⚠ Warning: Permission 'configurations:create' already exists in group 'Editor'
```

**Error examples:**
```bash
✗ Error: access ID '00000000-0000-0000-0000-000000000000' not found
✗ Error: permission 'invalid-resource:create' not found in database
✗ Error: access_id must be a valid UUID format
```

For detailed documentation, see: [`cmd/permission-manager/README.md`](cmd/permission-manager/README.md)
