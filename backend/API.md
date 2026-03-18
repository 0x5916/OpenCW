# OpenCW Backend API (v1)

Base URL: `https://<your-host>` (production) or `http://localhost:<PORT>` (local). The server listens on `:` + value of `configs.App.Port` (default configured in `.env`).

Authentication:
- Type: JWT access token (Bearer). Include header `Authorization: Bearer <access_token>` for protected endpoints.
- Access tokens expire ~15 minutes. Use the refresh endpoint to get a new access token + refresh token pair.
- Refresh tokens are opaque strings returned by the server. When refreshing, send the raw refresh token value in the `refresh_token` body field.

Common response shapes
- Success message: `{"message": "..."}`
- Error shape: `{"error": "..."}` — server returns appropriate HTTP status (400, 401, 403, 404, 409, 500).

Top-level routes
- `GET /api/v1/health` — health check. Returns: `{"status":"healthy","timestamp":<unix>}`

Auth
- `POST /api/v1/auth/register` — Register new user
  - Request JSON:
    - `username` (string, required, custom username validator)
    - `email` (string, required, email)
    - `password` (string, required, min 8)
  - Success (200): `AuthTokenPairResponse` — `{"refresh_token": "<raw>", "access_token": "<jwt>"}`
  - Errors: 400 (validation), 409 (username/email existing), 500 (db/crypto)

- `POST /api/v1/auth/login` — Login with username or email
  - Request JSON:
    - `identifier` (string, required) — username or email
    - `password` (string, required)
  - Success (200): `AuthTokenPairResponse`
  - Errors: 400, 401 (invalid credentials), 500

- `POST /api/v1/auth/refresh` — Exchange refresh token for new pair
  - Request JSON:
    - `refresh_token` (string, required) — raw refresh token previously returned
  - Success (200): `AuthTokenPairResponse` (new refresh token + access token)
  - Errors: 400, 401 (invalid/expired token), 500

User (Protected — requires `Authorization: Bearer <access_token>`)
- `GET /api/v1/user/me` — Get current user's basic info
  - Success (200): `UserInfoResponse` — `{"username":"...","email":"...","created_at":"..."}`

- `PUT /api/v1/user/email` — Update current user's email
  - Request JSON:
    - `email` (string, required, email)
  - Success (200): `{"message":"Email updated"}`
  - Errors: 400, 409 (email conflict), 500

- `PUT /api/v1/user/password` — Update password
  - Request JSON:
    - `old_password` (string, required)
    - `new_password` (string, required, min 8)
  - Success (200): `{"message":"Password updated"}`
  - Errors: 400, 401 (old password mismatch), 500

CW Settings (Protected; user loaded via middleware)
- `GET /api/v1/cw/settings` — Get user's CW (Morse/typing) settings
  - Success (200): `CWSettingsResponse` — `{"char_wpm":int,"eff_wpm":int,"freq":int,"start_delay":float}`
  - If no saved settings, defaults from `models.GetDefaultCWSettings()` are returned.

- `POST /api/v1/cw/settings` — Create or update CW settings
  - Request JSON (CWSettingsInput):
    - `char_wpm` (int, required, 5..50)
    - `eff_wpm` (int, required, 5..50)
    - `freq` (int, required, 300..2000)
    - `start_delay` (float, required, 0.0..10.0)
  - Success (200): `{"message":"Settings updated"}`
  - Errors: 400 (validation), 500

Page Settings (Protected)
- `GET /api/v1/page/settings` — Get user's page settings
  - Success (200): `PageSettingsResponse` — `{"theme":"auto|dark|light","language":"...","cur_lesson":int}`
  - If none, returns `models.GetDefaultPageSettings()`.

- `POST /api/v1/page/settings` — Create or update page settings
  - Request JSON (PageSettingsInput):
    - `theme` (string, required) — one of `auto`, `dark`, `light`
    - `language` (string, required)
    - `cur_lesson` (int, required)
  - Success (200): `{"message":"Settings updated"}`
  - Errors: 400, 500

Progress (Protected)
- `GET /api/v1/cw/progress` — Get all progress records for current user
  - Success (200): `{"data": [ProgressResponse, ...]}`
  - `ProgressResponse` shape: `{"lesson":string,"char_wpm":int,"eff_wpm":int,"accuracy":float,"created_at":"..."}`

- `PUT /api/v1/cw/progress` — Add a new progress entry
  - Request JSON (ProgressInput):
    - `lesson` (int, required)
    - `char_wpm` (int, required, 5..50)
    - `eff_wpm` (int, required, 5..50)
    - `accuracy` (float, required, 0.0..1.0)
  - Success (201): `{"message":"Progress Created"}`
  - Errors: 400, 500

Misc
- `GET /api/v1/hello` — Protected test endpoint, returns `{"message":"Hello, authenticated user!"}`

Headers
- `Content-Type: application/json` for JSON bodies.
- `Authorization: Bearer <access_token>` for protected endpoints.

Validation notes
- Validation tags are defined in `handlers/v1/common` input structs. See code:
  - `handlers/v1/common/input.go` ([file](handlers/v1/common/input.go))

- Examples
- Register (curl):

```bash
curl -X POST https://api.example.com/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","email":"alice@example.com","password":"s3curepass"}'
```

- Login (curl):

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"identifier":"alice","password":"s3curepass"}'
```

- Fetch CW settings (curl):

```bash
curl -H "Authorization: Bearer <access_token>" \
  http://localhost:8080/api/v1/cw/settings
```

Errors and Troubleshooting
- 401 Unauthorized — missing/invalid token or refresh token
- 409 Conflict — uniqueness constraint (username/email)
- 400 Bad Request — validation failures; response body contains `{"error":"<message>"}`
- 500 Internal Server Error — DB or cryptographic failures

Implementation pointers for frontend
- When storing tokens: store refresh token securely (httpOnly cookie recommended); access token can be stored in memory and refreshed when expired.
- Use `POST /api/v1/auth/refresh` with `refresh_token` to obtain a new `access_token`/`refresh_token` pair.

Server internals (Go models and DB schema) are intentionally omitted from this document — the frontend should rely on the JSON request/response shapes above. If you need machine-readable schemas or language-specific types (OpenAPI, JSON Schema, TypeScript), I can generate those instead.

If you want, I can:
- add OpenAPI/Swagger spec
- generate TypeScript client types
- include example responses for each endpoint
