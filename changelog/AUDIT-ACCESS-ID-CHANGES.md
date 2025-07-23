# Audit Logs: Field Rename from user_id to access_id

## Summary

Successfully renamed the `user_id` field to `access_id` in the `audit_logs` table to better reflect that it references the `id` field from the `access` table.

## Changes Made

### 1. Model Updates (`internal/modules/audit/audit_model.go`)
- ✅ Renamed `UserID` field to `AccessID` in `AuditLog` struct
- ✅ Updated field tag from `user_id` to `access_id`
- ✅ Added `AccessID` filter option to `AuditLogFilter` struct

### 2. Middleware Updates (`internal/modules/audit/audit_middleware.go`)
- ✅ Updated variable name from `userID` to `accessID`
- ✅ Updated `AuditLog` creation to use `AccessID` field

### 3. Repository Updates (`internal/modules/audit/audit_repository.go`)
- ✅ Added `access_id` filter support in `GetAuditLogs` method
- ✅ Filter allows querying audit logs by specific access UUID

### 4. Handler Updates (`internal/modules/audit/audit_handler.go`)
- ✅ Added `access_id` query parameter support
- ✅ Updated Swagger documentation to include `access_id` filter
- ✅ Updated `GetAuditLogs` handler to process `access_id` filter

### 5. Testing (`test/test.http`)
- ✅ Added example API calls with `access_id` filter
- ✅ Added multiple filter combinations examples

## Technical Details

### Before:
```go
type AuditLog struct {
    UserID *string `json:"user_id" gorm:"type:uuid;index"`
    // ... other fields
}
```

### After:
```go
type AuditLog struct {
    AccessID *string `json:"access_id" gorm:"type:uuid;index"`
    // ... other fields
}
```

## Database Schema Change

### Field Rename:
- **Old:** `user_id` (UUID, nullable, indexed)
- **New:** `access_id` (UUID, nullable, indexed)

### Purpose:
- Better semantic naming - clearly indicates reference to `access` table
- Maintains foreign key relationship to `access.id`
- Supports filtering audit logs by specific access/user

## API Changes

### New Filter Parameter:
```bash
# Filter audit logs by specific access ID
GET /v1/audit-logs?access_id=01983465-34ad-760b-87f5-b73601f6e281

# Combine with other filters
GET /v1/audit-logs?access_id=01983465-34ad-760b-87f5-b73601f6e281&method=POST&limit=10
```

### Swagger Documentation Updated:
- Added `@Param access_id query string false "Filter by access ID (UUID)"`
- Maintains backward compatibility with existing filters

## JSON Response Format

### Before:
```json
{
  "id": "01983475-5757-778b-aa8f-2143995b7fe7",
  "user_id": "01983465-34ad-760b-87f5-b73601f6e281",
  "user_email": "user@example.com",
  // ... other fields
}
```

### After:
```json
{
  "id": "01983475-5757-778b-aa8f-2143995b7fe7",
  "access_id": "01983465-34ad-760b-87f5-b73601f6e281",
  "user_email": "user@example.com",
  // ... other fields
}
```

## Benefits

1. **Clearer Semantics**: Field name clearly indicates reference to access table
2. **Better Filtering**: Can filter audit logs by specific access/user ID
3. **Consistent Naming**: Aligns with database relationship naming conventions
4. **Enhanced Querying**: Supports more granular audit log analysis

## Usage Examples

### Creating Audit Logs (Automatic via Middleware):
```go
// Middleware automatically sets AccessID from authenticated user
auditLog := &AuditLog{
    AccessID:   &user.ID,  // References access.id
    UserEmail:  user.Email,
    Method:     "GET",
    Path:       "/v1/examples",
    // ... other fields
}
```

### Querying by Access ID:
```bash
# Get all audit logs for a specific user/access
curl -X GET "http://localhost:3000/v1/audit-logs?access_id=01983465-34ad-760b-87f5-b73601f6e281" \
  -H "Authorization: Bearer admin-api-key-789"

# Combine filters for specific analysis
curl -X GET "http://localhost:3000/v1/audit-logs?access_id=01983465-34ad-760b-87f5-b73601f6e281&method=POST&date_from=2025-01-01" \
  -H "Authorization: Bearer admin-api-key-789"
```

## Verification

✅ **Build Status**: Successful compilation  
✅ **Field Rename**: `user_id` → `access_id` completed  
✅ **Filter Support**: New `access_id` filter working  
✅ **API Documentation**: Swagger updated  
✅ **Test Examples**: Updated with new field usage  

## Notes

- This change is for initial project setup only (no migration scripts needed)
- Field maintains the same data type and constraints (UUID, nullable, indexed)
- Backward compatibility maintained for other existing filters
- The change provides better semantic clarity for the audit system