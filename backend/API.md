# API Documentation
*Generated: 2026-03-21*
*Base URL: http://localhost:8080*

## Endpoints Overview
| Method | Path                      | Description                             |
|--------|---------------------------|-----------------------------------------|
| GET    | /v1/health                | Health check                             |
| POST   | /v1/auth/register         | Register new user                        |
| POST   | /v1/auth/login            | Login and receive tokens                 |
| POST   | /v1/auth/refresh          | Refresh access/refresh token pair        |
| GET    | /v1/settings/all          | Get all user settings (cw + page)        |
| GET    | /v1/settings/cw           | Get user's CW settings                   |
| GET    | /v1/settings/page         | Get user's page settings                 |
| POST   | /v1/settings/cw           | Update CW settings                       |
| POST   | /v1/settings/page         | Update Page settings                     |
| GET    | /v1/user/me               | Get authenticated user's info            |
| PUT    | /v1/user/email            | Update authenticated user's email        |
| PUT    | /v1/user/password         | Update authenticated user's password     |
| GET    | /v1/cw/progress           | List user's lesson progress              |
| PUT    | /v1/cw/progress           | Add new progress entry                   |
| GET    | /v1/hello                 | Simple authenticated hello endpoint      |

## Authentication / Middleware
- Protected endpoints (settings, user, cw progress, /v1/hello) require two middlewares: `AuthRequired()` and `LoadUser(db)`.
- `AuthRequired()` expects an `Authorization: Bearer <token>` header and validates a JWT (signed with HMAC using app JWT secret). Possible auth errors returned by middleware: `AUTH_HEADER_REQUIRED`, `INVALID_AUTH_HEADER_FORMAT`, `INVALID_TOKEN` (all return HTTP 401).
- `LoadUser(db)` loads the user from DB using the subject from token and returns `USER_NOT_FOUND` (HTTP 401) when the user cannot be found.

## Schemas (structures with json tags)

Request bodies and responses use the following Go structs (rendered as example JSON):

- RegisterInput
```json
{
  "username": "alan_yeung",
  "email": "alan@example.com",
  "password": "min8chars123"
}
```
Validation: username uses custom `username` validator (regex ^[a-zA-Z0-9][a-zA-Z0-9_-]{1,14}[a-zA-Z0-9]$) — length 3-16, may include letters, digits, underscores and hyphens but cannot start/end with underscore/hyphen. Email must be valid and <= 254 chars. Password min 8, max 256.

- LoginInput
```json
{
  "identifier": "alan_yeung or alan@example.com",
  "password": "min8chars123"
}
```

- RefreshInput
```json
{
  "refresh_token": "<raw_refresh_token>"
}
```

- CWSettingsInput
```json
{
  "char_wpm": 20,
  "eff_wpm": 10,
  "freq": 600,
  "start_delay": 0.5
}
```
Validation: char_wpm/eﬀ_wpm 5-50, freq 300-2000, start_delay 0.0-10.0

- PageSettingsInput
```json
{
  "theme": "auto",
  "language": "auto",
  "cur_lesson": 1
}
```
Validation: theme one of [auto, dark, light]

- ProgressInput
```json
{
  "lesson": 1,
  "char_wpm": 20,
  "eff_wpm": 10,
  "accuracy": 0.95
}
```

- AuthTokenPairResponse (success response for login/register/refresh)
```json
{
  "refresh_token": "<raw_refresh_token>",
  "access_token": "<jwt_access_token>"
}
```

- UserInfoResponse
```json
{
  "call_sign": null,
  "username": "alan_yeung",
  "email": "alan@example.com",
  "email_verified": false,
  "created_at": "2026-03-21T18:30:00Z"
}
```

- CWSettingsResponse
```json
{
  "char_wpm": 20,
  "eff_wpm": 10,
  "freq": 600,
  "start_delay": 0.5
}
```

- PageSettingsResponse
```json
{
  "theme": "auto",
  "language": "auto",
  "cur_lesson": 1
}
```

- ProgressResponse (list items)
```json
{
  "lesson": 1,
  "char_wpm": 20,
  "eff_wpm": 10,
  "accuracy": 0.95,
  "created_at": "2026-03-21T18:30:00Z"
}
```

