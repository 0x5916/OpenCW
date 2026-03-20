# API Documentation

**Base URL:** `https://api.example.com`
> Last updated: March 21, 2026

## Table of Contents
- [Auth](#auth)
  - [POST /v1/auth/register](#post-v1authregister)
  - [POST /v1/auth/login](#post-v1authlogin)
  - [POST /v1/auth/refresh](#post-v1authrefresh)
- [Settings](#settings)
  - [GET /v1/settings/all](#get-v1settingsall)
  - [GET /v1/settings/cw](#get-v1settingscw)
  - [GET /v1/settings/page](#get-v1settingspage)
  - [POST /v1/settings/cw](#post-v1settingscw)
  - [POST /v1/settings/page](#post-v1settingspage)
- [User](#user)
  - [GET /v1/user/me](#get-v1userme)
  - [PUT /v1/user/email](#put-v1useremail)
  - [PUT /v1/user/password](#put-v1userpassword)
- [Progress](#progress)
  - [GET /v1/cw/progress](#get-v1cwprogress)
  - [PUT /v1/cw/progress](#put-v1cwprogress)
- [Misc](#misc)
  - [GET /v1/health](#get-v1health)
  - [GET /v1/hello](#get-v1hello)

---

## Auth

### POST `/v1/auth/register`

**Description:** Register a new user account.

**Authentication:** Not required

### Request

**Headers:**
| Header | Required | Description |
|---|---|---|
| Content-Type | Yes | Must be `application/json` |

**Path Parameters:**
| Name | Type | Required | Description |
|---|---|---|---|

**Query Parameters:**
| Name | Type | Required | Default | Description |
|---|---|---|---|---|

**Request Body** (`application/json`):
```json
{
  "username": "morsefan",
  "email": "alice@example.com",
  "password": "MySecurePassw0rd!"
}
```

| Field | Type | Required | Description |
|---|---|---|---|
| username | string | Yes | Unique username, 3-32 chars, a-z, 0-9, _ allowed |
| email | string | Yes | Valid email address, max 254 chars |
| password | string | Yes | 8-256 chars |

### Response

**Success `[200]`:**
```json
{
  "refresh_token": "...",
  "access_token": "..."
}
```

| Field | Type | Description |
|---|---|---|
| refresh_token | string | Token for refreshing session |
| access_token | string | JWT for authenticated requests |

**Errors:**
| Status Code | Meaning |
|---|---|
| 400 | Invalid request body |
| 409 | Username or email already exists |
| 500 | Database or server error |

### Example (cURL)
```bash
curl -X POST https://api.example.com/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"morsefan","email":"alice@example.com","password":"MySecurePassw0rd!"}'
```

---

### POST `/v1/auth/login`

**Description:** Log in with username/email and password to receive tokens.

**Authentication:** Not required

### Request

**Headers:**
| Header | Required | Description |
|---|---|---|
| Content-Type | Yes | Must be `application/json` |

**Path Parameters:**
| Name | Type | Required | Description |
|---|---|---|---|

**Query Parameters:**
| Name | Type | Required | Default | Description |
|---|---|---|---|---|

**Request Body** (`application/json`):
```json
{
  "identifier": "alice@example.com",
  "password": "MySecurePassw0rd!"
}
```

| Field | Type | Required | Description |
|---|---|---|---|
| identifier | string | Yes | Username or email |
| password | string | Yes | User password |

### Response

**Success `[200]`:**
```json
{
  "refresh_token": "...",
  "access_token": "..."
}
```

| Field | Type | Description |
|---|---|---|
| refresh_token | string | Token for refreshing session |
| access_token | string | JWT for authenticated requests |

**Errors:**
| Status Code | Meaning |
|---|---|
| 400 | Invalid request body |
| 401 | Invalid credentials |
| 500 | Internal server error |

### Example (cURL)
```bash
curl -X POST https://api.example.com/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"identifier":"alice@example.com","password":"MySecurePassw0rd!"}'
```

---

### POST `/v1/auth/refresh`

**Description:** Obtain a new access token using a refresh token.

**Authentication:** Not required

### Request

**Headers:**
| Header | Required | Description |
|---|---|---|
| Content-Type | Yes | Must be `application/json` |

**Path Parameters:**
| Name | Type | Required | Description |
|---|---|---|---|

**Query Parameters:**
| Name | Type | Required | Default | Description |
|---|---|---|---|---|

**Request Body** (`application/json`):
```json
{
  "refresh_token": "..."
}
```

| Field | Type | Required | Description |
|---|---|---|---|
| refresh_token | string | Yes | Valid refresh token |

### Response

**Success `[200]`:**
```json
{
  "refresh_token": "...",
  "access_token": "..."
}
```

| Field | Type | Description |
|---|---|---|
| refresh_token | string | New refresh token |
| access_token | string | New access token |

**Errors:**
| Status Code | Meaning |
|---|---|
| 400 | Invalid request body |
| 401 | Invalid or expired refresh token |
| 500 | Internal server error |

### Example (cURL)
```bash
curl -X POST https://api.example.com/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{"refresh_token":"..."}'
```

---

## Settings

### GET `/v1/settings/all`

**Description:** Get all user settings (CW and page settings).

**Authentication:** Required (Bearer token)

### Request

**Headers:**
| Header | Required | Description |
|---|---|---|
| Authorization | Yes | Bearer access token |

**Path Parameters:**
| Name | Type | Required | Description |
|---|---|---|---|

**Query Parameters:**
| Name | Type | Required | Default | Description |
|---|---|---|---|---|

**Request Body** (`application/json`):

_None_

### Response

**Success `[200]`:**
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

| Field | Type | Description |
|---|---|---|
| cw_settings | object | User's CW (Morse) settings |
| page_settings | object | User's page settings |

**Errors:**
| Status Code | Meaning |
|---|---|
| 401 | Unauthorized |
| 500 | Failed to get settings |

### Example (cURL)
```bash
curl -X GET https://api.example.com/v1/settings/all \
  -H "Authorization: Bearer <token>"
```

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

