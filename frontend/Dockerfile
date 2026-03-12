# syntax=docker/dockerfile:1

FROM node:22-alpine AS build
WORKDIR /app

# Install dependencies first for better layer caching.
COPY package.json package-lock.json ./
RUN npm ci

COPY . .

# Required by $env/static/public imports during `npm run build`.
ARG PUBLIC_API_BASE=http://localhost:8000
ENV PUBLIC_API_BASE=$PUBLIC_API_BASE

# Generate .svelte-kit files in clean build environments, then build.
RUN npm run prepare && npm run build

FROM node:20-alpine AS runtime
WORKDIR /app
ENV NODE_ENV=production
ENV PORT=3000

# Install runtime-only dependencies for a smaller production image.
COPY package.json package-lock.json ./
RUN npm ci --omit=dev && npm cache clean --force

COPY --from=build /app/build ./build

USER node
EXPOSE 3000
CMD ["node", "build"]
