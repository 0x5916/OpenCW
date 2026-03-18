
# API Documentation

> Last updated: March 18, 2026

## Table of Contents
- Health
- Auth
  - POST /api/v1/auth/register
  - POST /api/v1/auth/login
  - POST /api/v1/auth/refresh
- CW
  - GET /api/v1/cw/settings
  - POST /api/v1/cw/settings
  - GET /api/v1/cw/progress
  - PUT /api/v1/cw/progress
- Page
  - GET /api/v1/page/settings
  - POST /api/v1/page/settings
- User
  - GET /api/v1/user/me
  - PUT /api/v1/user/email
  - PUT /api/v1/user/password

---

---
## GET `/api/v1/health`

**Description:** Simple health check endpoint. Returns basic server status and a unix timestamp.

**Authentication:** Not required

### Request

**Headers:**
| Header | Required | Description |
|---|---:|---|
| None | No | No special headers required. |

**Path Parameters:** None

**Query Parameters:** None

**Request Body:** None

### Response

**Success `200`:**
```json
{
  "status": "healthy",
  "timestamp": 1710681600
}
```

| Field | Type | Description |
|---|---|---|
| status | string | Always "healthy" when server is up |
| timestamp | integer | Unix timestamp (seconds) when response was generated |

**Errors:** None (this endpoint always returns 200 in normal operation)

### Example (cURL)
```bash
curl -X GET https://api.example.com/api/v1/health
```

---

## Auth

All auth endpoints live under /api/v1/auth and are unauthenticated.

### POST `/api/v1/auth/register`

**Description:** Create a new user account. Returns an access token and a refresh token on success.

**Authentication:** Not required

### Request

**Headers:**
| Header | Required | Description |
|---|---:|---|
| Content-Type | Yes | Must be `application/json` |

**Request Body** (`application/json`):
```json
{
  "username": "captain_jones",
  "email": "jones@example.com",
  "password": "Secur3P@ssw0rd"
}
```

| Field | Type | Required | Description |
|---|---|---:|---|
| username | string | Yes | 3–16 characters, letters/numbers, may include _ or - in middle. Custom validator `username` enforces pattern. |
| email | string | Yes | Must be valid email, max length 254 |
| password | string | Yes | 8–256 characters |

### Response

**Success `200`:**
```json
{
  "refresh_token": "<long-random-string>",
  "access_token": "<jwt-access-token>"
}
```

| Field | Type | Description |
|---|---|---|
| refresh_token | string | Opaque refresh token (raw). Keep this secure; used to obtain new access tokens. Expires in ~30 days. |
| access_token | string | JWT access token (Bearer). Expires in ~15 minutes. |

**Errors:**
| Status Code | Meaning |
|---:|---|
| 400 | Invalid request body or validation failed |
| 409 | Username or email already exists |
| 500 | Internal server error (database or hashing failure) |

### Example (cURL)
```bash
curl -X POST https://api.example.com/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"captain_jones","email":"jones@example.com","password":"Secur3P@ssw0rd"}'
```

---

### POST `/api/v1/auth/login`

**Description:** Authenticate a user by username or email and password. Returns an access token and a refresh token on success.

**Authentication:** Not required

### Request

**Headers:**
| Header | Required | Description |
|---|---:|---|
| Content-Type | Yes | `application/json` |

**Request Body** (`application/json`):
```json
{
  "identifier": "jones@example.com",
  "password": "Secur3P@ssw0rd"
}
```

| Field | Type | Required | Description |
|---|---|---:|---|
| identifier | string | Yes | Either username or email. If it contains an `@` it is treated as an email. |
| password | string | Yes | Plain-text password to verify. |

### Response

**Success `200`:**
```json
{
  "refresh_token": "<long-random-string>",
  "access_token": "<jwt-access-token>"
}
```

| Field | Type | Description |
|---|---|---|
| refresh_token | string | Raw refresh token (store safely). |
| access_token | string | JWT access token (use as `Authorization: Bearer <token>`) |

**Errors:**
| Status Code | Meaning |
|---:|---|
| 400 | Invalid request body |
| 401 | Invalid credentials (wrong identifier or password) |
| 500 | Internal server error |

### Example (cURL)
```bash
curl -X POST https://api.example.com/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"identifier":"jones@example.com","password":"Secur3P@ssw0rd"}'
```

---

### POST `/api/v1/auth/refresh`

**Description:** Exchange a valid refresh token for a new refresh token + access token pair. Refresh tokens are one-time use: submitting a valid refresh token revokes it and issues a brand-new refresh token.

**Authentication:** Not required (provides refresh token in body)

### Request

