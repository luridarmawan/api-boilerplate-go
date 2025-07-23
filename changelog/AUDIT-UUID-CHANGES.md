# Audit Logs Table ID Migration to UUIDv7

## Summary

Successfully migrated the `audit_logs` table's `id` field from `uint` to UUIDv7 string type for better performance and time-sortable properties.

## Changes Made

### 1. Model Updates (`internal/modules/audit/audit_model.go`)
- ✅ Changed `AuditLog.ID` from `uint` to `string` with UUID type
- ✅ Updated `AuditLogResponse.ID` from `uint` to `string`
- ✅ Added `BeforeCreate` hook to generate UUIDv7 automatically
- ✅ Added imports for `utils` and `gorm`

### 2. Repository Updates (`internal/modules/audit/audit_repository.go`)
- ✅ Updated `Repository` interface `GetAuditLogByID` method to accept `string` ID
- ✅ Updated repository implementation to work with string IDs

### 3. Handler Updates (`internal/modules/audit/audit_handler.go`)
- ✅ Updated `GetAuditLog` handler to work with string ID parameters
- ✅ Removed `strconv.ParseUint` calls for ID parsing
- ✅ Updated Swagger documentation to use `string` for ID parameters

### 4. Testing (`test/uuid_test.go`)
- ✅ Added `TestAuditLogUUIDGeneration` test function
- ✅ Verified time-sortable properties of generated UUIDs
- ✅ All tests passing

### 5. API Testing (`test/test.http`)
- ✅ Added example API calls with UUID format for audit logs

## Technical Details

### Before (Integer ID):
```go
type AuditLog struct {
    ID uint `json:"id" gorm:"primaryKey"`
    // ... other fields
}
```

### After (UUIDv7):
```go
type AuditLog struct {
    ID string `json:"id" gorm:"type:uuid;primaryKey"`
    // ... other fields
}

func (a *AuditLog) BeforeCreate(tx *gorm.DB) error {
    if a.ID == "" {
        a.ID = utils.GenerateUUIDv7()
    }
    return nil
}
```

## API Changes

### Before (Integer ID):
```
GET /v1/audit-logs/123
```

### After (UUID):
```
GET /v1/audit-logs/01983475-5757-778b-aa8f-2143995b7fe7
```

## Benefits Achieved

1. **Time-sortable**: Audit logs are naturally ordered by creation time
2. **Better performance**: Reduced database index fragmentation
3. **Globally unique**: No collision risk across distributed systems
4. **Debugging friendly**: Timestamp embedded in ID for easier troubleshooting
5. **Consistent with access table**: Both tables now use UUIDv7

## UUIDv7 Format for Audit Logs

```
01983475-5757-778b-aa8f-2143995b7fe7
├─────────────┤ ├─┤ ├─────────────────┤
│             │ │  │
│             │ │  └─ Random data (62 bits)
│             │ └─ Version (4 bits) + Random (12 bits)
└─ Timestamp (48 bits) - milliseconds since Unix epoch
```

## Verification

✅ **Build Status**: Successful compilation  
✅ **Tests**: All UUID tests passing  
✅ **Time-sortable**: Verified lexicographic ordering by timestamp  
✅ **Uniqueness**: Verified all generated UUIDs are unique  

## Usage Examples

### Creating Audit Logs
```go
// UUIDs are generated automatically via BeforeCreate hook
auditLog := &AuditLog{
    UserID:     &userID,
    Method:     "GET",
    Path:       "/v1/examples",
    StatusCode: 200,
    // ID will be auto-generated as UUIDv7
}
```

### Querying by ID
```go
// Repository method now accepts string ID
log, err := repo.GetAuditLogByID("01983475-5757-778b-aa8f-2143995b7fe7")
```

### API Endpoints
```bash
# Get specific audit log
curl -X GET "http://localhost:3000/v1/audit-logs/01983475-5757-778b-aa8f-2143995b7fe7" \
  -H "Authorization: Bearer admin-api-key-789"
```

## Notes

- This change is for initial project setup only (no migration scripts needed)
- UUIDs are automatically generated when creating new audit log entries
- Time-sortable property ensures chronological ordering in database queries
- Consistent with the access table UUID implementation