# OpenCW

OpenCW is a full-stack Morse code training platform with:

- Go backend API (Gin + GORM)
- SvelteKit frontend
- PostgreSQL database
- Docker Compose orchestration

## Project Structure

- [backend](backend)
- [frontend](frontend)
- [api_test](api_test)
- [docker-compose.yaml](docker-compose.yaml)
- [example.env](example.env)
- [DEPLOYMENT.md](DEPLOYMENT.md)

## Quick Start (Docker)

1. Copy environment template:

   cp example.env .env

2. Generate JWT secret:

   openssl rand -base64 32

3. Put the generated value into .env as JWT_SECRET.

4. Start all services:

   docker compose up -d --build

5. Open the app:

   http://localhost:3000

## Services

- Frontend: http://localhost:3000
- Backend API base: http://localhost:8080/api/v1
- Health check: http://localhost:8080/api/v1/health

## Common Commands

Start / rebuild:

docker compose up -d --build

View logs:

docker compose logs -f

Restart backend only:

docker compose up -d backend

Stop without removing data:

docker compose down

Stop and remove data volumes:

docker compose down -v

## Environment Variables

See [example.env](example.env) for the full list. Required variables:

- POSTGRES_USER
- POSTGRES_PASSWORD
- POSTGRES_DB
- JWT_SECRET (must be base64)
- CORS_ORIGINS
- PUBLIC_API_BASE

## API Testing

API request collections and environment files are in [api_test](api_test).

## Production Deployment

Use [DEPLOYMENT.md](DEPLOYMENT.md) for production setup, reverse proxy, backup, and update procedures.

## Notes

- The frontend uses PUBLIC_API_BASE at build time. If this value changes, rebuild the frontend image.
- If browser requests are blocked by CORS, ensure CORS_ORIGINS contains the exact frontend origin you open in the browser.
- PostgreSQL settings in [docker-compose.yaml](docker-compose.yaml) are tuned for an 8 GB RAM / 4 vCPU Linux server (LXC container profile). Re-tune if your host resources differ.