## Error Codes (from common/error.go)
All string error codes defined in the project:

```
INVALID_REQUEST_BODY
INTERNAL_SERVER_ERROR
DATABASE_FAILURE
INVALID_CREDENTIALS
CONFLICT
INVALID_TOKEN
EXPIRED_TOKEN
AUTH_HEADER_REQUIRED
INVALID_AUTH_HEADER_FORMAT
USER_NOT_FOUND
SETTINGS_FETCH_FAILED
SETTINGS_UPDATE_FAILED
PROGRESS_QUERY_FAILED
PROGRESS_CREATE_FAILED
PASSWORD_HASH_FAILED
TOKEN_ISSUE_FAILED
EMAIL_ALREADY_IN_USE
EMAIL_UNCHANGED
```

Error response shape (used by handlers):
```json
{
  "code": "USERNAME_TAKEN_OR_OTHER_CODE",
  "error": "Human readable error message"
}
```

Validation error shape (422-like semantic used by project on validation):
```json
{
  "error": "VALIDATION_FAILED",
  "message": "Field validation errors",
  "status": 422,
  "details": {
	"username": "must be 3-16 characters",
	"password": "must contain 8+ characters"
  }
}
```

## Endpoints

### POST /v1/auth/register
**Registers a new user**

**Authentication**: None

**Path Parameters**: None

**Query Parameters**: None

**Request Body**:
```json
{
  "username": "alan_yeung",
  "email": "alan@example.com",
  "password": "securepass123"
}
```
*Validation: username (regex, 3-16 chars), email (valid format, max 254), password (min 8 chars)*

**Response (200)**:
```json
{
  "refresh_token": "<raw_refresh_token>",
  "access_token": "<jwt_access_token>"
}
```

**Error Responses**:
| Status | Error Code              | Message                          |
|--------|-------------------------|----------------------------------|
| 400    | `INVALID_REQUEST_BODY`  | "Invalid request body"         |
| 409    | `CONFLICT`              | "Username and/or email exists" |
| 500    | `DATABASE_FAILURE`      | "Database failure"             |
| 500    | `PASSWORD_HASH_FAILED`  | "Failed to hash password"      |
| 500    | `TOKEN_ISSUE_FAILED`    | "Failed to issue token"        |

**Example cURL**:
```bash
curl -X POST http://localhost:8080/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"alan_yeung","email":"alan@example.com","password":"securepass123"}'
```

### POST /v1/auth/login
**Logs in a user and returns an access + refresh token pair**

**Authentication**: None

**Request Body**:
```json
{
  "identifier": "alan_yeung",
  "password": "securepass123"
}
```

**Response (200)**:
```json
{
  "refresh_token": "<raw_refresh_token>",
  "access_token": "<jwt_access_token>"
}
```

**Error Responses**:
| Status | Error Code              | Message                          |
|--------|-------------------------|----------------------------------|
| 400    | `INVALID_REQUEST_BODY`  | "Invalid request body"         |
| 401    | `INVALID_CREDENTIALS`   | "Invalid credentials"          |
| 500    | `INTERNAL_SERVER_ERROR` | "internal error"               |
| 500    | `TOKEN_ISSUE_FAILED`    | "Failed to issue token"        |

**Example cURL**:
```bash
curl -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"identifier":"alan_yeung","password":"securepass123"}'
```

### POST /v1/auth/refresh
**Exchanges a refresh token for a new pair**

**Authentication**: None (sends refresh token in body)

**Request Body**:
```json
{
  "refresh_token": "<raw_refresh_token>"
}
```

**Response (200)**:
```json
{
  "refresh_token": "<new_raw_refresh_token>",
  "access_token": "<jwt_access_token>"
}
```

**Error Responses**:
| Status | Error Code           | Message                  |
|--------|----------------------|--------------------------|
| 400    | `INVALID_REQUEST_BODY` | "Invalid request body" |
| 401    | `INVALID_TOKEN`      | "Invalid refresh token" |
| 401    | `EXPIRED_TOKEN`      | "Refresh token expired" |
| 500    | `INTERNAL_SERVER_ERROR` | "Failed to generate token" |

