# Configuration Module Permissions

This guide explains how to set up secure permissions for the Configuration module, ensuring only Admin users can access configuration endpoints.

## Overview

The Configuration module contains sensitive system settings that should only be accessible to administrators. This document provides scripts and procedures to:

- Add configuration permissions only to Admin group
- Remove configuration permissions from other groups
- Verify security settings
- Maintain proper access control

## Security Model

### Access Levels

| Group | Configuration Access | Reason |
|-------|---------------------|---------|
| **Admin** | ✅ Full Access | System administrators need to manage configurations |
| **Editor** | ❌ No Access | Content editors don't need system configuration access |
| **Viewer** | ❌ No Access | Read-only users should not see system configurations |
| **General Client** | ❌ No Access | External clients should not access internal configurations |

### Permissions Created

The following permissions are created for the Configuration module:

| Permission | Resource | Action | Description |
|------------|----------|--------|-------------|
| Create Configurations | configurations | create | Create new configuration entries |
| Read Configurations | configurations | read | View configuration settings |
| Update Configurations | configurations | update | Modify existing configurations |
| Delete Configurations | configurations | delete | Remove configuration entries |
| Manage Configurations | configurations | manage | Full configuration management |

## Setup Scripts

### 1. Safe Admin-Only Setup (Recommended)

Use this script to add configuration permissions only to Admin group:

#### Linux/macOS
```bash
# Run the safe setup script
./scripts/seed-configuration-permissions.sh
```

#### Windows
```cmd
REM Run the safe setup script
scripts\seed-configuration-permissions.bat
```

#### Manual Execution
```bash
# Add admin-only permissions
go run scripts/add-admin-only-configuration-permissions.go

# Optional: Clean up permissions from other groups
go run scripts/remove-configuration-permissions-from-non-admin.go
```

### 2. Direct Database Seeding

If you're setting up a fresh database, the main seeder already includes admin-only configuration permissions:

```bash
# Run main seeder (includes configuration permissions for Admin only)
go run cmd/api/main.go --seed
```

## Verification

### Check Current Permissions

Run this script to verify current permission assignments:

```bash
go run scripts/add-admin-only-configuration-permissions.go
```

The script will show:
- ✅ Which permissions exist
- ✅ Which groups have configuration access
- ⚠️ Any security issues found

### Expected Output

```
🔍 Found Admin group with ID: 1
✅ Assigned create permission to Admin group
✅ Assigned read permission to Admin group
✅ Assigned update permission to Admin group
✅ Assigned delete permission to Admin group
✅ Assigned manage permission to Admin group

🔍 Checking other groups for configuration permissions...
✅ Editor group has no configuration permissions (correct)
✅ Viewer group has no configuration permissions (correct)
✅ General client group has no configuration permissions (correct)

📊 Summary:
✅ Configuration permissions created: 5
✅ Permissions assigned to Admin group: 5

🔒 Security Status:
✅ Configuration module is now ADMIN-ONLY
✅ Only Admin group can access configuration endpoints
✅ Other groups (Editor, Viewer, General client) have NO access
```

## Cleanup Scripts

### Remove Permissions from Non-Admin Groups

If configuration permissions were accidentally given to other groups:

```bash
# Remove configuration permissions from non-admin groups
go run scripts/remove-configuration-permissions-from-non-admin.go
```

This script will:
- ✅ Find all non-admin groups
- ✅ Remove any configuration permissions from them
- ✅ Verify Admin group still has permissions
- ✅ Provide security audit report

## API Endpoints Security

### Protected Endpoints

With proper permissions in place, these endpoints will be Admin-only:

```
POST   /v1/configurations          # Create configuration
GET    /v1/configurations          # List configurations  
GET    /v1/configurations/{id}     # Get specific configuration
PUT    /v1/configurations/{id}     # Update configuration
DELETE /v1/configurations/{id}     # Delete configuration
```

### Permission Middleware

Ensure your configuration routes use permission middleware:

```go
// Example route registration
configuration.RegisterConfigurationRoutes(
    app, 
    configHandler, 
    authMiddleware, 
    rateLimitMiddleware, 
    middleware.RequirePermission
)
```

### Route Implementation Example

