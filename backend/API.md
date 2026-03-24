# OpenCW Backend API

Base URL: `http://localhost:8080/v1`

## Authentication

Protected endpoints require:

- Header: `Authorization: Bearer <access_token>`
- Middleware: `AuthRequired` + `LoadUser`

Error response shape:

```json
{
  "code": "INVALID_TOKEN",
  "error": "invalid token"
}
```

## CW Progress

### GET `/cw/progress`

Returns all progress rows for the authenticated user.

Response `200 OK`:

```json
{
  "data": [
    {
      "lesson": "1",
      "char_wpm": 20,
      "eff_wpm": 15,
      "accuracy": 0.95,
      "created_at": "2026-03-24T10:00:00Z",
      "client_created_at": "2026-03-24T09:58:00Z"
    }
  ]
}
```

Error responses:

- `500 PROGRESS_QUERY_FAILED`

---

### PUT `/cw/progress`

Creates progress rows for the authenticated user from a **non-empty JSON array**.

> Breaking change: this endpoint now accepts **multiple progress items** in one request body (batch insert) for offline sync.

Request body:

```json
[
  {
    "lesson": 8,
    "char_wpm": 24,
    "eff_wpm": 20,
    "accuracy": 0.93,
    "client_created_at": "2026-03-24T09:55:10Z"
  },
  {
    "lesson": 9,
    "char_wpm": 25,
    "eff_wpm": 20,
    "accuracy": 0.91,
    "client_created_at": "2026-03-24T09:58:42Z"
  }
]
```

Field rules per item:

- `lesson`: required
- `char_wpm`: required, `5..50`
- `eff_wpm`: required, `5..50`
- `accuracy`: required, `0.0..1.0`
- `client_created_at`: optional, RFC3339 timestamp (recommended for offline ordering)

Response `201 Created`:

```json
{
  "message": "Progress Created",
  "count": 2
}
```

Error responses:

- `400 INVALID_REQUEST_BODY` - malformed JSON, wrong shape (object instead of array), or empty array
- `401 AUTH_HEADER_REQUIRED|INVALID_AUTH_HEADER_FORMAT|INVALID_TOKEN|EXPIRED_TOKEN|USER_NOT_FOUND`
- `500 PROGRESS_CREATE_FAILED`

## Common Error Codes (used by progress endpoints)

- `INVALID_REQUEST_BODY`
- `AUTH_HEADER_REQUIRED`
- `INVALID_AUTH_HEADER_FORMAT`
- `INVALID_TOKEN`
- `EXPIRED_TOKEN`
- `USER_NOT_FOUND`
- `PROGRESS_QUERY_FAILED`
- `PROGRESS_CREATE_FAILED`

