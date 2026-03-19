# API Documentation

**Base URL:** `https://api.example.com`
> Last updated: March 20, 2026

## Table of Contents
- [Auth](#auth)
	- [POST `/api/v1/auth/register`](#post-apiv1authregister)
	- [POST `/api/v1/auth/login`](#post-apiv1authlogin)
	- [POST `/api/v1/auth/refresh`](#post-apiv1authrefresh)
- [Health](#health)
	- [GET `/api/v1/health`](#get-apiv1health)
- [Settings](#settings)
	- [GET `/api/v1/settings/`](#get-apiv1settings)
	- [GET `/api/v1/settings/cw`](#get-apiv1settingscw)
	- [GET `/api/v1/settings/page`](#get-apiv1settingspage)
	- [POST `/api/v1/settings/cw`](#post-apiv1settingscw)
	- [POST `/api/v1/settings/page`](#post-apiv1settingspage)
- [Users](#users)
	- [GET `/api/v1/user/me`](#get-apiv1userme)
	- [PUT `/api/v1/user/email`](#put-apiv1useremail)
	- [PUT `/api/v1/user/password`](#put-apiv1userpassword)
- [CW Progress](#cw-progress)
	- [GET `/api/v1/cw/progress`](#get-apiv1cwprogress)
	- [PUT `/api/v1/cw/progress`](#put-apiv1cwprogress)
- [Utilities](#utilities)
	- [GET `/api/v1/hello`](#get-apiv1hello)

---

## Auth

---

## [POST] `/api/v1/auth/register`

**Description:** Create a new user account and immediately return an access token and refresh token.

**Authentication:** Not required

### Request

**Headers:**
| Header | Required | Description |
|---|---|---|
| Content-Type: application/json | Yes | Body must be JSON |

**Path Parameters:**
| Name | Type | Required | Description |
|---|---|---|---|
| None | - | - | This endpoint does not use path parameters |

**Query Parameters:**
| Name | Type | Required | Default | Description |
|---|---|---|---|---|
| None | - | - | - | This endpoint does not use query parameters |

**Request Body** (`application/json`):
```json
{
	"username": "alan_morse",
	"email": "alan@example.com",
	"password": "MyStrongPass!2026"
}
```

| Field | Type | Required | Description |
|---|---|---|---|
| username | string | Yes | 3-16 chars, must match `^[a-zA-Z0-9][a-zA-Z0-9_-]{1,14}[a-zA-Z0-9]$` |
| email | string | Yes | Valid email format, max 254 characters |
| password | string | Yes | 8-256 characters |

### Response

**Success `200`:**
```json
{
	"refresh_token": "d18f5e0956c7400fd1f7be47f98f83b4d4de2f0f74a68446f78d4ef287ccb947",
	"access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

| Field | Type | Description |
|---|---|---|
| refresh_token | string | Opaque token used to renew access tokens |
| access_token | string | JWT bearer token, expires in 15 minutes |

**Errors:**
| Status Code | Meaning |
|---|---|
| 400 | Invalid JSON body or validation failure |
| 409 | Username and/or email already exists |
| 500 | User creation or token issuance failed |

### Example (cURL)
```bash
curl -X POST https://api.example.com/api/v1/auth/register \
	-H "Content-Type: application/json" \
	-d '{"username":"alan_morse","email":"alan@example.com","password":"MyStrongPass!2026"}'
```

---

## [POST] `/api/v1/auth/login`

**Description:** Authenticate a user with username or email plus password, then return a new token pair.

**Authentication:** Not required

### Request

**Headers:**
| Header | Required | Description |
|---|---|---|
| Content-Type: application/json | Yes | Body must be JSON |

**Path Parameters:**
| Name | Type | Required | Description |
|---|---|---|---|
| None | - | - | This endpoint does not use path parameters |

**Query Parameters:**
| Name | Type | Required | Default | Description |
|---|---|---|---|---|
| None | - | - | - | This endpoint does not use query parameters |

**Request Body** (`application/json`):
```json
{
	"identifier": "alan@example.com",
	"password": "MyStrongPass!2026"
}
```

| Field | Type | Required | Description |
|---|---|---|---|
| identifier | string | Yes | Username or email |
| password | string | Yes | Plain-text password |

### Response

**Success `200`:**
```json
{
	"refresh_token": "c1606b1983af7ca6fce8e4c5c1d65d34f8e20ec1d1636279ca8bc0f335127823",
	"access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

| Field | Type | Description |
|---|---|---|
| refresh_token | string | Opaque token used to renew access tokens |
| access_token | string | JWT bearer token, expires in 15 minutes |

**Errors:**
| Status Code | Meaning |
|---|---|
| 400 | Invalid JSON body |
| 401 | Invalid credentials |
| 500 | Password verification or token issuance failed |

### Example (cURL)
```bash
curl -X POST https://api.example.com/api/v1/auth/login \
	-H "Content-Type: application/json" \
	-d '{"identifier":"alan@example.com","password":"MyStrongPass!2026"}'
```

---

## [POST] `/api/v1/auth/refresh`

**Description:** Exchange a valid refresh token for a new access token and a new refresh token. The previous refresh token is revoked.

**Authentication:** Not required

### Request

**Headers:**
| Header | Required | Description |
|---|---|---|
| Content-Type: application/json | Yes | Body must be JSON |

**Path Parameters:**
| Name | Type | Required | Description |
|---|---|---|---|
| None | - | - | This endpoint does not use path parameters |

**Query Parameters:**
| Name | Type | Required | Default | Description |
|---|---|---|---|---|
| None | - | - | - | This endpoint does not use query parameters |

**Request Body** (`application/json`):
```json
{
	"refresh_token": "d18f5e0956c7400fd1f7be47f98f83b4d4de2f0f74a68446f78d4ef287ccb947"
}
```

| Field | Type | Required | Description |
|---|---|---|---|
| refresh_token | string | Yes | Refresh token previously returned by register/login/refresh |

### Response

**Success `200`:**
```json
{
	"refresh_token": "f00baf9f67d06b8a4727363831c281302d2f13f7f4b2d453f6537673578db4ac",
	"access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

| Field | Type | Description |
|---|---|---|
| refresh_token | string | New refresh token; replaces the previous token |
| access_token | string | New JWT bearer token |

**Errors:**
| Status Code | Meaning |
|---|---|
| 400 | Invalid JSON body |
| 401 | Invalid or expired refresh token |
| 500 | Refresh transaction or token generation failed |

### Example (cURL)
```bash
curl -X POST https://api.example.com/api/v1/auth/refresh \
	-H "Content-Type: application/json" \
	-d '{"refresh_token":"d18f5e0956c7400fd1f7be47f98f83b4d4de2f0f74a68446f78d4ef287ccb947"}'
```

---

## Health

---

## [GET] `/api/v1/health`

**Description:** Basic server health check.

**Authentication:** Not required

### Request

**Headers:**
| Header | Required | Description |
|---|---|---|
| None | - | No special headers required |

**Path Parameters:**
| Name | Type | Required | Description |
|---|---|---|---|
| None | - | - | This endpoint does not use path parameters |

**Query Parameters:**
| Name | Type | Required | Default | Description |
|---|---|---|---|---|
| None | - | - | - | This endpoint does not use query parameters |

### Response

**Success `200`:**
```json
{
	"status": "healthy",
	"timestamp": 1774028400
}
```

| Field | Type | Description |
|---|---|---|
| status | string | Health state |
| timestamp | integer | Unix timestamp in seconds |

**Errors:**
| Status Code | Meaning |
|---|---|
| None | This endpoint always responds with HTTP 200 in normal operation |

### Example (cURL)
```bash
curl -X GET https://api.example.com/api/v1/health
```

---

## Settings

Protected settings endpoints require JWT auth middleware and loaded user context.

Authentication header format:

`Authorization: Bearer <access_token>`

If the header is missing/invalid or the token is invalid, the API returns `401`.

---

## [GET] `/api/v1/settings/`

**Description:** Return both CW settings and page settings for the authenticated user. If records are missing, defaults are returned.

**Authentication:** Required (`Authorization: Bearer <access_token>`)

### Request

**Headers:**
| Header | Required | Description |
|---|---|---|
| Authorization: Bearer <token> | Yes | Access token |

**Path Parameters:**
| Name | Type | Required | Description |
|---|---|---|---|
| None | - | - | This endpoint does not use path parameters |

**Query Parameters:**
| Name | Type | Required | Default | Description |
|---|---|---|---|---|
| None | - | - | - | This endpoint does not use query parameters |

### Response

**Success `200`:**
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

| Field | Type | Description |
|---|---|---|
| cw_settings | object | CW playback/training preferences |
| cw_settings.char_wpm | integer | Character speed in WPM |
| cw_settings.eff_wpm | integer | Effective speed in WPM |
| cw_settings.freq | integer | Tone frequency in Hz |
| cw_settings.start_delay | number | Delay before playback starts (seconds) |
| page_settings | object | UI/user page preferences |
| page_settings.theme | string | `auto`, `dark`, or `light` |
| page_settings.language | string | Language code/identifier used by frontend |
| page_settings.cur_lesson | integer | Current lesson index |

**Errors:**
| Status Code | Meaning |
|---|---|
| 401 | Missing/invalid token or user not found |
| 500 | Failed to fetch settings |

### Example (cURL)
```bash
curl -X GET https://api.example.com/api/v1/settings/ \
	-H "Authorization: Bearer <token>"
```

---

## [GET] `/api/v1/settings/cw`

**Description:** Return only CW settings for the authenticated user. If missing, default CW settings are returned.

**Authentication:** Required (`Authorization: Bearer <access_token>`)

### Request

**Headers:**
| Header | Required | Description |
|---|---|---|
| Authorization: Bearer <token> | Yes | Access token |

**Path Parameters:**
| Name | Type | Required | Description |
|---|---|---|---|
| None | - | - | This endpoint does not use path parameters |

**Query Parameters:**
| Name | Type | Required | Default | Description |
|---|---|---|---|---|
| None | - | - | - | This endpoint does not use query parameters |

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
| char_wpm | integer | Character speed in WPM |
| eff_wpm | integer | Effective speed in WPM |
| freq | integer | Tone frequency in Hz |
| start_delay | number | Delay before playback starts (seconds) |

**Errors:**
| Status Code | Meaning |
|---|---|
| 401 | Missing/invalid token or user not found |
| 500 | Failed to fetch settings |

### Example (cURL)
```bash
curl -X GET https://api.example.com/api/v1/settings/cw \
	-H "Authorization: Bearer <token>"
```

---

## [GET] `/api/v1/settings/page`

**Description:** Return only page settings for the authenticated user. If missing, default page settings are returned.

**Authentication:** Required (`Authorization: Bearer <access_token>`)

### Request

**Headers:**
| Header | Required | Description |
|---|---|---|
| Authorization: Bearer <token> | Yes | Access token |

**Path Parameters:**
| Name | Type | Required | Description |
|---|---|---|---|
| None | - | - | This endpoint does not use path parameters |

**Query Parameters:**
| Name | Type | Required | Default | Description |
|---|---|---|---|---|
| None | - | - | - | This endpoint does not use query parameters |

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
| theme | string | `auto`, `dark`, or `light` |
| language | string | Language code/identifier used by frontend |
| cur_lesson | integer | Current lesson index |

**Errors:**
| Status Code | Meaning |
|---|---|
| 401 | Missing/invalid token or user not found |
| 500 | Failed to fetch settings |

### Example (cURL)
```bash
curl -X GET https://api.example.com/api/v1/settings/page \
	-H "Authorization: Bearer <token>"
```

---

## [POST] `/api/v1/settings/cw`

**Description:** Create or update CW settings for the authenticated user.

**Authentication:** Required (`Authorization: Bearer <access_token>`)

### Request

**Headers:**
| Header | Required | Description |
|---|---|---|
| Authorization: Bearer <token> | Yes | Access token |
| Content-Type: application/json | Yes | Body must be JSON |

**Path Parameters:**
| Name | Type | Required | Description |
|---|---|---|---|
| None | - | - | This endpoint does not use path parameters |

**Query Parameters:**
| Name | Type | Required | Default | Description |
|---|---|---|---|---|
| None | - | - | - | This endpoint does not use query parameters |

**Request Body** (`application/json`):
```json
{
	"char_wpm": 24,
	"eff_wpm": 18,
	"freq": 700,
	"start_delay": 0.3
}
```

| Field | Type | Required | Description |
|---|---|---|---|
| char_wpm | integer | Yes | Character speed, min 5 max 50 |
| eff_wpm | integer | Yes | Effective speed, min 5 max 50 |
| freq | integer | Yes | Frequency in Hz, min 300 max 2000 |
| start_delay | number | Yes | Start delay in seconds, min 0.0 max 10.0 |

### Response

**Success `200`:**
```json
{
	"message": "Settings updated"
}
```

| Field | Type | Description |
|---|---|---|
| message | string | Operation result message |

**Errors:**
| Status Code | Meaning |
|---|---|
| 400 | Invalid JSON body or validation failure |
| 401 | Missing/invalid token or user not found |
| 500 | Failed to update settings |

### Example (cURL)
```bash
curl -X POST https://api.example.com/api/v1/settings/cw \
	-H "Authorization: Bearer <token>" \
	-H "Content-Type: application/json" \
	-d '{"char_wpm":24,"eff_wpm":18,"freq":700,"start_delay":0.3}'
```

---

## [POST] `/api/v1/settings/page`

**Description:** Create or update page settings for the authenticated user.

**Authentication:** Required (`Authorization: Bearer <access_token>`)

### Request

**Headers:**
| Header | Required | Description |
|---|---|---|
| Authorization: Bearer <token> | Yes | Access token |
| Content-Type: application/json | Yes | Body must be JSON |

**Path Parameters:**
| Name | Type | Required | Description |
|---|---|---|---|
| None | - | - | This endpoint does not use path parameters |

**Query Parameters:**
| Name | Type | Required | Default | Description |
|---|---|---|---|---|
| None | - | - | - | This endpoint does not use query parameters |

**Request Body** (`application/json`):
```json
{
	"theme": "dark",
	"language": "en-US",
	"cur_lesson": 12
}
```

| Field | Type | Required | Description |
|---|---|---|---|
| theme | string | Yes | Allowed values: `auto`, `dark`, `light` |
| language | string | Yes | UI language identifier |
| cur_lesson | integer | Yes | Current lesson index |

### Response

**Success `200`:**
```json
{
	"message": "Settings updated"
}
```

| Field | Type | Description |
|---|---|---|
| message | string | Operation result message |

**Errors:**
| Status Code | Meaning |
|---|---|
| 400 | Invalid JSON body or validation failure |
| 401 | Missing/invalid token or user not found |
| 500 | Failed to update settings |

### Example (cURL)
```bash
curl -X POST https://api.example.com/api/v1/settings/page \
	-H "Authorization: Bearer <token>" \
	-H "Content-Type: application/json" \
	-d '{"theme":"dark","language":"en-US","cur_lesson":12}'
```

---

## Users

---

## [GET] `/api/v1/user/me`

**Description:** Return the authenticated user's profile summary.

**Authentication:** Required (`Authorization: Bearer <access_token>`)

### Request

**Headers:**
| Header | Required | Description |
|---|---|---|
| Authorization: Bearer <token> | Yes | Access token |

**Path Parameters:**
| Name | Type | Required | Description |
|---|---|---|---|
| None | - | - | This endpoint does not use path parameters |

**Query Parameters:**
| Name | Type | Required | Default | Description |
|---|---|---|---|---|
| None | - | - | - | This endpoint does not use query parameters |

### Response

**Success `200`:**
```json
{
	"username": "alan_morse",
	"email": "alan@example.com",
	"created_at": "2026-03-20T15:04:05Z"
}
```

| Field | Type | Description |
|---|---|---|
| username | string | Account username |
| email | string | Account email |
| created_at | string (ISO 8601 datetime) | Account creation timestamp |

**Errors:**
| Status Code | Meaning |
|---|---|
| 401 | Missing/invalid token or user not found |

### Example (cURL)
```bash
curl -X GET https://api.example.com/api/v1/user/me \
	-H "Authorization: Bearer <token>"
```

---

## [PUT] `/api/v1/user/email`

**Description:** Update the authenticated user's email.

**Authentication:** Required (`Authorization: Bearer <access_token>`)

### Request

**Headers:**
| Header | Required | Description |
|---|---|---|
| Authorization: Bearer <token> | Yes | Access token |
| Content-Type: application/json | Yes | Body must be JSON |

**Path Parameters:**
| Name | Type | Required | Description |
|---|---|---|---|
| None | - | - | This endpoint does not use path parameters |

**Query Parameters:**
| Name | Type | Required | Default | Description |
|---|---|---|---|---|
| None | - | - | - | This endpoint does not use query parameters |

**Request Body** (`application/json`):
```json
{
	"email": "alan.new@example.com"
}
```

| Field | Type | Required | Description |
|---|---|---|---|
| email | string | Yes | Valid email format, max 254 characters |

### Response

**Success `200`:**
```json
{
	"message": "Email updated"
}
```

| Field | Type | Description |
|---|---|---|
| message | string | Operation result message |

**Errors:**
| Status Code | Meaning |
|---|---|
| 400 | Invalid body or new email equals current email |
| 401 | Missing/invalid token or user not found |
| 409 | Email already in use |
| 500 | Failed to update email |

### Example (cURL)
```bash
curl -X PUT https://api.example.com/api/v1/user/email \
	-H "Authorization: Bearer <token>" \
	-H "Content-Type: application/json" \
	-d '{"email":"alan.new@example.com"}'
```

---

## [PUT] `/api/v1/user/password`

**Description:** Update the authenticated user's password.

**Authentication:** Required (`Authorization: Bearer <access_token>`)

### Request

**Headers:**
| Header | Required | Description |
|---|---|---|
| Authorization: Bearer <token> | Yes | Access token |
| Content-Type: application/json | Yes | Body must be JSON |

**Path Parameters:**
| Name | Type | Required | Description |
|---|---|---|---|
| None | - | - | This endpoint does not use path parameters |

**Query Parameters:**
| Name | Type | Required | Default | Description |
|---|---|---|---|---|
| None | - | - | - | This endpoint does not use query parameters |

**Request Body** (`application/json`):
```json
{
	"old_password": "MyStrongPass!2026",
	"new_password": "ASecondStrongPass!2026"
}
```

| Field | Type | Required | Description |
|---|---|---|---|
| old_password | string | Yes | Current password, 8-256 characters |
| new_password | string | Yes | New password, 8-256 characters |

### Response

**Success `200`:**
```json
{
	"message": "Password updated"
}
```

| Field | Type | Description |
|---|---|---|
| message | string | Operation result message |

**Errors:**
| Status Code | Meaning |
|---|---|
| 400 | Invalid JSON body or validation failure |
| 401 | Invalid credentials, missing/invalid token, or user not found |
| 500 | Password hashing/comparison or database update failed |

### Example (cURL)
```bash
curl -X PUT https://api.example.com/api/v1/user/password \
	-H "Authorization: Bearer <token>" \
	-H "Content-Type: application/json" \
	-d '{"old_password":"MyStrongPass!2026","new_password":"ASecondStrongPass!2026"}'
```

---

## CW Progress

---

## [GET] `/api/v1/cw/progress`

**Description:** Return all stored lesson progress rows for the authenticated user.

**Authentication:** Required (`Authorization: Bearer <access_token>`)

### Request

**Headers:**
| Header | Required | Description |
|---|---|---|
| Authorization: Bearer <token> | Yes | Access token |

**Path Parameters:**
| Name | Type | Required | Description |
|---|---|---|---|
| None | - | - | This endpoint does not use path parameters |

**Query Parameters:**
| Name | Type | Required | Default | Description |
|---|---|---|---|---|
| None | - | - | - | This endpoint does not use query parameters |

### Response

**Success `200`:**
```json
{
	"data": [
		{
			"lesson": "12",
			"char_wpm": 24,
			"eff_wpm": 18,
			"accuracy": 0.96,
			"created_at": "2026-03-20T15:18:40Z"
		}
	]
}
```

| Field | Type | Description |
|---|---|---|
| data | array | List of progress rows |
| data[].lesson | string | Lesson identifier stored in response as string |
| data[].char_wpm | integer | Character speed in WPM |
| data[].eff_wpm | integer | Effective speed in WPM |
| data[].accuracy | number | Accuracy value between `0.0` and `1.0` |
| data[].created_at | string (ISO 8601 datetime) | Progress record creation time |

**Errors:**
| Status Code | Meaning |
|---|---|
| 401 | Missing/invalid token or user not found |
| 500 | Failed to query progress |

### Example (cURL)
```bash
curl -X GET https://api.example.com/api/v1/cw/progress \
	-H "Authorization: Bearer <token>"
```

---

## [PUT] `/api/v1/cw/progress`

**Description:** Create a new progress row for the authenticated user.

**Authentication:** Required (`Authorization: Bearer <access_token>`)

### Request

**Headers:**
| Header | Required | Description |
|---|---|---|
| Authorization: Bearer <token> | Yes | Access token |
| Content-Type: application/json | Yes | Body must be JSON |

**Path Parameters:**
| Name | Type | Required | Description |
|---|---|---|---|
| None | - | - | This endpoint does not use path parameters |

**Query Parameters:**
| Name | Type | Required | Default | Description |
|---|---|---|---|---|
| None | - | - | - | This endpoint does not use query parameters |

**Request Body** (`application/json`):
```json
{
	"lesson": 12,
	"char_wpm": 24,
	"eff_wpm": 18,
	"accuracy": 0.96
}
```

| Field | Type | Required | Description |
|---|---|---|---|
| lesson | integer | Yes | Lesson index |
| char_wpm | integer | Yes | Character speed, min 5 max 50 |
| eff_wpm | integer | Yes | Effective speed, min 5 max 50 |
| accuracy | number | Yes | Value between `0.0` and `1.0` |

### Response

**Success `201`:**
```json
{
	"message": "Progress Created"
}
```

| Field | Type | Description |
|---|---|---|
| message | string | Operation result message |

**Errors:**
| Status Code | Meaning |
|---|---|
| 400 | Invalid JSON body or validation failure |
| 401 | Missing/invalid token or user not found |
| 500 | Failed to create progress |

### Example (cURL)
```bash
curl -X PUT https://api.example.com/api/v1/cw/progress \
	-H "Authorization: Bearer <token>" \
	-H "Content-Type: application/json" \
	-d '{"lesson":12,"char_wpm":24,"eff_wpm":18,"accuracy":0.96}'
```

---

## Utilities

---

## [GET] `/api/v1/hello`

**Description:** Simple authenticated test endpoint that confirms token auth and loaded user context.

**Authentication:** Required (`Authorization: Bearer <access_token>`)

### Request

**Headers:**
| Header | Required | Description |
|---|---|---|
| Authorization: Bearer <token> | Yes | Access token |

**Path Parameters:**
| Name | Type | Required | Description |
|---|---|---|---|
| None | - | - | This endpoint does not use path parameters |

**Query Parameters:**
| Name | Type | Required | Default | Description |
|---|---|---|---|---|
| None | - | - | - | This endpoint does not use query parameters |

### Response

**Success `200`:**
```json
{
	"message": "Hello, authenticated user {alan_morse}!"
}
```

| Field | Type | Description |
|---|---|---|
| message | string | Greeting including the authenticated username |

**Errors:**
| Status Code | Meaning |
|---|---|
| 401 | Missing/invalid token or user not found |

### Example (cURL)
```bash
curl -X GET https://api.example.com/api/v1/hello \
	-H "Authorization: Bearer <token>"
```

---

## Authentication and Constraints

- Protected endpoints are guarded by JWT middleware and user-loading middleware.
- Required auth header format: `Authorization: Bearer <token>`.
- Access tokens are JWTs with 15-minute expiry.
- Refresh tokens are opaque strings with ~30-day expiry and are rotated on each `/auth/refresh` request.
- No explicit HTTP rate limiting middleware is configured in this codebase.

## Error Reference

| Status Code | Meaning |
|---|---|
| 400 | Bad Request — Invalid or missing parameters |
| 401 | Unauthorized — Missing or invalid token |
| 403 | Forbidden — Insufficient permissions |
| 404 | Not Found — Resource does not exist |
| 429 | Too Many Requests — Rate limit exceeded |
| 500 | Internal Server Error — Unexpected server failure |
