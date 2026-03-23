# API Documentation
*Generated: March 23, 2026*
*Base URL: http://localhost:8080*

## Endpoints Overview
| Method | Path                      | Description                        |
|--------|---------------------------|-------------------------------------|
| GET    | /v1/health                | Health check                       |
| POST   | /v1/auth/register         | Register new user                  |
| POST   | /v1/auth/login            | Login and get tokens               |
| POST   | /v1/auth/logout           | Logout and revoke refresh token     |
| POST   | /v1/auth/refresh          | Refresh access/refresh tokens       |
| POST   | /v1/auth/send-verification-email | Send verification code to current email |
| POST   | /v1/auth/verify-email     | Verify current email with OTP code  |
| GET    | /v1/settings/all          | Get all user settings               |
| GET    | /v1/settings/cw           | Get CW settings                     |
| GET    | /v1/settings/page         | Get page settings                   |
| POST   | /v1/settings/cw           | Update CW settings                  |
| POST   | /v1/settings/page         | Update page settings                |
| GET    | /v1/user/me               | Get current user info               |
| PUT    | /v1/user/callsign         | Update user call sign               |
| PUT    | /v1/user/email            | Update user email                   |
| PUT    | /v1/user/password         | Update user password                |
| GET    | /v1/cw/progress           | Get all progress                    |
| PUT    | /v1/cw/progress           | Add progress                        |
| GET    | /v1/hello                 | Authenticated hello                 |

---

## Authentication
Public endpoints: `/v1/health`, `/v1/auth/register`, `/v1/auth/login`, `/v1/auth/logout`, `/v1/auth/refresh`.

All other endpoints require Bearer JWT authentication.

---

## Auth Endpoints

### POST /v1/auth/register
**Registers a new user**

**Request Body:**
```json
{
  "username": "string (3-20 chars)",
  "email": "string (valid email)",
  "password": "string (min 8 chars)"
}
```
*Validation: username (required, 3-20 chars), email (required, valid), password (required, min 8 chars)*

**Response (200):**
```json
{
  "refresh_token": "...",
  "access_token": "..."
}
```

**Error Responses:**
| Status | Error Code                        | Message                          |
|--------|-----------------------------------|----------------------------------|
| 400    | `INVALID_REQUEST_BODY`            | "Invalid request body"           |
| 409    | `USERNAME_ALREADY_IN_USE`         | "Username already exists"        |
| 409    | `EMAIL_VERIFIED_BY_ANOTHER_ACCOUNT` | "This email is already verified by another account. Please change your email." |
| 500    | `DATABASE_FAILURE`                | "Database failure"               |
| 500    | `PASSWORD_HASH_FAILED`            | "Failed to hash password"        |
| 500    | `INTERNAL_SERVER_ERROR`           | "Failed to create user"          |

**Example cURL:**
```bash
curl -X POST http://localhost:8080/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test@example.com","password":"password123"}'
```

### POST /v1/auth/login
**Login and receive tokens**

**Request Body:**
```json
{
  "identifier": "username or email",
  "password": "string"
}
```

**Response (200):**
```json
{
  "refresh_token": "...",
  "access_token": "..."
}
```

**Error Responses:**
| Status | Error Code                | Message                |
|--------|---------------------------|------------------------|
| 400    | `INVALID_REQUEST_BODY`    | "Invalid request body" |
| 401    | `INVALID_CREDENTIALS`     | "Invalid credentials"  |
| 500    | `INTERNAL_SERVER_ERROR`   | "internal error"       |

**Example cURL:**
```bash
curl -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"identifier":"testuser","password":"password123"}'
```

### POST /v1/auth/logout
**Logout and revoke refresh token**

**Request Body:**
```json
{
  "refresh_token": "..."
}
```

**Response (200):**
```json
{
  "message": "Logged out"
}
```

**Error Responses:**
| Status | Error Code              | Message                  |
|--------|-------------------------|--------------------------|
| 400    | `INVALID_REQUEST_BODY`  | "Invalid request body"   |
| 401    | `INVALID_TOKEN`         | "Invalid refresh token"  |
| 500    | `INTERNAL_SERVER_ERROR` | "Failed to logout"       |

**Example cURL:**
```bash
curl -X POST http://localhost:8080/v1/auth/logout \
  -H "Content-Type: application/json" \
  -d '{"refresh_token":"..."}'
```

### POST /v1/auth/refresh
**Refresh access and refresh tokens**

**Request Body:**
```json
{
  "refresh_token": "..."
}
```

**Response (200):**
```json
{
  "refresh_token": "...",
  "access_token": "..."
}
```

