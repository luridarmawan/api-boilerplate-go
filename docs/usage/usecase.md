
# USE CASE DOCUMENTATION

### 1. Setup Test Data

[Lihat bagian seeder](README.md#%EF%B8%8F-seeder--test-data)


### 2. Test Endpoints dengan cURL

**Get Profile (Valid Token):**
```bash
curl -X GET "http://localhost:3000/v1/profile" \
  -H "Authorization: Bearer test-api-key-123"
```

**Get Examples (Valid Token):**
```bash
curl -X GET "http://localhost:3000/v1/examples" \
  -H "Authorization: Bearer test-api-key-123"
```

**Create Example (Valid Token):**
```bash
curl -X POST "http://localhost:3000/v1/examples" \
  -H "Authorization: Bearer test-api-key-123" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Example",
    "description": "This is a test example"
  }'
```

**Test Invalid Token:**
```bash
curl -X GET "http://localhost:3000/v1/profile" \
  -H "Authorization: Bearer invalid-token"
```

**Test Missing Token:**
```bash
curl -X GET "http://localhost:3000/v1/profile"
```

### 3. Test Permission System dengan Different User Roles

**Admin User - Can Access Everything:**
```bash
# Get all permissions (Admin only)
curl -X GET "http://localhost:3000/v1/permissions" \
  -H "Authorization: Bearer admin-api-key-789"

# Create new permission (Admin only)
curl -X POST "http://localhost:3000/v1/permissions" \
  -H "Authorization: Bearer admin-api-key-789" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Delete Examples",
    "description": "Permission to delete examples",
    "resource": "examples",
    "action": "delete"
  }'

# Get all groups (Admin only)
curl -X GET "http://localhost:3000/v1/groups" \
  -H "Authorization: Bearer admin-api-key-789"

# Create new group (Admin only)
curl -X POST "http://localhost:3000/v1/groups" \
  -H "Authorization: Bearer admin-api-key-789" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Manager",
    "description": "Manager role with limited admin access",
    "permission_ids": [1, 2, 3]
  }'
```

**Editor User - Can Create/Read Examples:**
```bash
# Can create examples
curl -X POST "http://localhost:3000/v1/examples" \
  -H "Authorization: Bearer test-api-key-123" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Editor Example",
    "description": "Created by editor user"
  }'

# Can read examples
curl -X GET "http://localhost:3000/v1/examples" \
  -H "Authorization: Bearer test-api-key-123"

# CANNOT access permissions (will get 403 Forbidden)
curl -X GET "http://localhost:3000/v1/permissions" \
  -H "Authorization: Bearer test-api-key-123"
```

**Viewer User - Read Only Access:**
```bash
# Can read examples
curl -X GET "http://localhost:3000/v1/examples" \
  -H "Authorization: Bearer test-api-key-456"

# Can view profile
curl -X GET "http://localhost:3000/v1/profile" \
  -H "Authorization: Bearer test-api-key-456"

# CANNOT create examples (will get 403 Forbidden)
curl -X POST "http://localhost:3000/v1/examples" \
  -H "Authorization: Bearer test-api-key-456" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Viewer Example",
    "description": "This should fail"
  }'
```

### 4. Test Group Permission Management

**Update Group Permissions (Admin only):**
```bash
# Add more permissions to Editor group
curl -X PUT "http://localhost:3000/v1/groups/2/permissions" \
  -H "Authorization: Bearer admin-api-key-789" \
  -H "Content-Type: application/json" \
  -d '{
    "permission_ids": [1, 2, 3, 4]
  }'

# Get specific group details
curl -X GET "http://localhost:3000/v1/groups/2" \
  -H "Authorization: Bearer admin-api-key-789"
```
