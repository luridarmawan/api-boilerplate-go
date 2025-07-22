# Set rate limit untuk user (300 requests per menit)
curl -X PUT "http://localhost:3000/v1/access/1/rate-limit" \
  -H "Authorization: Bearer admin-api-key-789" \
  -H "Content-Type: application/json" \
  -d '{
    "rate_limit": 300
  }'