**Error Responses:**
| Status | Error Code              | Message                        |
|--------|-------------------------|--------------------------------|
| 400    | `INVALID_REQUEST_BODY`  | "Invalid request body"         |
| 401    | `INVALID_TOKEN`         | "Invalid refresh token"        |
| 401    | `EXPIRED_TOKEN`         | "Refresh token expired"        |
| 500    | `INTERNAL_SERVER_ERROR` | "Failed to generate token"     |

**Example cURL:**
```bash
curl -X POST http://localhost:8080/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{"refresh_token":"..."}'
```

### POST /v1/auth/send-verification-email
**Send a verification OTP to the authenticated user's current email**

**Authentication:** Bearer JWT

**Request Body:**
None

**Response (200):**
```json
{
  "message": "Verification email sent"
}
```

**Error Responses:**
| Status | Error Code                 | Message                               |
|--------|----------------------------|---------------------------------------|
| 400    | `EMAIL_ALREADY_VERIFIED`   | "Email is already verified"          |
| 409    | `EMAIL_VERIFIED_BY_ANOTHER_ACCOUNT` | "This email is already verified by another account. Please change your email." |
| 429    | `VERIFICATION_RATE_LIMITED`| "Please wait before requesting another verification email" |
| 500    | `DATABASE_FAILURE`         | "Database failure"                   |
| 500    | `VERIFICATION_SEND_FAILED` | "Failed to send verification email"  |

**Example cURL:**
```bash
curl -X POST http://localhost:8080/v1/auth/send-verification-email \
  -H "Authorization: Bearer <token>"
```

### POST /v1/auth/verify-email
**Verify authenticated user's current email using OTP code**

**Authentication:** Bearer JWT

**Request Body:**
```json
{
  "code": "123456"
}
```

**Response (200):**
```json
{
  "message": "Email verified"
}
```

**Error Responses:**
| Status | Error Code                   | Message                       |
|--------|------------------------------|-------------------------------|
| 400    | `INVALID_REQUEST_BODY`       | "Invalid request body"       |
| 400    | `EMAIL_ALREADY_VERIFIED`     | "Email is already verified"  |
| 400    | `VERIFICATION_CODE_INVALID`  | "Invalid verification code"  |
| 400    | `VERIFICATION_CODE_EXPIRED`  | "Verification code expired"  |
| 409    | `EMAIL_VERIFIED_BY_ANOTHER_ACCOUNT` | "This email is already verified by another account. Please change your email." |
| 500    | `DATABASE_FAILURE`           | "Database failure"           |

**Example cURL:**
```bash
curl -X POST http://localhost:8080/v1/auth/verify-email \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"code":"123456"}'
```

---

## Settings Endpoints

### GET /v1/settings/all
**Get all user settings**

**Authentication:** Bearer JWT

**Response (200):**
```json
{
  "cw_settings": {
    "char_wpm": 20,
    "eff_wpm": 15,
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

**Error Responses:**
| Status | Error Code                | Message                  |
|--------|---------------------------|--------------------------|
| 500    | `SETTINGS_FETCH_FAILED`   | "Failed to get settings" |

**Example cURL:**
```bash
curl -X GET http://localhost:8080/v1/settings/all \
  -H "Authorization: Bearer <token>"
```

### GET /v1/settings/cw
**Get CW settings**

**Authentication:** Bearer JWT

**Response (200):**
```json
{
  "char_wpm": 20,
  "eff_wpm": 15,
  "freq": 600,
  "start_delay": 0.5
}
```

**Error Responses:**
| Status | Error Code                | Message                  |
|--------|---------------------------|--------------------------|
| 500    | `SETTINGS_FETCH_FAILED`   | "Failed to get settings" |

**Example cURL:**
```bash
curl -X GET http://localhost:8080/v1/settings/cw \
  -H "Authorization: Bearer <token>"
```

### GET /v1/settings/page
**Get page settings**

**Authentication:** Bearer JWT

**Response (200):**
```json
{
  "theme": "auto",
  "language": "auto",
  "cur_lesson": 1
}
```

**Error Responses:**
| Status | Error Code                | Message                  |
|--------|---------------------------|--------------------------|
| 500    | `SETTINGS_FETCH_FAILED`   | "Failed to get settings" |

**Example cURL:**
```bash
curl -X GET http://localhost:8080/v1/settings/page \
  -H "Authorization: Bearer <token>"
