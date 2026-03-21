# API Documentation
*Generated: 2026-03-21*
*Base URL: http://localhost:8080*

## Endpoints Overview
| Method | Path                      | Description                                   |
|--------|---------------------------|-----------------------------------------------|
| GET    | /v1/health                | Health check                                  |
| POST   | /v1/auth/register         | Register a new user / issue tokens            |
| POST   | /v1/auth/login            | Login with username/email and password        |
| POST   | /v1/auth/refresh          | Refresh access token using refresh token     |
| GET    | /v1/settings/all          | Get both CW and Page settings for user       |
| GET    | /v1/settings/cw           | Get CW (typing) settings                     |
| GET    | /v1/settings/page         | Get page settings                            |
| POST   | /v1/settings/cw           | Update CW settings                           |
| POST   | /v1/settings/page         | Update page settings                         |
| GET    | /v1/user/me               | Get authenticated user's info                |
| PUT    | /v1/user/email            | Update authenticated user's email            |
| PUT    | /v1/user/password         | Update authenticated user's password         |
| GET    | /v1/cw/progress           | List user's progress entries                 |
| PUT    | /v1/cw/progress           | Create a new progress entry                  |
| GET    | /v1/hello                 | Hello message for authenticated user         |

## Authentication Endpoints

### POST /v1/auth/register
**Registers a new user and returns an access + refresh token pair**

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
*Validation: username (required, custom `username` validator), email (required, valid email, max 254), password (required, min 8 chars, max 256)*

**Response (200)**:
```json
{
  "refresh_token": "<raw_refresh_token>",
  "access_token": "<jwt_access_token>"
}
```

**Error Responses**:
| Status | Error Code                          | Message                                      |
|--------|-------------------------------------|----------------------------------------------|
| 400    | `INVALID_REQUEST_BODY`              | "Invalid request body"                      |
| 409    | `USERNAME_EMAIL_ALREADY_IN_USE`     | "Username and email already exists"        |
| 409    | `USERNAME_ALREADY_IN_USE`           | "Username already exists"                  |
| 409    | `EMAIL_ALREADY_IN_USE`              | "Email already exists"                     |
| 409    | `CONFLICT`                          | "Registration conflict, please try again"  |
| 500    | `DATABASE_FAILURE`                  | "Database failure"                          |
| 500    | `PASSWORD_HASH_FAILED`              | "Failed to hash password"                   |
| 500    | `TOKEN_ISSUE_FAILED`                | "Failed to issue token, try to login."     |

**Error Response Shape** (actual service shape):
```json
{
  "code": "USERNAME_ALREADY_IN_USE",
  "error": "Username already exists"
}
```

**Example cURL**:
```bash
curl -X POST http://localhost:8080/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
	"username": "alan_yeung",
	"email": "alan@example.com",
	"password": "securepass123"
  }'
```

---

### POST /v1/auth/login
**Log in with username or email and password**

**Authentication**: None

**Request Body**:
```json
{
  "identifier": "alan_yeung_or_email@example.com",
  "password": "yourpassword"
}
```
*Validation: identifier (required), password (required)*

**Response (200)**:
```json
{
  "refresh_token": "<raw_refresh_token>",
  "access_token": "<jwt_access_token>"
}
```

**Error Responses**:
| Status | Error Code              | Message                  |
|--------|-------------------------|--------------------------|
| 400    | `INVALID_REQUEST_BODY`  | "Invalid request body"  |
| 401    | `INVALID_CREDENTIALS`   | "Invalid credentials"   |
| 500    | `INTERNAL_SERVER_ERROR` | "internal error"        |
| 500    | `TOKEN_ISSUE_FAILED`    | "Failed to issue token" |

**Example cURL**:
```bash
curl -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
	"identifier": "alan_yeung",
	"password": "securepass123"
  }'
```

---

### POST /v1/auth/refresh
**Exchange a refresh token for a new access token and rotated refresh token**

**Authentication**: None (uses refresh token in request body)

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
  "access_token": "<new_access_token>"
}
```

**Error Responses**:
| Status | Error Code              | Message                      |
|--------|-------------------------|------------------------------|
| 400    | `INVALID_REQUEST_BODY`  | "Invalid request body"      |
| 401    | `INVALID_TOKEN`         | "Invalid refresh token"     |
| 401    | `EXPIRED_TOKEN`         | "Refresh token expired"     |
| 500    | `INTERNAL_SERVER_ERROR` | "Failed to generate access token" |

**Example cURL**:
```bash
curl -X POST http://localhost:8080/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{"refresh_token":"<raw_refresh_token>"}'
```

---

## Settings

All /v1/settings/* endpoints require Authorization: Bearer <access_token> (JWT) and the user must be authenticated.

### GET /v1/settings/all
**Fetch both CW and Page settings for the authenticated user**

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
| Status | Error Code                    | Message                  |
|--------|-------------------------------|--------------------------|
| 401    | `AUTH_HEADER_REQUIRED` / `INVALID_TOKEN` | "Authorization header is required" / "Invalid token" |
| 401    | `USER_NOT_FOUND`              | "User not found"        |
| 500    | `SETTINGS_FETCH_FAILED`       | "Failed to get settings"|

**Example cURL**:
```bash
curl -X GET http://localhost:8080/v1/settings/all \
  -H "Authorization: Bearer <access_token>"