**Headers:**
| Header | Required | Description |
|---|---:|---|
| Content-Type | Yes | `application/json` |

**Request Body** (`application/json`):
```json
{
  "refresh_token": "<raw-refresh-token>"
}
```

| Field | Type | Required | Description |
|---|---|---:|---|
| refresh_token | string | Yes | Raw refresh token previously issued by the server. |

### Response

**Success `200`:**
```json
{
  "refresh_token": "<new-raw-refresh-token>",
  "access_token": "<new-jwt-access-token>"
}
```

| Field | Type | Description |
|---|---|---|
| refresh_token | string | New raw refresh token (previous one is revoked) |
| access_token | string | New JWT access token |

**Errors:**
| Status Code | Meaning |
|---:|---|
| 400 | Invalid request body |
| 401 | Invalid refresh token or expired refresh token |
| 500 | Internal server error (DB failure) |

### Notes
- Access tokens expire after ~15 minutes. Refresh tokens expire after ~30 days and are single-use (they are revoked when exchanged).

### Example (cURL)
```bash
curl -X POST https://api.example.com/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{"refresh_token":"<raw-refresh-token>"}'
```

---

## CW

All /api/v1/cw endpoints require a valid access token in the Authorization header (see "Authentication" below). Additionally the server loads the full user model into the request context.

Authentication for protected endpoints:
- Header: Authorization: Bearer <access_token>
- Access tokens are JWTs signed by the server and validated by middleware. If missing, malformed, expired, or invalid → 401.

### GET `/api/v1/cw/settings`

**Description:** Retrieve the current user's CW (typing) settings. If the user has not saved settings yet, default server settings are returned.

**Authentication:** Required — Bearer access token in `Authorization` header

### Request

**Headers:**
| Header | Required | Description |
|---|---:|---|
| Authorization | Yes | `Bearer <access_token>` |

**Request Body:** None

### Response

**Success `200`:**
```json
{
  "char_wpm": 20,
  "eff_wpm": 10,
  "freq": 600,
  "start_delay": 0.5
}
```

| Field | Type | Description |
|---|---|---|
| char_wpm | integer | Target character WPM (5–50) |
| eff_wpm | integer | Effective WPM (5–50) |
| freq | integer | Frequency / pacing value (300–2000) |
| start_delay | float | Seconds to wait before starting (0.0–10.0) |

**Errors:**
| Status Code | Meaning |
|---:|---|
| 401 | Missing or invalid access token |
| 500 | Failed to read settings from database |

### Example (cURL)
```bash
curl -X GET https://api.example.com/api/v1/cw/settings \
  -H "Authorization: Bearer <access_token>"
```

---

### POST `/api/v1/cw/settings`

**Description:** Create or update the user's CW settings.

**Authentication:** Required — Bearer access token

### Request

**Headers:**
| Header | Required | Description |
|---|---:|---|
| Authorization | Yes | `Bearer <access_token>` |
| Content-Type | Yes | `application/json` |

**Request Body** (`application/json`):
```json
{
  "char_wpm": 25,
  "eff_wpm": 12,
  "freq": 800,
  "start_delay": 0.7
}
```

| Field | Type | Required | Validation | Description |
|---|---|---:|---|---|
| char_wpm | integer | Yes | 5–50 | Desired character WPM |
| eff_wpm | integer | Yes | 5–50 | Effective WPM |
| freq | integer | Yes | 300–2000 | Frequency / pacing |
| start_delay | number | Yes | 0.0–10.0 | Start delay in seconds |

### Response

**Success `200`:**
```json
{
  "message": "Settings updated"
}
```

| Field | Type | Description |
|---|---|---|
| message | string | Human-readable confirmation |

**Errors:**
| Status Code | Meaning |
|---:|---|
| 400 | Invalid request body or validation failure |
| 401 | Missing or invalid access token |
| 500 | Failed to update settings in database |

### Example (cURL)
```bash
curl -X POST https://api.example.com/api/v1/cw/settings \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{"char_wpm":25,"eff_wpm":12,"freq":800,"start_delay":0.7}'
```

---

### GET `/api/v1/cw/progress`

**Description:** Retrieve all saved lesson progress records for the authenticated user.

**Authentication:** Required — Bearer access token

### Request

**Headers:**
| Header | Required | Description |
|---|---:|---|
| Authorization | Yes | `Bearer <access_token>` |

### Response

**Success `200`:**
```json
{
  "data": [
	{
	  "lesson": "1",
	  "char_wpm": 24,
	  "eff_wpm": 12,
	  "accuracy": 0.94,
	  "created_at": "2026-03-01T12:34:56Z"
	}
  ]
}
```