```

### POST /v1/settings/cw
**Update CW settings**

**Authentication:** Bearer JWT

**Request Body:**
```json
{
  "char_wpm": 20,
  "eff_wpm": 15,
  "freq": 600,
  "start_delay": 0.5
}
```
*Validation: char_wpm (5-50), eff_wpm (5-50), freq (300-2000), start_delay (0.0-10.0)*

**Response (200):**
```json
{
  "message": "Settings updated"
}
```

**Error Responses:**
| Status | Error Code                  | Message                    |
|--------|-----------------------------|----------------------------|
| 400    | `INVALID_REQUEST_BODY`      | "Invalid request body"     |
| 500    | `SETTINGS_UPDATE_FAILED`    | "Failed to update settings"|

**Example cURL:**
```bash
curl -X POST http://localhost:8080/v1/settings/cw \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"char_wpm":20,"eff_wpm":15,"freq":600,"start_delay":0.5}'
```

### POST /v1/settings/page
**Update page settings**

**Authentication:** Bearer JWT

**Request Body:**
```json
{
  "theme": "auto",
  "language": "auto",
  "cur_lesson": 1
}
```
*Validation: theme (auto|dark|light), language (required), cur_lesson (required)*

**Response (200):**
```json
{
  "message": "Settings updated"
}
```

**Error Responses:**
| Status | Error Code                  | Message                    |
|--------|-----------------------------|----------------------------|
| 400    | `INVALID_REQUEST_BODY`      | "Invalid request body"     |
| 500    | `SETTINGS_UPDATE_FAILED`    | "Failed to update settings"|

**Example cURL:**
```bash
curl -X POST http://localhost:8080/v1/settings/page \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"theme":"auto","language":"auto","cur_lesson":1}'
```

---

## User Endpoints

### GET /v1/user/me
**Get current user info**

**Authentication:** Bearer JWT

**Response (200):**
```json
{
  "call_sign": null,
  "username": "testuser",
  "email": "test@example.com",
  "email_verified": false,
  "created_at": "2026-03-21T18:30:00Z"
}
```

**Example cURL:**
```bash
curl -X GET http://localhost:8080/v1/user/me \
  -H "Authorization: Bearer <token>"
```

### PUT /v1/user/callsign
**Update user call sign**

**Authentication:** Bearer JWT

**Request Body:**
```json
{
  "call_sign": "K1ABC"
}
```

**Response (200):**
```json
{
  "message": "Call sign updated"
}
```

**Error Responses:**
| Status | Error Code                  | Message                    |
|--------|-----------------------------|----------------------------|
| 400    | `INVALID_REQUEST_BODY`      | "Invalid request body"     |
| 409    | `CALL_SIGN_ALREADY_IN_USE`  | "Call sign already in use" |
| 500    | `INTERNAL_SERVER_ERROR`     | "Failed to update call sign"|

**Example cURL:**
```bash
curl -X PUT http://localhost:8080/v1/user/callsign \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"call_sign":"K1ABC"}'
```

### PUT /v1/user/email
**Update user email**

**Authentication:** Bearer JWT

**Request Body:**
```json
{
  "email": "new@example.com"
}
```

**Response (200):**
```json
{
  "message": "Email updated"
}
```

**Error Responses:**
| Status | Error Code                  | Message                    |
|--------|-----------------------------|----------------------------|
| 400    | `INVALID_REQUEST_BODY`      | "Invalid request body"     |
| 400    | `EMAIL_UNCHANGED`           | "New email must be different from current email"|
| 409    | `EMAIL_VERIFIED_BY_ANOTHER_ACCOUNT` | "This email is already verified by another account. Please change your email." |
| 409    | `EMAIL_ALREADY_IN_USE`      | "Email already in use"     |
| 500    | `INTERNAL_SERVER_ERROR`     | "Failed to update email"   |

**Example cURL:**
```bash
curl -X PUT http://localhost:8080/v1/user/email \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"email":"new@example.com"}'
```

### PUT /v1/user/password
**Update user password**

**Authentication:** Bearer JWT

**Request Body:**
```json
{
  "old_password": "oldpass123",
  "new_password": "newpass456"
}
```

**Response (200):**
```json
{
  "message": "Password updated"
}
```

**Error Responses:**
| Status | Error Code                  | Message                    |
|--------|-----------------------------|----------------------------|
| 400    | `INVALID_REQUEST_BODY`      | "Invalid request body"     |
| 401    | `INVALID_CREDENTIALS`       | "Invalid credentials"      |
| 500    | `PASSWORD_HASH_FAILED`      | "Failed to hash password"  |
| 500    | `INTERNAL_SERVER_ERROR`     | "Failed to update password"|

**Example cURL:**
```bash
curl -X PUT http://localhost:8080/v1/user/password \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"old_password":"oldpass123","new_password":"newpass456"}'
```

---

## Progress Endpoints

### GET /v1/cw/progress
**Get all progress for current user**

**Authentication:** Bearer JWT

**Response (200):**
```json
{
  "data": [
    {
      "lesson": 1,
      "char_wpm": 20,
      "eff_wpm": 15,
      "accuracy": 0.98,
      "created_at": "2026-03-21T18:30:00Z",
      "client_created_at": "2026-03-21T18:29:50Z"
    }
  ]
}
```

**Error Responses:**
| Status | Error Code                  | Message                    |
|--------|-----------------------------|----------------------------|
| 500    | `PROGRESS_QUERY_FAILED`     | "failed to query progress" |

**Example cURL:**
```bash
curl -X GET http://localhost:8080/v1/cw/progress \
  -H "Authorization: Bearer <token>"
