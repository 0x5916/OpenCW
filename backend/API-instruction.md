You are a senior backend engineer and technical writer. You have full access to this codebase. Your task is to autonomously read the source code and generate complete API documentation for frontend developers, saved as `API.md`.

---

### Step 1 — Explore the Project

Begin by reading the project structure. Then:
1. Open and read `main.go` — identify the router, all middleware, and all route groups
2. Follow every registered route to its handler file and read the full handler logic
3. Read all request/response struct definitions referenced by those handlers
4. Read all middleware (auth, rate limiting, etc.) to infer constraints
5. Repeat until every registered route is fully understood

Do not skip any file. If a route references a function in another file, go read that file.

---

### Step 2 — Extract Per-Endpoint Information

For each endpoint, extract:
- HTTP method and full path (with route group prefixes)
- All parameters: path, query, and body — name, type, required, validation rules
- Required request headers
- Successful response body — all fields, types, and meanings
- All possible error responses with status codes
- Auth requirements from middleware
- Rate limiting rules if present

---

### Step 3 — Write `API.md`

Write the entire output as valid Markdown, ready to save as `API.md`.

Start the file with:

# API Documentation

**Base URL:** `https://api.example.com`
> Last updated: [today's date]

## Table of Contents
[generate based on all discovered endpoints, grouped by resource]

---

Then document each endpoint using this template:

---

## [METHOD] `[full path]`

**Description:** [Plain English explanation of what this endpoint does]

**Authentication:** [Required / Not required — specify method if required]

### Request

**Headers:**
| Header | Required | Description |
|---|---|---|

**Path Parameters:**
| Name | Type | Required | Description |
|---|---|---|---|

**Query Parameters:**
| Name | Type | Required | Default | Description |
|---|---|---|---|---|

**Request Body** (`application/json`):
\```json
{
  // realistic example payload
}
\```

| Field | Type | Required | Description |
|---|---|---|---|

### Response

**Success `[status code]`:**
\```json
{
  // realistic example response
}
\```

| Field | Type | Description |
|---|---|---|

**Errors:**
| Status Code | Meaning |
|---|---|

### Example (cURL)
\```bash
curl -X [METHOD] https://api.example.com/[path] \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{...}'
\```

---

Group all endpoints by resource to match the route groups in `main.go` (e.g., `## Auth`, `## Users`).

---

### Step 4 — Append Global Error Reference

At the bottom of `API.md`, always append:

## Error Reference

| Status Code | Meaning |
|---|---|
| 400 | Bad Request — Invalid or missing parameters |
| 401 | Unauthorized — Missing or invalid token |
| 403 | Forbidden — Insufficient permissions |
| 404 | Not Found — Resource does not exist |
| 429 | Too Many Requests — Rate limit exceeded |
| 500 | Internal Server Error — Unexpected server failure |

---

### Writing Rules
- Audience is frontend developers — never expose internal implementation details
- Use plain, active-voice English
- Use realistic example values, never "foo", "test", or "123"
- Note any response field that may be unintuitive to a frontend developer
- Mark deprecated fields with ⚠️ **Deprecated**

---

Begin now. Start with `main.go`, read all referenced files, and produce the complete `API.md`.
