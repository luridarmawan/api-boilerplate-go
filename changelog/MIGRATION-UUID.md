# Migration: Access Table ID to UUIDv7

This document explains the migration process to change the `access` table's `id` field from `uint` to UUIDv7.

## What Changed

1. **Model Changes:**
   - `User.ID` field changed from `uint` to `string` (UUID)
   - Added UUIDv7 generation in `BeforeCreate` hook
   - Updated `GetID()` method to return `string`

2. **Repository Changes:**
   - All repository methods now accept `string` ID parameters
   - Updated method signatures for `UpdateExpiredDate`, `UpdateRateLimit`, and `GetUserByID`

3. **Handler Changes:**
   - Removed `strconv.ParseUint` calls
   - Updated Swagger documentation to use `string` for ID parameters
   - Direct string ID handling from URL parameters

4. **Interface Changes:**
   - Updated `types.User` interface `GetID()` method to return `string`

## UUIDv7 Benefits

- **Time-sortable**: UUIDs are generated with timestamp prefix, maintaining chronological order
- **Globally unique**: No collision risk across distributed systems
- **Better performance**: Sortable nature reduces database index fragmentation
- **Future-proof**: Scalable across multiple servers/databases
- **Database-friendly**: Better for primary keys than random UUIDs (v4)
- **Debugging-friendly**: Timestamp embedded in ID makes troubleshooting easier

## UUIDv7 Format

```
01983465-34ad-760b-87f5-b73601f6e281
├─────────────┤ ├─┤ ├─────────────────┤
│             │ │  │
│             │ │  └─ Random data (62 bits)
│             │ └─ Version (4 bits) + Random (12 bits)
└─ Timestamp (48 bits) - milliseconds since Unix epoch
```

The timestamp prefix ensures that newer records have lexicographically larger IDs, making them naturally sortable by creation time.

## Migration Steps

### 1. Pre-Migration Checklist

- [ ] Backup your database
- [ ] Test the migration script in a staging environment
- [ ] Notify API consumers about the breaking change
- [ ] Schedule maintenance window

### 2. Database Migration (REQUIRED)

Before deploying the new code, run the SQL migration:

```bash
# Connect to your database
psql -d your_database -f scripts/migrate-access-to-uuid.sql
```

### 3. Code Deployment

Deploy the updated code after the database migration is complete.

### 4. Verification

After deployment, verify the changes:

```sql
-- Check that IDs are now UUIDs
SELECT id, name, email FROM access LIMIT 5;

-- Verify data type
SELECT pg_typeof(id) FROM access LIMIT 1;

-- Check audit_logs table
SELECT user_id, user_email FROM audit_logs WHERE user_id IS NOT NULL LIMIT 5;

-- Verify UUIDv7 time-sortable property
SELECT id, created_at FROM access ORDER BY id LIMIT 10;
```

### 5. Post-Migration Tasks

- [ ] Update API documentation
- [ ] Notify API consumers that migration is complete
- [ ] Monitor application logs for any issues
- [ ] Clean up backup tables after confirming everything works

## API Changes

### Before (Integer ID):
```
PUT /v1/access/123/rate-limit
PUT /v1/access/123/expired-date
DELETE /v1/access/123/expired-date
```

### After (UUID):
```
PUT /v1/access/550e8400-e29b-41d4-a716-446655440000/rate-limit
PUT /v1/access/550e8400-e29b-41d4-a716-446655440000/expired-date
DELETE /v1/access/550e8400-e29b-41d4-a716-446655440000/expired-date
```

## Important Notes

1. **Breaking Change**: This is a breaking change for any clients using the API
2. **Backup**: The migration script creates a backup table (`access_backup`)
3. **Foreign Keys**: Check for any tables referencing `access.id` and update them
4. **Testing**: Test thoroughly in staging environment before production deployment

## Rollback Plan

If rollback is needed:

1. Stop the application
2. Restore from `access_backup` table
3. Deploy the previous version of the code

```sql
-- Rollback (if needed)
DROP TABLE access;
ALTER TABLE access_backup RENAME TO access;
```