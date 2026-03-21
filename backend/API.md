# API Documentation
*Generated: 2026-03-21*
*Base URL: http://localhost:8080*

## Endpoints Overview
| Method | Path | Description |
|--------|------|-------------|
| GET    | /v1/health | Health check |
| POST   | /v1/auth/register | Register new user (returns refresh + access tokens) |
| POST   | /v1/auth/login | Login with username/email + password (returns refresh + access tokens) |
| POST   | /v1/auth/refresh | Exchange refresh token for new refresh + access tokens |
| GET    | /v1/settings/all | Get both CW and Page settings for authenticated user |
| GET    | /v1/settings/cw | Get CW settings for authenticated user |
| GET    | /v1/settings/page | Get Page settings for authenticated user |
| POST   | /v1/settings/cw | Update CW settings (authenticated) |
| POST   | /v1/settings/page | Update Page settings (authenticated) |
| GET    | /v1/user/me | Get current authenticated user's info |
| PUT    | /v1/user/callsign | Update current user's call sign |
| PUT    | /v1/user/email | Update current user's email |
| PUT    | /v1/user/password | Update current user's password |
| GET    | /v1/cw/progress | Get user's CW progress records |
| PUT    | /v1/cw/progress | Add a progress record for the user |
| GET    | /v1/hello | Simple authenticated hello endpoint (for testing auth + loaduser)

## Authentication
- Protected routes under the protected group use two middlewares: `AuthRequired()` and `LoadUser(db)`.
- `AuthRequired()` requires an `Authorization: Bearer <token>` header and validates a JWT signed with the app secret.
- `LoadUser()` reads the user ID from the validated token and loads the `User` record into the request context as `user`.

## Error codes (global)
The project defines centralized string error codes in `common/error.go`. Common codes used across endpoints:

| Error Code | Meaning |
|------------|---------|
| `INVALID_REQUEST_BODY` | Request JSON is invalid or fails binding |
| `INTERNAL_SERVER_ERROR` | Generic server error |
| `DATABASE_FAILURE` | Low-level DB failure |
| `INVALID_CREDENTIALS` | Username/email or password incorrect |
| `CONFLICT` | Generic conflict on create/update |
| `INVALID_TOKEN` | Token invalid |
| `EXPIRED_TOKEN` | Refresh token expired |
| `AUTH_HEADER_REQUIRED` | Authorization header missing |
| `INVALID_AUTH_HEADER_FORMAT` | Authorization header not in `Bearer <token>` format |
| `USER_NOT_FOUND` | User record not found |
| `SETTINGS_FETCH_FAILED` | Failed to fetch settings |
| `SETTINGS_UPDATE_FAILED` | Failed to update settings |
| `PROGRESS_QUERY_FAILED` | Failed to query progress records |
| `PROGRESS_CREATE_FAILED` | Failed to create a progress record |
| `PASSWORD_HASH_FAILED` | Password hashing failed |
| `TOKEN_ISSUE_FAILED` | Failed to issue/generate tokens |
| `EMAIL_ALREADY_IN_USE` | Email uniqueness conflict |
| `USERNAME_ALREADY_IN_USE` | Username uniqueness conflict |
| `USERNAME_EMAIL_ALREADY_IN_USE` | Username and email both in use |
| `EMAIL_UNCHANGED` | New email equals current email |
| `CALL_SIGN_ALREADY_IN_USE` | Call sign uniqueness conflict |

---

## Health

### GET /v1/health
**Simple health check**

**Authentication**: None

**Path Parameters**: None

**Query Parameters**: None

**Request Body**: None

**Response (200)**:
```json
{
  "status": "healthy",
  "timestamp": 1711060200
}
```

**Example cURL**:
```bash
curl -X GET http://localhost:8080/v1/health
```

## Authentication Endpoints

### POST /v1/auth/register
**Registers a new user and returns refresh + access tokens**

**Authentication**: None

**Path Parameters**: None

**Query Parameters**: None

**Request Body**:
```json
{
  "username": "alan_yeung",
  "email": "alan@example.com",
  "password": "min8chars123"
}
```
*Validation: username (custom `username` rule: ASCII letters/numbers, underscores/dashes allowed, length approx 3-16 as enforced by regex), email (valid format), password (min 8 chars)*

**Response (200)**:
```json
{
  "refresh_token": "<refresh-token-raw>",
  "access_token": "<jwt-access-token>"
}
```

