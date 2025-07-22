
## ⏳ API Key Expiration System

This API features a comprehensive API key expiration management system:

### Key Features:
- ✅ **Flexible Expiration**: API keys can be set to expire at specific dates
- ✅ **Never Expires**: API keys can be configured to never expire (NULL value)
- ✅ **Auto Validation**: Expired API keys are automatically rejected
- ✅ **Management API**: Dedicated endpoints to set and clear expiration dates
- ✅ **Permission Based**: Only admins with "access:manage" permission can configure



**Set Expiration Date:**
```bash
# Set API key untuk expired dalam 30 hari
curl -X PUT "http://localhost:3000/v1/access/1/expired-date" \
  -H "Authorization: Bearer admin-api-key-789" \
  -H "Content-Type: application/json" \
  -d '{
    "expired_date": "2025-08-22T00:00:00Z"
  }'
```

**Remove Expiration (Never Expires):**
```bash
# Set API key untuk tidak pernah expired
curl -X DELETE "http://localhost:3000/v1/access/1/expired-date" \
  -H "Authorization: Bearer admin-api-key-789"
```


Manual example:
```sql
UPDATE access SET token_expired_at = NOW() + INTERVAL '2 days';
```