| Field | Type | Description |
|---|---|---|
| data | array | List of progress entries |
| lesson | string | NOTE: returned as string (the database stores lesson as integer). Frontend should treat as string or convert to number as needed. |
| char_wpm | integer | Character-level WPM recorded |
| eff_wpm | integer | Effective WPM recorded |
| accuracy | float | Value between 0.0 and 1.0 |
| created_at | string (ISO 8601) | Timestamp when the record was created |

**Errors:**
| Status Code | Meaning |
|---:|---|
| 401 | Missing or invalid access token |
| 500 | Database read failure |

### Example (cURL)
```bash
curl -X GET https://api.example.com/api/v1/cw/progress \
  -H "Authorization: Bearer <access_token>"
```

---

### PUT `/api/v1/cw/progress`

**Description:** Submit a new lesson progress record for the authenticated user.

**Authentication:** Required — Bearer access token

### Request

**Headers:**
| Header | Required | Description |
|---|---:|---|
| Authorization | Yes | `Bearer <access_token>` |
| Content-Type | Yes | `application/json` |

**Request Body** (`application/json`):
```json
{
  "lesson": 1,
  "char_wpm": 24,
  "eff_wpm": 12,
  "accuracy": 0.94
}
```

| Field | Type | Required | Validation | Description |
|---|---|---:|---|---|
| lesson | integer | Yes | required | Lesson identifier (integer) |
| char_wpm | integer | Yes | 5–50 | Character WPM measured |
| eff_wpm | integer | Yes | 5–50 | Effective WPM measured |
| accuracy | number | Yes | 0.0–1.0 | Accuracy ratio (0.0–1.0) |

### Response

**Success `201`:**
```json
{
  "message": "Progress Created"
}
```

| Field | Type | Description |
|---|---|---|
| message | string | Confirmation message |

**Errors:**
| Status Code | Meaning |
|---:|---|
| 400 | Invalid request body or validation failure |
| 401 | Missing or invalid access token |
| 500 | Failed to write progress to database |

### Example (cURL)
```bash
curl -X PUT https://api.example.com/api/v1/cw/progress \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{"lesson":1,"char_wpm":24,"eff_wpm":12,"accuracy":0.94}'
```

---

## Page

Routes under `/api/v1/page` manage user-visible page settings.

### GET `/api/v1/page/settings`

**Description:** Retrieve the current user's page settings. If not set, server defaults are returned.

**Authentication:** Required — Bearer access token

### Request

**Headers:**
| Header | Required | Description |
|---|---:|---|
| Authorization | Yes | `Bearer <access_token>` |

### Response

**Success `200`:**
```json
{
  "theme": "auto",
  "language": "auto",
  "cur_lesson": 0
}
```

| Field | Type | Description |
|---|---|---|
| theme | string | One of `auto`, `dark`, `light` |
| language | string | Language code or `auto` |
| cur_lesson | integer | Currently selected lesson index |

**Errors:**
| Status Code | Meaning |
|---:|---|
| 401 | Missing or invalid access token |
| 500 | Database read failure |

### Example (cURL)
```bash
curl -X GET https://api.example.com/api/v1/page/settings \
  -H "Authorization: Bearer <access_token>"
```

---

### POST `/api/v1/page/settings`

**Description:** Create or update the user's page settings.

**Authentication:** Required — Bearer access token

### Request

**Headers:**
| Header | Required | Description |
|---|---:|---|
| Authorization | Yes | `Bearer <access_token>` |
| Content-Type | Yes | `application/json` |

**Request Body** (`application/json`):
```json
{
  "theme": "dark",
  "language": "en",
  "cur_lesson": 3
}
```

| Field | Type | Required | Validation | Description |
|---|---|---:|---|---|
| theme | string | Yes | one of `auto`, `dark`, `light` | UI theme preference |
| language | string | Yes | required | Language code or `auto` |
| cur_lesson | integer | Yes | required | Current lesson index |

### Response

**Success `200`:**
```json
{
  "message": "Settings updated"
}
```

**Errors:**
| Status Code | Meaning |
|---:|---|
| 400 | Invalid input or validation failure |
| 401 | Missing or invalid access token |
| 500 | Database write failure |

### Example (cURL)
```bash
curl -X POST https://api.example.com/api/v1/page/settings \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{"theme":"dark","language":"en","cur_lesson":3}'
```

---

## User

User-related endpoints are under `/api/v1/user` and require a valid access token.

### GET `/api/v1/user/me`

**Description:** Retrieve information about the currently authenticated user.

**Authentication:** Required — Bearer access token

### Request

**Headers:**
| Header | Required | Description |
|---|---:|---|
| Authorization | Yes | `Bearer <access_token>` |