```go
func RegisterConfigurationRoutes(
    app *fiber.App,
    handler *Handler,
    authMiddleware fiber.Handler,
    rateLimitMiddleware fiber.Handler,
    requirePermission func(string) fiber.Handler,
) {
    api := app.Group("/v1/configurations")
    
    // Apply middleware
    api.Use(authMiddleware)
    api.Use(rateLimitMiddleware)
    
    // Admin-only endpoints
    api.Post("/", requirePermission("configurations:create"), handler.Create)
    api.Get("/", requirePermission("configurations:read"), handler.GetAll)
    api.Get("/:id", requirePermission("configurations:read"), handler.GetByID)
    api.Put("/:id", requirePermission("configurations:update"), handler.Update)
    api.Delete("/:id", requirePermission("configurations:delete"), handler.Delete)
}
```

## Testing Access Control

### Test Admin Access

```bash
# Test with admin API key
curl -X GET "http://localhost:3000/v1/configurations" \
  -H "Authorization: Bearer admin-api-key-789"

# Expected: 200 OK with configuration data
```

### Test Non-Admin Access

```bash
# Test with editor API key
curl -X GET "http://localhost:3000/v1/configurations" \
  -H "Authorization: Bearer test-api-key-123"

# Expected: 403 Forbidden
```

```bash
# Test with viewer API key  
curl -X GET "http://localhost:3000/v1/configurations" \
  -H "Authorization: Bearer viewer-api-key-456"

# Expected: 403 Forbidden
```

## Troubleshooting

### Common Issues

1. **All users can access configurations**
   ```bash
   # Check if permissions are properly assigned
   go run scripts/add-admin-only-configuration-permissions.go
   
   # Remove permissions from non-admin groups
   go run scripts/remove-configuration-permissions-from-non-admin.go
   ```

2. **Admin cannot access configurations**
   ```bash
   # Verify admin has permissions
   go run scripts/add-admin-only-configuration-permissions.go
   
   # Check admin API key is valid and not expired
   ```

3. **Permission middleware not working**
   ```go
   // Ensure routes use permission middleware
   api.Get("/", requirePermission("configurations:read"), handler.GetAll)
   ```

### Database Queries for Manual Check

```sql
-- Check configuration permissions
SELECT p.name, p.resource, p.action 
FROM permissions p 
WHERE p.resource = 'configurations';

-- Check which groups have configuration permissions
SELECT g.name as group_name, p.name as permission_name, p.action
FROM groups g
JOIN group_permissions gp ON g.id = gp.group_id
JOIN permissions p ON gp.permission_id = p.id
WHERE p.resource = 'configurations';

-- Should only show Admin group
```

## Security Best Practices

1. **Regular Audits**: Periodically check permission assignments
2. **Principle of Least Privilege**: Only Admin should have configuration access
3. **Monitoring**: Log configuration access attempts
4. **Testing**: Regularly test access control with different user roles
5. **Documentation**: Keep permission changes documented

## Migration from Existing Setup

If you previously ran the unsafe script (`scripts/add-configuration-permissions.go`):

1. **Audit Current State**:
   ```bash
   go run scripts/add-admin-only-configuration-permissions.go
   ```

2. **Remove Unsafe Permissions**:
   ```bash
   go run scripts/remove-configuration-permissions-from-non-admin.go
   ```

3. **Verify Security**:
   ```bash
   # Test that only admin can access
   curl -X GET "http://localhost:3000/v1/configurations" \
     -H "Authorization: Bearer admin-api-key-789"
   
   # Test that others cannot access
   curl -X GET "http://localhost:3000/v1/configurations" \
     -H "Authorization: Bearer test-api-key-123"
   ```

## Automation

### CI/CD Integration

Add permission verification to your CI/CD pipeline:

```yaml
# GitHub Actions example
- name: Verify Configuration Security
  run: |
    go run scripts/add-admin-only-configuration-permissions.go
    # Add tests to verify only admin can access
```

### Monitoring Script

Create a monitoring script to check permissions regularly:

```bash
#!/bin/bash
# monitor-config-permissions.sh

echo "🔍 Configuration Permissions Audit - $(date)"
go run scripts/add-admin-only-configuration-permissions.go

# Add to cron for regular checks
# 0 9 * * * /path/to/monitor-config-permissions.sh >> /var/log/config-audit.log
```

---

## Summary

- ✅ Use `scripts/seed-configuration-permissions.sh` for safe setup
- ✅ Configuration access is Admin-only by design
- ✅ Regular auditing ensures security compliance
- ✅ Proper middleware enforcement prevents unauthorized access
- ✅ Testing verifies access control works correctly

For questions or issues, refer to the main [API Documentation](../README.md) or [Security Guide](./SECURITY.md).