**Example cURL**:
```bash
curl -X POST http://localhost:8080/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{"refresh_token":"<raw_refresh_token>"}'
```

### GET /v1/health
**Health check**

**Authentication**: None

**Response (200)**:
```json
{
  "status": "healthy",
  "timestamp": 1679500000
}
```

**Example cURL**:
```bash
curl http://localhost:8080/v1/health
```

### GET /v1/settings/all
**Returns both CW and Page settings for the authenticated user**

**Authentication**: Bearer JWT

**Path Parameters**: None

**Query Parameters**: None

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
	"cur_lesson": 1
  }
}
```

**Error Responses**:
| Status | Error Code                  | Message                      |
|--------|-----------------------------|------------------------------|
| 401    | `AUTH_HEADER_REQUIRED` / `INVALID_TOKEN` | "Authorization required or invalid token" |
| 401    | `USER_NOT_FOUND`            | "User not found"            |
| 500    | `SETTINGS_FETCH_FAILED`     | "Failed to get settings"    |

**Example cURL**:
```bash
curl http://localhost:8080/v1/settings/all \
  -H "Authorization: Bearer <access_token>"
```

### GET /v1/settings/cw
**Get CW (typing) settings for authenticated user**

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
| Status | Error Code              | Message                      |
|--------|-------------------------|------------------------------|
| 401    | `AUTH_HEADER_REQUIRED` / `INVALID_TOKEN` | "Authorization required or invalid token" |
| 500    | `SETTINGS_FETCH_FAILED` | "Failed to get settings"    |

**Example cURL**:
```bash
curl http://localhost:8080/v1/settings/cw \
  -H "Authorization: Bearer <access_token>"
```

### GET /v1/settings/page
**Get page settings for authenticated user**

**Authentication**: Bearer JWT

**Response (200)**:
```json
{
  "theme": "auto",
  "language": "auto",
  "cur_lesson": 1
}
```

**Error Responses**:
| Status | Error Code              | Message                      |
|--------|-------------------------|------------------------------|
| 401    | `AUTH_HEADER_REQUIRED` / `INVALID_TOKEN` | "Authorization required or invalid token" |
| 500    | `SETTINGS_FETCH_FAILED` | "Failed to get settings"    |

**Example cURL**:
```bash
curl http://localhost:8080/v1/settings/page \
  -H "Authorization: Bearer <access_token>"
```

### POST /v1/settings/cw
**Update or create CW settings for authenticated user**

**Authentication**: Bearer JWT

**Request Body**: (see CWSettingsInput above)

**Response (200)**:
```json
{ "message": "Settings updated" }
```

**Error Responses**:
| Status | Error Code                  | Message                      |
|--------|-----------------------------|------------------------------|
| 400    | `INVALID_REQUEST_BODY`      | "Invalid request body"      |
| 401    | `AUTH_HEADER_REQUIRED` / `INVALID_TOKEN` | "Authorization required or invalid token" |
| 500    | `SETTINGS_UPDATE_FAILED`    | "Failed to update settings" |

**Example cURL**:
```bash
curl -X POST http://localhost:8080/v1/settings/cw \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{"char_wpm":20,"eff_wpm":10,"freq":600,"start_delay":0.5}'
```

### POST /v1/settings/page
**Update or create page settings for authenticated user**

**Authentication**: Bearer JWT

**Request Body**: (see PageSettingsInput above)

**Response (200)**:
```json
{ "message": "Settings updated" }
```

**Error Responses**:
| Status | Error Code                  | Message                      |
|--------|-----------------------------|------------------------------|
| 400    | `INVALID_REQUEST_BODY`      | "Invalid request body"      |
| 401    | `AUTH_HEADER_REQUIRED` / `INVALID_TOKEN` | "Authorization required or invalid token" |
| 500    | `SETTINGS_UPDATE_FAILED`    | "Failed to update settings" |

**Example cURL**:
```bash
curl -X POST http://localhost:8080/v1/settings/page \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{"theme":"auto","language":"auto","cur_lesson":1}'
```

### GET /v1/user/me
**Retrieves the authenticated user's profile**

**Authentication**: Bearer JWT

**Response (200)**: see `UserInfoResponse` above

**Error Responses**:
| Status | Error Code              | Message                      |
|--------|-------------------------|------------------------------|
| 401    | `AUTH_HEADER_REQUIRED` / `INVALID_TOKEN` | "Authorization required or invalid token" |
| 401    | `USER_NOT_FOUND`        | "User not found"            |

**Example cURL**:
```bash
curl http://localhost:8080/v1/user/me \
  -H "Authorization: Bearer <access_token>"