**Error Responses**:
| Status | Error Code | Message |
|--------|------------|---------|
| 400 | `INVALID_REQUEST_BODY` | "Invalid request body" |
| 409 | `USERNAME_EMAIL_ALREADY_IN_USE` | "Username and email already exists" |
| 409 | `USERNAME_ALREADY_IN_USE` | "Username already exists" |
| 409 | `EMAIL_ALREADY_IN_USE` | "Email already exists" |
| 409 | `CONFLICT` | "Registration conflict, please try again" |
| 500 | `DATABASE_FAILURE` | "Database failure" |
| 500 | `PASSWORD_HASH_FAILED` | "Failed to hash password" |
| 500 | `INTERNAL_SERVER_ERROR` | "Failed to create user" |
| 500 | `TOKEN_ISSUE_FAILED` | "Failed to issue token, try to login." |

**Example cURL**:
```bash
curl -X POST http://localhost:8080/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"alan_yeung","email":"alan@example.com","password":"securepass123"}'
```

### POST /v1/auth/login
**Logs in a user and returns refresh + access tokens**

**Authentication**: None

**Request Body**:
```json
{
  "identifier": "alan_yeung or alan@example.com",
  "password": "userpassword"
}
```

**Response (200)**:
```json
{
  "refresh_token": "<refresh-token-raw>",
  "access_token": "<jwt-access-token>"
}
```

**Error Responses**:
| Status | Error Code | Message |
|--------|------------|---------|
| 400 | `INVALID_REQUEST_BODY` | "Invalid request body" |
| 401 | `INVALID_CREDENTIALS` | "Invalid credentials" |
| 500 | `INTERNAL_SERVER_ERROR` | "internal error" |
| 500 | `TOKEN_ISSUE_FAILED` | "Failed to issue token" |

**Example cURL**:
```bash
curl -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"identifier":"alan_yeung","password":"securepass123"}'
```

### POST /v1/auth/refresh
**Exchange a refresh token for a new refresh token + access token**

**Authentication**: None (uses refresh token in body)

**Request Body**:
```json
{
  "refresh_token": "<raw-refresh-token>"
}
```

**Response (200)**:
```json
{
  "refresh_token": "<new-refresh-token-raw>",
  "access_token": "<new-access-token>"
}
```

**Error Responses**:
| Status | Error Code | Message |
|--------|------------|---------|
| 400 | `INVALID_REQUEST_BODY` | "Invalid request body" |
| 401 | `INVALID_TOKEN` | "Invalid refresh token" |
| 401 | `EXPIRED_TOKEN` | "Refresh token expired" |
| 500 | `INTERNAL_SERVER_ERROR` | "Failed to revoke/replace refresh token" |
| 500 | `TOKEN_ISSUE_FAILED` | "Failed to generate access token" |

**Example cURL**:
```bash
curl -X POST http://localhost:8080/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{"refresh_token":"<raw-refresh-token>"}'
```

## Settings

### GET /v1/settings/all
**Returns CW and Page settings for the authenticated user**

**Authentication**: Bearer JWT

**Response (200)**:
```json
{
  "cw_settings": {
	"char_wpm": 20,
	"eff_wpm": 10,
	"freq": 600,
	"start_delay": 0.5
  },
  "page_settings": {
	"theme": "auto",
	"language": "auto",
	"cur_lesson": 0
  }
}
```

**Error Responses**:
| Status | Error Code | Message |
|--------|------------|---------|
| 500 | `SETTINGS_FETCH_FAILED` | "Failed to get settings" |

**Example cURL**:
```bash
curl -X GET http://localhost:8080/v1/settings/all \
  -H "Authorization: Bearer <access-token>"
```

### GET /v1/settings/cw
**Get CW settings**

**Authentication**: Bearer JWT

**Response (200)**:
```json
{
  "char_wpm": 20,
  "eff_wpm": 10,
  "freq": 600,
  "start_delay": 0.5
}
```

**Error Responses**:
| Status | Error Code | Message |
|--------|------------|---------|
| 500 | `SETTINGS_FETCH_FAILED` | "Failed to get settings" |

**Example cURL**:
```bash
curl -X GET http://localhost:8080/v1/settings/cw \
  -H "Authorization: Bearer <access-token>"
```

### GET /v1/settings/page
**Get Page settings**

**Authentication**: Bearer JWT

**Response (200)**:
```json
{
  "theme": "auto",
  "language": "auto",
  "cur_lesson": 0
}
```