```

---

### GET /v1/settings/cw
**Fetch user's CW settings**

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

**Errors**:
| Status | Error Code                | Message                  |
|--------|---------------------------|--------------------------|
| 401    | `AUTH_HEADER_REQUIRED` / `INVALID_TOKEN` | "Authorization required / invalid token" |
| 500    | `SETTINGS_FETCH_FAILED`   | "Failed to get settings"|

**Example cURL**:
```bash
curl -X GET http://localhost:8080/v1/settings/cw \
  -H "Authorization: Bearer <access_token>"
```

---

### GET /v1/settings/page
**Fetch user's page settings**

**Authentication**: Bearer JWT

**Response (200)**:
```json
{
  "theme": "auto",
  "language": "auto",
  "cur_lesson": 0
}
```

**Errors**:
| Status | Error Code                | Message                  |
|--------|---------------------------|--------------------------|
| 401    | `AUTH_HEADER_REQUIRED` / `INVALID_TOKEN` | "Authorization required / invalid token" |
| 500    | `SETTINGS_FETCH_FAILED`   | "Failed to get settings"|

**Example cURL**:
```bash
curl -X GET http://localhost:8080/v1/settings/page \
  -H "Authorization: Bearer <access_token>"
```

---

### POST /v1/settings/cw
**Create or update CW settings for authenticated user**

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
*Validation: char_wpm (required, 5-50), eff_wpm (required, 5-50), freq (required, 300-2000), start_delay (required, 0.0-10.0)*

**Response (200)**:
```json
{
  "message": "Settings updated"
}
```

**Errors**:
| Status | Error Code                    | Message                  |
|--------|-------------------------------|--------------------------|
| 400    | `INVALID_REQUEST_BODY`        | "Invalid request body"  |
| 401    | `INVALID_TOKEN` / `AUTH_HEADER_REQUIRED` | "Invalid or missing token" |
| 500    | `SETTINGS_UPDATE_FAILED`      | "Failed to update settings"|

**Example cURL**:
```bash
curl -X POST http://localhost:8080/v1/settings/cw \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{"char_wpm":25,"eff_wpm":12,"freq":700,"start_delay":0.5}'
```

---

### POST /v1/settings/page
**Create or update page settings for authenticated user**

**Authentication**: Bearer JWT

**Request Body**:
```json
{
  "theme": "dark",
  "language": "en",
  "cur_lesson": 3
}
```
*Validation: theme (required, oneof: auto,dark,light), language (required), cur_lesson (required)*

**Response (200)**:
```json
{
  "message": "Settings updated"
}
```

**Errors**:
| Status | Error Code                    | Message                  |
|--------|-------------------------------|--------------------------|
| 400    | `INVALID_REQUEST_BODY`        | "Invalid request body"  |
| 401    | `INVALID_TOKEN` / `AUTH_HEADER_REQUIRED` | "Invalid or missing token" |
| 500    | `SETTINGS_UPDATE_FAILED`      | "Failed to update settings"|

**Example cURL**:
```bash
curl -X POST http://localhost:8080/v1/settings/page \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{"theme":"dark","language":"en","cur_lesson":3}'
```

---

## User Management

All `/v1/user/*` endpoints require Authorization: Bearer <access_token> (JWT).

### GET /v1/user/me
**Retrieves the authenticated user's profile information**

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

**Errors**:
| Status | Error Code            | Message              |
|--------|-----------------------|----------------------|
| 401    | `AUTH_HEADER_REQUIRED` / `INVALID_TOKEN` | "Missing or invalid token" |
| 401    | `USER_NOT_FOUND`      | "User not found"    |

**Example cURL**:
```bash
curl -X GET http://localhost:8080/v1/user/me \
  -H "Authorization: Bearer <access_token>"
```

---

### PUT /v1/user/email
**Update authenticated user's email**

**Authentication**: Bearer JWT

**Request Body**:
```json
{
  "email": "new_email@example.com"
}
```
*Validation: email (required, valid email, max 254)*

**Response (200)**:
```json
{
  "message": "Email updated"
}
```

**Error Responses**:
| Status | Error Code                    | Message                        |
|--------|-------------------------------|--------------------------------|
| 400    | `INVALID_REQUEST_BODY`        | "Invalid request body"        |
| 400    | `EMAIL_UNCHANGED`             | "New email must be different from current email" |
| 409    | `EMAIL_ALREADY_IN_USE`        | "Email already in use"        |
| 500    | `INTERNAL_SERVER_ERROR`       | "Failed to update email"      |

**Example cURL**:
```bash
curl -X PUT http://localhost:8080/v1/user/email \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{"email":"new_email@example.com"}'
```

---

### PUT /v1/user/password
**Update authenticated user's password**

**Authentication**: Bearer JWT

**Request Body**:
```json
{
  "old_password": "currentpass",
  "new_password": "newsecurepass"
}
```
*Validation: old_password (required, min 8), new_password (required, min 8)*

**Response (200)**:
```json
{
  "message": "Password updated"
}
```

**Error Responses**:
| Status | Error Code                    | Message                    |
|--------|-------------------------------|----------------------------|
| 400    | `INVALID_REQUEST_BODY`        | "Invalid request body"    |
| 401    | `INVALID_CREDENTIALS`         | "Invalid credentials"     |
| 500    | `PASSWORD_HASH_FAILED`        | "Failed to hash password" |
| 500    | `INTERNAL_SERVER_ERROR`       | "Failed to update password"|

**Example cURL**:
```bash
curl -X PUT http://localhost:8080/v1/user/password \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{"old_password":"currentpass","new_password":"newsecurepass"}'
```

---

## Progress

All `/v1/cw/*` progress endpoints require Authorization: Bearer <access_token> (JWT).

### GET /v1/cw/progress
**List all progress entries for the authenticated user**

**Authentication**: Bearer JWT

**Response (200)**:
```json
{
  "data": [
	{
	  "lesson": "Lesson 1",
	  "char_wpm": 25,
	  "eff_wpm": 12,
	  "accuracy": 0.95,
	  "created_at": "2026-03-21T18:30:00Z"
	}
  ]
}
```

**Error Responses**:
| Status | Error Code                    | Message                  |
|--------|-------------------------------|--------------------------|
| 401    | `INVALID_TOKEN` / `AUTH_HEADER_REQUIRED` | "Missing or invalid token" |
| 500    | `PROGRESS_QUERY_FAILED`        | "failed to query progress" |

**Example cURL**:
```bash
curl -X GET http://localhost:8080/v1/cw/progress \
  -H "Authorization: Bearer <access_token>"
```

---

### PUT /v1/cw/progress
**Create a new progress entry for the authenticated user**

**Authentication**: Bearer JWT

**Request Body**:
```json
{
  "lesson": 1,
  "char_wpm": 25,
  "eff_wpm": 12,
  "accuracy": 0.95
}
```
*Validation: lesson (required), char_wpm (required, 5-50), eff_wpm (required, 5-50), accuracy (required, 0.0-1.0)*

**Response (201)**:
```json
{
  "message": "Progress Created"
}
```

**Error Responses**:
| Status | Error Code                    | Message                  |
|--------|-------------------------------|--------------------------|
| 400    | `INVALID_REQUEST_BODY`        | "Invalid request body"  |
| 401    | `INVALID_TOKEN`               | "Missing or invalid token" |
| 500    | `PROGRESS_CREATE_FAILED`      | "Failed to create progress" |

**Example cURL**:
```bash
curl -X PUT http://localhost:8080/v1/cw/progress \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{"lesson":1,"char_wpm":25,"eff_wpm":12,"accuracy":0.95}'
```

---

## Misc

### GET /v1/health
**Simple health check**

**Authentication**: None

**Response (200)**:
```json
{
  "status": "healthy",
  "timestamp": 1700000000
}
```

### GET /v1/hello
**Return a greeting for authenticated user**

**Authentication**: Bearer JWT

**Response (200)**:
```json
{
  "message": "Hello, authenticated user {username}!"
}
```

## Structs / Schemas

Key request/response structs (from `common` package):

- RegisterInput: {username, email, password}
- LoginInput: {identifier, password}
- RefreshInput: {refresh_token}
- UpdateCallSignInput: {call_sign}
- UpdateEmailInput: {email}
- UpdatePasswordInput: {old_password, new_password}
- CWSettingsInput: {char_wpm, eff_wpm, freq, start_delay}
- PageSettingsInput: {theme, language, cur_lesson}
- ProgressInput: {lesson, char_wpm, eff_wpm, accuracy}

Responses:
- AuthTokenPairResponse: {refresh_token, access_token}
- UserInfoResponse: {call_sign, username, email, email_verified, created_at}
- CWSettingsResponse: {char_wpm, eff_wpm, freq, start_delay}
- PageSettingsResponse: {theme, language, cur_lesson}
- ProgressResponse: {lesson, char_wpm, eff_wpm, accuracy, created_at}
- MessageResponse: {message}
- ErrorResponse: {code, error}

## Error Codes

All string error codes defined in the project (from `common/error.go`):

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
USERNAME_ALREADY_IN_USE
USERNAME_EMAIL_ALREADY_IN_USE
EMAIL_UNCHANGED
CALL_SIGN_ALREADY_IN_USE
```

## Pagination Pattern
All list endpoints in this project currently return a simple array in `data` (for example `/v1/cw/progress`). The repository contains no global pagination middleware in the code inspected. If pagination is added later it will follow the common pattern: `?page=1&limit=20` and responses would include `{ "data": [...], "total": <n>, "page": <p>, "pages": <m> }`.


