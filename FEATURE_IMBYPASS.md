# URL Duplication Feature with imbypass=true

## Overview

This feature automatically duplicates URLs by adding the query parameter `imbypass=true` when purging content through Akamai.

## Implementation Details

### When it activates
- **Only for `purgeType: "urls"`** (not for cache-tags)
- Applies to both Akamai purge requests and post-purge GET requests

### Code Flow

```
User Request
    ↓
[PurgeHandler receives request]
    ↓
[Check if purgeType == "urls"]
    ↓
[duplicateURLsWithBypass()]
    ↓
Original URLs → [URL1, URL2, URL3]
    ↓
Duplicated → [URL1, URL1?imbypass=true, URL2, URL2?imbypass=true, URL3, URL3?imbypass=true]
    ↓
[Send to Akamai API]
    ↓
[If postPurgeRequest enabled]
    ↓
[executePurgeRequest() - duplicates URLs again for GET requests]
```

## Examples

### Example 1: Simple URL
**Input:**
```json
{
  "purgeType": "urls",
  "actionType": "invalidate",
  "environment": "production",
  "paths": ["https://example.com/page"]
}
```

**Akamai receives:**
- `https://example.com/page`
- `https://example.com/page?imbypass=true`

---

### Example 2: URL with existing query parameters
**Input:**
```json
{
  "purgeType": "urls",
  "actionType": "delete",
  "environment": "staging",
  "paths": ["https://example.com/page?foo=bar&baz=qux"]
}
```

**Akamai receives:**
- `https://example.com/page?foo=bar&baz=qux`
- `https://example.com/page?baz=qux&foo=bar&imbypass=true`

(Note: Query parameters may be reordered alphabetically by the URL parser)

---

### Example 3: Multiple URLs with post-purge requests
**Input:**
```json
{
  "purgeType": "urls",
  "actionType": "invalidate",
  "environment": "production",
  "postPurgeRequest": true,
  "paths": [
    "https://example.com/page1",
    "https://example.com/page2?id=123"
  ]
}
```

**Akamai receives (purge):**
- `https://example.com/page1`
- `https://example.com/page1?imbypass=true`
- `https://example.com/page2?id=123`
- `https://example.com/page2?id=123&imbypass=true`

**GET requests sent (post-purge, after 5 seconds):**
- GET `https://example.com/page1`
- GET `https://example.com/page1?imbypass=true`
- GET `https://example.com/page2?id=123`
- GET `https://example.com/page2?id=123&imbypass=true`

---

### Example 4: Cache tags (not affected)
**Input:**
```json
{
  "purgeType": "cache-tags",
  "actionType": "invalidate",
  "environment": "production",
  "paths": ["tag1", "tag2"]
}
```

**Akamai receives:**
- `tag1`
- `tag2`

(No duplication - feature only applies to URLs)

## Functions Added

### 1. `duplicateURLsWithBypass(paths []string, ctx v1alpha1.Context) []string`
- Creates a new slice with double the capacity
- For each path, adds the original and a version with `imbypass=true`
- Logs warnings if URL parsing fails

### 2. `addQueryParam(urlStr, key, value string) (string, error)`
- Safely parses URLs using `net/url.Parse()`
- Adds query parameter preserving existing ones
- Returns error for invalid URLs

## Error Handling

- Invalid URLs: Logged as warnings, original URL still processed
- URL parsing errors: Gracefully skipped, doesn't break entire request
- Continues processing other URLs even if one fails

## Testing

Run tests with:
```bash
go test ./internal/api/... -v
```

Test coverage includes:
- URLs without query parameters
- URLs with single query parameter
- URLs with multiple query parameters
- Path-only URLs (relative)
- Invalid URLs

All tests passing ✓

## Files Modified

- `internal/api/api.go` - Core implementation
- `README.md` - User documentation
- `internal/api/api_test.go` - Unit tests (new)
- `examples/purge_example.sh` - Usage example (new)

## Backward Compatibility

✅ Fully backward compatible
- Existing cache-tag purging unchanged
- Existing URL purging behavior extended (duplicates URLs)
- No breaking changes to API

## Performance Impact

- **Memory**: 2x allocation for URLs list (negligible for typical use cases)
- **Network**: 2x requests to Akamai (intentional - purges both variants)
- **Latency**: Minimal (URL parsing is fast)

## Configuration

No configuration changes required. Feature is always active for `purgeType: "urls"`.

---

**Status:** ✅ Implemented and tested
**Date:** 2026-01-30