**Error Responses**:
| Status | Error Code | Message |
|--------|------------|---------|
| 500 | `SETTINGS_FETCH_FAILED` | "Failed to get settings" |

**Example cURL**:
```bash
curl -X GET http://localhost:8080/v1/settings/page \
  -H "Authorization: Bearer <access-token>"
```

### POST /v1/settings/cw
**Update or create CW settings for the authenticated user**

**Authentication**: Bearer JWT

**Request Body**:
```json
{
  "char_wpm": 25,
  "eff_wpm": 12,
  "freq": 700,
  "start_delay": 0.5
}
```
*Validation: char_wpm 5-50, eff_wpm 5-50, freq 300-2000, start_delay 0.0-10.0*

**Response (200)**:
```json
{
  "message": "Settings updated"
}
```

**Error Responses**:
| Status | Error Code | Message |
|--------|------------|---------|
| 400 | `INVALID_REQUEST_BODY` | "Invalid request body" |
| 500 | `SETTINGS_UPDATE_FAILED` | "Failed to update settings" |

**Example cURL**:
```bash
curl -X POST http://localhost:8080/v1/settings/cw \
  -H "Authorization: Bearer <access-token>" \
  -H "Content-Type: application/json" \
  -d '{"char_wpm":25,"eff_wpm":12,"freq":700,"start_delay":0.5}'
```

### POST /v1/settings/page
**Update or create Page settings for the authenticated user**

**Authentication**: Bearer JWT

**Request Body**:
```json
{
  "theme": "dark",
  "language": "en",
  "cur_lesson": 1
}
```
*Validation: theme one of [auto, dark, light]*

**Response (200)**:
```json
{
  "message": "Settings updated"
}
```

**Error Responses**:
| Status | Error Code | Message |
|--------|------------|---------|
| 400 | `INVALID_REQUEST_BODY` | "Invalid request body" |
| 500 | `SETTINGS_UPDATE_FAILED` | "Failed to update settings" |

**Example cURL**:
```bash
curl -X POST http://localhost:8080/v1/settings/page \
  -H "Authorization: Bearer <access-token>" \
  -H "Content-Type: application/json" \
  -d '{"theme":"dark","language":"en","cur_lesson":1}'
```

## User Management

### GET /v1/user/me
**Returns profile for the authenticated user**

**Authentication**: Bearer JWT

**Response (200)**:
```json
{
  "call_sign": null,
  "username": "alan_yeung",
  "email": "alan@example.com",
  "email_verified": false,
  "created_at": "2026-03-21T18:30:00Z"
}
```

**Example cURL**:
```bash
curl -X GET http://localhost:8080/v1/user/me \
  -H "Authorization: Bearer <access-token>"
```

### PUT /v1/user/callsign
**Update the authenticated user's call sign**

**Authentication**: Bearer JWT

**Request Body**:
```json
{
  "call_sign": "KJ7ABC"
}
```

**Response (200)**:
```json
{
  "message": "Call sign updated"
}
```

**Error Responses**:
| Status | Error Code | Message |
|--------|------------|---------|
| 400 | `INVALID_REQUEST_BODY` | "Invalid request body" |
| 409 | `CALL_SIGN_ALREADY_IN_USE` | "Call sign already in use" |
| 500 | `INTERNAL_SERVER_ERROR` | "Failed to update call sign" |

**Example cURL**:
```bash
curl -X PUT http://localhost:8080/v1/user/callsign \
  -H "Authorization: Bearer <access-token>" \
  -H "Content-Type: application/json" \
  -d '{"call_sign":"KJ7ABC"}'
```

### PUT /v1/user/email
**Update the authenticated user's email**

**Authentication**: Bearer JWT

**Request Body**:
```json
{
  "email": "newemail@example.com"
}
```

**Response (200)**:
```json
{
  "message": "Email updated"
}
```

**Error Responses**:
| Status | Error Code | Message |
|--------|------------|---------|
| 400 | `INVALID_REQUEST_BODY` | "Invalid request body" |
| 400 | `EMAIL_UNCHANGED` | "New email must be different from current email" |
| 409 | `EMAIL_ALREADY_IN_USE` | "Email already in use" |
| 500 | `INTERNAL_SERVER_ERROR` | "Failed to update email" |

**Example cURL**:
```bash
curl -X PUT http://localhost:8080/v1/user/email \
  -H "Authorization: Bearer <access-token>" \
  -H "Content-Type: application/json" \
  -d '{"email":"newemail@example.com"}'
```