```

### PUT /v1/cw/progress
**Add a new progress record**

**Authentication:** Bearer JWT

**Request Body:**
```json
{
  "lesson": 1,
  "char_wpm": 20,
  "eff_wpm": 15,
  "accuracy": 0.98,
  "client_created_at": "2026-03-21T18:29:50Z"
}
```
*Validation: lesson (required), char_wpm (5-50), eff_wpm (5-50), accuracy (0.0-1.0), client_created_at (optional RFC3339 timestamp)*

**Response (201):**
```json
{
  "message": "Progress Created"
}
```

**Error Responses:**
| Status | Error Code                  | Message                    |
|--------|-----------------------------|----------------------------|
| 400    | `INVALID_REQUEST_BODY`      | "Invalid request body"     |
| 500    | `PROGRESS_CREATE_FAILED`    | "Failed to create progress"|

**Example cURL:**
```bash
curl -X PUT http://localhost:8080/v1/cw/progress \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"lesson":1,"char_wpm":20,"eff_wpm":15,"accuracy":0.98,"client_created_at":"2026-03-21T18:29:50Z"}'
```

---

## Miscellaneous

### GET /v1/hello
**Authenticated hello endpoint**

**Authentication:** Bearer JWT

**Response (200):**
```json
{
  "message": "Hello, authenticated user {username}!"
}
```

**Example cURL:**
```bash
curl -X GET http://localhost:8080/v1/hello \
  -H "Authorization: Bearer <token>"
```

### GET /v1/health
**Health check endpoint**

**Response (200):**
```json
{
  "status": "healthy",
  "timestamp": 1711132200
}
```

**Example cURL:**
```bash
curl -X GET http://localhost:8080/v1/health
```

---

## Error Response Shape
```json
{
  "code": "ERROR_CODE",
  "error": "Error message"
}
```

## Validation Error Shape (example)
```json
{
  "code": "INVALID_REQUEST_BODY",
  "error": "Invalid request body",
  "details": {
    "field": "validation error message"
  }
}
```

## Error Codes
| Code                           | Description                       |
|--------------------------------|-----------------------------------|
| INVALID_REQUEST_BODY           | Malformed or invalid request body |
| INTERNAL_SERVER_ERROR          | Unexpected server error           |
| DATABASE_FAILURE               | Database operation failed         |
| INVALID_CREDENTIALS            | Invalid username/email or password|
| CONFLICT                       | Resource conflict                 |
| INVALID_TOKEN                  | Invalid or expired token          |
| EXPIRED_TOKEN                  | Token expired                     |
| AUTH_HEADER_REQUIRED           | Missing Authorization header      |
| INVALID_AUTH_HEADER_FORMAT     | Bad Authorization header format   |
| USER_NOT_FOUND                 | User not found                    |
| SETTINGS_FETCH_FAILED          | Failed to fetch settings          |
| SETTINGS_UPDATE_FAILED         | Failed to update settings         |
| PROGRESS_QUERY_FAILED          | Failed to query progress          |
| PROGRESS_CREATE_FAILED         | Failed to create progress         |
| PASSWORD_HASH_FAILED           | Password hashing failed           |
| TOKEN_ISSUE_FAILED             | Token generation failed           |
| EMAIL_ALREADY_IN_USE           | Email already in use              |
| EMAIL_VERIFIED_BY_ANOTHER_ACCOUNT | Email already verified by another account |
| USERNAME_ALREADY_IN_USE        | Username already in use           |
| EMAIL_UNCHANGED                | Email is unchanged                |
| EMAIL_ALREADY_VERIFIED         | Email already verified            |
| VERIFICATION_CODE_INVALID      | Verification code is invalid      |
| VERIFICATION_CODE_EXPIRED      | Verification code is expired      |
| VERIFICATION_SEND_FAILED       | Failed to send verification email |
| VERIFICATION_RATE_LIMITED      | Verification email request throttled |
| CALL_SIGN_ALREADY_IN_USE       | Call sign already in use          |

---