```

### PUT /v1/user/email
**Updates authenticated user's email**

**Authentication**: Bearer JWT

**Request Body**:
```json
{
  "email": "new@example.com"
}
```

**Response (200)**:
```json
{ "message": "Email updated" }
```

**Error Responses**:
| Status | Error Code                  | Message                                |
|--------|-----------------------------|----------------------------------------|
| 400    | `INVALID_REQUEST_BODY`      | "Invalid request body"               |
| 400    | `EMAIL_UNCHANGED`           | "New email must be different"        |
| 409    | `EMAIL_ALREADY_IN_USE`      | "Email already in use"               |
| 500    | `INTERNAL_SERVER_ERROR`     | "Failed to update email"             |

**Example cURL**:
```bash
curl -X PUT http://localhost:8080/v1/user/email \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{"email":"new@example.com"}'
```

### PUT /v1/user/password
**Updates authenticated user's password**

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
{ "message": "Password updated" }
```

**Error Responses**:
| Status | Error Code                  | Message                                |
|--------|-----------------------------|----------------------------------------|
| 400    | `INVALID_REQUEST_BODY`      | "Invalid request body"               |
| 401    | `INVALID_CREDENTIALS`       | "Invalid credentials"                |
| 500    | `PASSWORD_HASH_FAILED`      | "Failed to hash password"            |
| 500    | `INTERNAL_SERVER_ERROR`     | "Failed to update password"          |

**Example cURL**:
```bash
curl -X PUT http://localhost:8080/v1/user/password \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{"old_password":"current","new_password":"newsecure"}'
```

### GET /v1/cw/progress
**Lists authenticated user's progress entries**

**Authentication**: Bearer JWT

**Response (200)**:
```json
{ "data": [ /* array of ProgressResponse */ ] }
```

**Error Responses**:
| Status | Error Code                  | Message                      |
|--------|-----------------------------|------------------------------|
| 401    | `AUTH_HEADER_REQUIRED` / `INVALID_TOKEN` | "Authorization required or invalid token" |
| 500    | `PROGRESS_QUERY_FAILED`     | "failed to query progress"  |

**Example cURL**:
```bash
curl http://localhost:8080/v1/cw/progress \
  -H "Authorization: Bearer <access_token>"
```

### PUT /v1/cw/progress
**Adds a new progress entry for authenticated user**

**Authentication**: Bearer JWT

**Request Body**: see `ProgressInput` above

**Response (201)**:
```json
{ "message": "Progress Created" }
```

**Error Responses**:
| Status | Error Code                  | Message                      |
|--------|-----------------------------|------------------------------|
| 400    | `INVALID_REQUEST_BODY`      | "Invalid request body"      |
| 500    | `PROGRESS_CREATE_FAILED`    | "Failed to create progress" |

**Example cURL**:
```bash
curl -X PUT http://localhost:8080/v1/cw/progress \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{"lesson":1,"char_wpm":20,"eff_wpm":10,"accuracy":0.95}'
```

### GET /v1/hello
**Simple authenticated greeting that returns the username**

**Authentication**: Bearer JWT

**Response (200)**:
```json
{ "message": "Hello, authenticated user {username}!" }
```

**Error Responses**:
| Status | Error Code              | Message                      |
|--------|-------------------------|------------------------------|
| 401    | `AUTH_HEADER_REQUIRED` / `INVALID_TOKEN` | "Authorization required or invalid token" |

**Example cURL**:
```bash
curl http://localhost:8080/v1/hello \
  -H "Authorization: Bearer <access_token>"
```

## Pagination Pattern
All list endpoints support:
- `?page=1&limit=20` (default: page 1, limit 20)
- Response includes `{"data": [...], "total": 150, "page": 1, "pages": 8}`


