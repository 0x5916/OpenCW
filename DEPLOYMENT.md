# Production Deployment

OpenCW is deployed as three Docker containers orchestrated by Docker Compose:

| Service    | Image                            | Port |
|------------|----------------------------------|------|
| `db`       | `postgres:18-alpine`             | —    |
| `backend`  | built from `./backend`           | 8080 |
| `frontend` | built from `./frontend`          | 3000 |

## Requirements

- Docker ≥ 26 with the Compose plugin (`docker compose version`)
- A reverse proxy (nginx, Caddy, etc.) in front of port 3000 and 8080 for TLS termination

## 1. Clone and configure

```bash
git clone <repo-url>
cd OpenCW
cp example.env .env
```

Edit `.env` and fill in every value:

```dotenv
# ── PostgreSQL ────────────────────────────────────────────────────────────────
POSTGRES_USER=opencw
POSTGRES_PASSWORD=<strong-random-password>
POSTGRES_DB=opencw

# ── Backend ───────────────────────────────────────────────────────────────────
# Must be base64-encoded bytes (≥ 32 raw bytes).
JWT_SECRET=<output of: openssl rand -base64 32>

# Comma-separated list of origins the browser sends the API requests from.
# Set this to your public frontend URL(s), e.g.:
CORS_ORIGINS=https://opencw.example.com

# ── Frontend (build-time) ─────────────────────────────────────────────────────
# The URL the *browser* uses to reach the backend API.
# Must be publicly reachable — this is baked into the frontend bundle at build time.
PUBLIC_API_BASE=https://api.opencw.example.com/api/v1
```

### Generating a JWT secret

```bash
openssl rand -base64 32
```

The value must be base64-encoded. The backend decodes it at startup and exits if it is missing or malformed.

## 2. Build and start

```bash
docker compose up -d --build
```

Compose will:
1. Pull `postgres:18-alpine` and start the database.
2. Wait for the database healthcheck to pass.
3. Build and start the backend (Go, distroless image, `GIN_MODE=release`).
4. Build and start the frontend (SvelteKit Node adapter, `NODE_ENV=production`).

Check everything is healthy:

```bash
docker compose ps
docker compose logs --tail=50
```

## 3. Reverse proxy

Neither service should be exposed to the public internet directly. Terminate TLS in your reverse proxy and forward:

| Public URL                            | Upstream                 |
|---------------------------------------|--------------------------|
| `https://opencw.example.com`          | `localhost:3000`         |
| `https://api.opencw.example.com`      | `localhost:8080`         |

> If you host both under a single domain (e.g. `/api/v1` path prefix), update
> `PUBLIC_API_BASE` accordingly and rebuild the frontend image.

### Example: Caddy

```
opencw.example.com {
    reverse_proxy localhost:3000
}

api.opencw.example.com {
    reverse_proxy localhost:8080
}
```

### Example: nginx

```nginx
server {
    listen 443 ssl;
    server_name opencw.example.com;
    location / { proxy_pass http://localhost:3000; }
}

server {
    listen 443 ssl;
    server_name api.opencw.example.com;
    location / { proxy_pass http://localhost:8080; }
}
```

## 4. Updates

```bash
git pull
docker compose up -d --build
```

Compose will rebuild changed images and recreate only the affected containers. The `pgdata` volume is preserved across updates.

## 5. Backups

The database is stored in the `opencw_pgdata` Docker volume. Back it up with:

```bash
docker compose exec db pg_dump -U "$POSTGRES_USER" "$POSTGRES_DB" | gzip > opencw_$(date +%F).sql.gz
```

Restore:

```bash
gunzip -c opencw_<date>.sql.gz | docker compose exec -T db psql -U "$POSTGRES_USER" "$POSTGRES_DB"
```

## 6. Stopping / removing

Stop without removing data:
```bash
docker compose down
```

Stop and **delete all data** (destructive):
```bash
docker compose down -v
```

## Environment variable reference

| Variable          | Required | Description |
|-------------------|----------|-------------|
| `POSTGRES_USER`   | yes      | PostgreSQL superuser name |
| `POSTGRES_PASSWORD` | yes    | PostgreSQL superuser password |
| `POSTGRES_DB`     | yes      | Database name |
| `JWT_SECRET`      | yes      | Base64-encoded secret (≥ 32 raw bytes). Backend exits on startup if missing. |
| `CORS_ORIGINS`    | yes      | Comma-separated allowed browser origins, e.g. `https://opencw.example.com` |
| `PUBLIC_API_BASE` | yes      | Browser-visible backend URL, baked into the frontend at build time, e.g. `https://api.opencw.example.com/api/v1` |