### PUT /v1/user/password
**Update the authenticated user's password**

**Authentication**: Bearer JWT

**Request Body**:
```json
{
  "old_password": "currentpass",
  "new_password": "newsecurepass"
}
```

**Response (200)**:
```json
{
  "message": "Password updated"
}
```

**Error Responses**:
| Status | Error Code | Message |
|--------|------------|---------|
| 400 | `INVALID_REQUEST_BODY` | "Invalid request body" |
| 401 | `INVALID_CREDENTIALS` | "Invalid credentials" |
| 500 | `PASSWORD_HASH_FAILED` | "Failed to hash password" |
| 500 | `INTERNAL_SERVER_ERROR` | "Failed to update password" |

**Example cURL**:
```bash
curl -X PUT http://localhost:8080/v1/user/password \
  -H "Authorization: Bearer <access-token>" \
  -H "Content-Type: application/json" \
  -d '{"old_password":"currentpass","new_password":"newsecurepass"}'
```

## Progress

### GET /v1/cw/progress
**Get all progress records for the authenticated user**

**Authentication**: Bearer JWT

**Response (200)**:
```json
{
  "data": [
	{
	  "lesson": "1",
	  "char_wpm": 20,
	  "eff_wpm": 12,
	  "accuracy": 0.95,
	  "created_at": "2026-03-21T18:30:00Z"
	}
  ]
}
```

**Error Responses**:
| Status | Error Code | Message |
|--------|------------|---------|
| 500 | `PROGRESS_QUERY_FAILED` | "failed to query progress" |

**Example cURL**:
```bash
curl -X GET http://localhost:8080/v1/cw/progress \
  -H "Authorization: Bearer <access-token>"
```

### PUT /v1/cw/progress
**Create a new progress record for the authenticated user**

**Authentication**: Bearer JWT

**Request Body**:
```json
{
  "lesson": 1,
  "char_wpm": 20,
  "eff_wpm": 12,
  "accuracy": 0.95
}
```

**Response (201)**:
```json
{
  "message": "Progress Created"
}
```

**Error Responses**:
| Status | Error Code | Message |
|--------|------------|---------|
| 400 | `INVALID_REQUEST_BODY` | "Invalid request body" |
| 500 | `PROGRESS_CREATE_FAILED` | "Failed to create progress" |

**Example cURL**:
```bash
curl -X PUT http://localhost:8080/v1/cw/progress \
  -H "Authorization: Bearer <access-token>" \
  -H "Content-Type: application/json" \
  -d '{"lesson":1,"char_wpm":20,"eff_wpm":12,"accuracy":0.95}'
```

## Misc

### GET /v1/hello
**Authenticated hello endpoint; demonstrates user loaded via middleware**

**Authentication**: Bearer JWT

**Response (200)**:
```json
{
  "message": "Hello, authenticated user {<username>}!"
}
```

**Example cURL**:
```bash
curl -X GET http://localhost:8080/v1/hello \
  -H "Authorization: Bearer <access-token>"
```

## Structs and Schemas (selected)
- common.RegisterInput
```json
{
  "username": "string",
  "email": "string",
  "password": "string"
}
```
- common.LoginInput
```json
{
  "identifier": "username or email",
  "password": "string"
}
```
- common.RefreshInput
```json
{
  "refresh_token": "string"
}
```
- common.CWSettingsInput
```json
{
  "char_wpm": 20,
  "eff_wpm": 10,
  "freq": 600,
  "start_delay": 0.5
}
```
- common.PageSettingsInput
```json
{
  "theme": "auto|dark|light",
  "language": "string",
  "cur_lesson": 0
}
```
- common.ProgressInput
```json
{
  "lesson": 1,
  "char_wpm": 20,
  "eff_wpm": 12,
  "accuracy": 0.95
}
```

## Error Response Shape
```json
{
  "error": "USERNAME_TAKEN",
  "message": "Username already exists",
  "status": 409
}
```

## Validation Error Shape (422)
(Note: current handlers return 400 for binding errors; included here as an example validation payload shape used by many APIs)
```json
{
  "error": "VALIDATION_FAILED",
  "message": "Field validation errors",
  "status": 422,
  "details": {
	"username": "must be 3-20 characters",
	"password": "must contain 8+ characters"
  }
}
```

## Pagination Pattern
All list endpoints support:
- `?page=1&limit=20` (default: page 1, limit 20)
- Response includes `{"data": [...], "total": 150, "page": 1, "pages": 8}`


