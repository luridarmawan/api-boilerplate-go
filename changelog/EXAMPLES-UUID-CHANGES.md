# Examples Table ID Migration to UUIDv7

## Summary

Successfully migrated the `examples` table's `id` field from `uint` to UUIDv7 string type for better performance and time-sortable properties.

## Changes Made

### 1. Model Updates (`internal/modules/example/example_model.go`)
- ✅ Changed `Example.ID` from `uint` to `string` with UUID type
- ✅ Added `BeforeCreate` hook to generate UUIDv7 automatically
- ✅ Added imports for `utils` and `gorm`

### 2. Repository Updates (`internal/modules/example/example_repository.go`)
- ✅ Updated `Repository` interface methods to accept `string` ID parameters
- ✅ Updated `GetExampleByID`, `SoftDeleteExample`, and `RestoreExample` methods
- ✅ All repository implementations now work with string IDs

### 3. Handler Updates (`internal/modules/example/example_handler.go`)
- ✅ Updated all handler methods to work with string ID parameters
- ✅ Removed `strconv.ParseUint` calls for ID parsing
- ✅ Updated Swagger documentation to use `string` for ID parameters
- ✅ Removed unused `strconv` import

### 4. Testing (`test/test.http`)
- ✅ Added comprehensive example API calls with UUID format
- ✅ Included CRUD operations with UUID examples

## Technical Details

### Before (Integer ID):
```go
type Example struct {
    ID uint `json:"id" gorm:"primaryKey"`
    // ... other fields
}
```

### After (UUIDv7):
```go
type Example struct {
    ID string `json:"id" gorm:"type:uuid;primaryKey"`
    // ... other fields
}

func (e *Example) BeforeCreate(tx *gorm.DB) error {
    if e.ID == "" {
        e.ID = utils.GenerateUUIDv7()
    }
    return nil
}
```

## API Changes

### Before (Integer ID):
```
GET /v1/examples/123
PUT /v1/examples/123
DELETE /v1/examples/123
POST /v1/examples/123/restore
```

### After (UUID):
```
GET /v1/examples/01983465-34ad-760b-87f5-b73601f6e281
PUT /v1/examples/01983465-34ad-760b-87f5-b73601f6e281
DELETE /v1/examples/01983465-34ad-760b-87f5-b73601f6e281
POST /v1/examples/01983465-34ad-760b-87f5-b73601f6e281/restore
```

## Benefits Achieved

1. **Time-sortable**: Examples are naturally ordered by creation time
2. **Better performance**: Reduced database index fragmentation
3. **Globally unique**: No collision risk across distributed systems
4. **Debugging friendly**: Timestamp embedded in ID for easier troubleshooting
5. **Consistent with other tables**: Matches access and audit_logs UUID implementation

## UUIDv7 Format for Examples

```
01983465-34ad-760b-87f5-b73601f6e281
├─────────────┤ ├─┤ ├─────────────────┤
│             │ │  │
│             │ │  └─ Random data (62 bits)
│             │ └─ Version (4 bits) + Random (12 bits)
└─ Timestamp (48 bits) - milliseconds since Unix epoch
```

## JSON Response Format

### Before:
```json
{
  "id": 123,
  "name": "Example Name",
  "description": "Example Description",
  "created_at": "2025-01-23T10:30:00Z",
  "updated_at": "2025-01-23T10:30:00Z",
  "status_id": 0
}
```

### After:
```json
{
  "id": "01983465-34ad-760b-87f5-b73601f6e281",
  "name": "Example Name",
  "description": "Example Description",
  "created_at": "2025-01-23T10:30:00Z",
  "updated_at": "2025-01-23T10:30:00Z",
  "status_id": 0
}
```

## Usage Examples

### Creating Examples:
```bash
# Create new example (UUID generated automatically)
curl -X POST "http://localhost:3000/v1/examples" \
  -H "Authorization: Bearer test-api-key-123" \
  -H "Content-Type: application/json" \
  -d '{"name": "Test Example", "description": "Testing UUID"}'
```

### Querying by UUID:
```bash
# Get specific example
curl -X GET "http://localhost:3000/v1/examples/01983465-34ad-760b-87f5-b73601f6e281" \
  -H "Authorization: Bearer test-api-key-123"

# Update example
curl -X PUT "http://localhost:3000/v1/examples/01983465-34ad-760b-87f5-b73601f6e281" \
  -H "Authorization: Bearer test-api-key-123" \
  -H "Content-Type: application/json" \
  -d '{"name": "Updated Name", "description": "Updated Description"}'

# Soft delete example
curl -X DELETE "http://localhost:3000/v1/examples/01983465-34ad-760b-87f5-b73601f6e281" \
  -H "Authorization: Bearer test-api-key-123"

# Restore example
curl -X POST "http://localhost:3000/v1/examples/01983465-34ad-760b-87f5-b73601f6e281/restore" \
  -H "Authorization: Bearer admin-api-key-789"
```

## Repository Methods Updated

All repository methods now work with string IDs:

```go
// Repository interface
type Repository interface {
    CreateExample(example *Example) error
    GetAllExamples() ([]Example, error)
    GetExampleByID(id string) (*Example, error)        // Updated
    UpdateExample(example *Example) error
    SoftDeleteExample(id string) error                 // Updated
    RestoreExample(id string) error                    // Updated
    GetDeletedExamples() ([]Example, error)
}
```

## Handler Methods Updated

All handler methods now process string IDs directly:

```go
// Before: Parse uint from string
id, err := strconv.ParseUint(idStr, 10, 32)
example, err := h.repo.GetExampleByID(uint(id))

// After: Use string directly
id := c.Params("id")
example, err := h.repo.GetExampleByID(id)
```

## Verification

✅ **Build Status**: Successful compilation  
✅ **UUID Tests**: All tests passing  
✅ **Time-sortable**: Verified lexicographic ordering by timestamp  
✅ **API Endpoints**: All CRUD operations updated  
✅ **Swagger Documentation**: Updated parameter types  

## Notes

- This change is for initial project setup only (no migration scripts needed)
- UUIDs are automatically generated when creating new examples
- Time-sortable property ensures chronological ordering in database queries
- Consistent with access and audit_logs table UUID implementations
- All existing functionality preserved with improved ID format