### Response

**Success `200`:**
```json
{
  "username": "captain_jones",
  "email": "jones@example.com",
  "created_at": "2026-02-15T09:21:00Z"
}
```

| Field | Type | Description |
|---|---|---|
| username | string | User's username |
| email | string | User's email address |
| created_at | string (ISO 8601) | Account creation timestamp |

**Errors:**
| Status Code | Meaning |
|---:|---|
| 401 | Missing or invalid access token or user not found |

### Example (cURL)
```bash
curl -X GET https://api.example.com/api/v1/user/me \
  -H "Authorization: Bearer <access_token>"
```

---

### PUT `/api/v1/user/email`

**Description:** Update the authenticated user's email address.

**Authentication:** Required — Bearer access token

### Request

**Headers:**
| Header | Required | Description |
|---|---:|---|
| Authorization | Yes | `Bearer <access_token>` |
| Content-Type | Yes | `application/json` |

**Request Body** (`application/json`):
```json
{
  "email": "new.email@example.com"
}
```

| Field | Type | Required | Validation | Description |
|---|---|---:|---|---|
| email | string | Yes | valid email, max 254 | New email address (must be different from current) |

### Response

**Success `200`:**
```json
{
  "message": "Email updated"
}
```

**Errors:**
| Status Code | Meaning |
|---:|---|
| 400 | Invalid request body or new email equals current email |
| 401 | Missing or invalid access token |
| 409 | Email already in use by another account |
| 500 | Database error while updating email |

### Example (cURL)
```bash
curl -X PUT https://api.example.com/api/v1/user/email \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{"email":"new.email@example.com"}'
```

---

### PUT `/api/v1/user/password`

**Description:** Change the authenticated user's password. The endpoint requires the current password for verification.

**Authentication:** Required — Bearer access token

### Request

**Headers:**
| Header | Required | Description |
|---|---:|---|
| Authorization | Yes | `Bearer <access_token>` |
| Content-Type | Yes | `application/json` |

**Request Body** (`application/json`):
```json
{
  "old_password": "Secur3P@ssw0rd",
  "new_password": "N3wSecur3P@ss"
}
```

| Field | Type | Required | Validation | Description |
|---|---|---:|---|---|
| old_password | string | Yes | 8–256 chars | Current password (verified) |
| new_password | string | Yes | 8–256 chars | New password to set |

### Response

**Success `200`:**
```json
{
  "message": "Password updated"
}
```

**Errors:**
| Status Code | Meaning |
|---:|---|
| 400 | Invalid request body |
| 401 | Invalid current password or missing/invalid access token |
| 500 | Internal server error (hashing or DB update failure) |

### Example (cURL)
```bash
curl -X PUT https://api.example.com/api/v1/user/password \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{"old_password":"Secur3P@ssw0rd","new_password":"N3wSecur3P@ss"}'
```

---

## Additional Notes

- Authentication middleware requirements:
  - All protected endpoints require the `Authorization` header in the form `Bearer <token>`.
  - Tokens are validated as JWTs and must include a numeric subject (the user ID).
  - If the Authorization header is missing, malformed, or token is invalid/expired, endpoints return 401 with a JSON error.
- The `LoadUser` middleware fetches the full user record from the database using the ID from the access token and places it in request context under the key `user`. If the user is not found, request returns 401.
- CORS: the server sets Access-Control-Allow-Origin based on configuration (explicit allowlist via env or, in non-release mode, allows all origins). The server includes Authorization in Access-Control-Allow-Headers.
- Rate limiting: none implemented in codebase (no per-route rate limits found).

---

## Protected miscellaneous

### GET `/api/v1/hello`

**Description:** Simple protected endpoint that returns a greeting message for authenticated users. Useful for quick token checks.

**Authentication:** Required — Bearer access token

### Request

**Headers:**
| Header | Required | Description |
|---|---:|---|
| Authorization | Yes | `Bearer <access_token>` |

### Response

**Success `200`:**
```json
{
  "message": "Hello, authenticated user!"
}
```

| Field | Type | Description |
|---|---|---|
| message | string | Greeting message |

**Errors:**
| Status Code | Meaning |
|---:|---|
| 401 | Missing or invalid access token |


---

## Error Reference

| Status Code | Meaning |
|---|---|
| 400 | Bad Request — Invalid or missing parameters |
| 401 | Unauthorized — Missing or invalid token |
| 403 | Forbidden — Insufficient permissions |
| 404 | Not Found — Resource does not exist |
| 429 | Too Many Requests — Rate limit exceeded |
| 500 | Internal Server Error — Unexpected server failure |


