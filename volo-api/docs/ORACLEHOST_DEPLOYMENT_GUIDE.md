# OracleHost Deployment Guide

This document is the OracleHost-specific preparation guide for Volo.

Its purpose is simple:

- another agent can use it to make Volo compatible with the OracleHost deployment model
- final host wiring, real infra values, and live server operations stay in private infrastructure notes

This guide intentionally avoids publishing live host details such as real IPs, hostnames, deployed app inventory, backup schedules, or private routing state.

## Why This Exists

`volo-api/docs/DEPLOYMENT.md` describes a generic standalone VM layout under `/opt/volo`.

OracleHost uses a different contract. The important differences are:

- per-app slots under `/opt/apps/<app-id>`
- server-local env files under `/etc/apps/<app-id>/<app-env-file>`
- Caddy route snippets managed at the platform level
- backend deploys triggered through `deploy-oracle-app`
- GitHub Actions plus GHCR for backend delivery
- frontend and backend treated as separate deployment concerns

If you are preparing Volo for OracleHost, use this guide first and treat the generic deployment doc as background only.

## OracleHost Contract Summary

OracleHost expects backend apps to follow this model:

1. The backend is packaged as a container image.
2. GitHub Actions builds and pushes that image to GHCR.
3. CI deploys over SSH as a narrow deploy user.
4. The deploy user runs:

```bash
sudo /usr/local/sbin/deploy-oracle-app <app-id> <ghcr-image-ref>
```

5. `deploy-oracle-app` updates `/opt/apps/<app-id>/stack/.deploy.env`.
6. Docker Compose is run from `/opt/apps/<app-id>/stack`.
7. The deploy verifies:

```text
http://127.0.0.1:${PORT}/api/health
```

Platform assumptions:

- the backend listens on the runtime port env var
- the backend exposes `/api/health`
- secrets stay out of Git
- Caddy is the public HTTP entrypoint
- databases remain private to the host

## Recommended Environment Model

For this project, keep the deployment model simple:

- `local` is the real development and staging environment
- `production` is the public deployment

There is no need to maintain a separate hosted `dev` environment unless the project grows enough to justify it.

That means the OracleHost preparation in this repo should focus on one hosted app slot:

### Hosted production

```text
APP_ID=volo
APP_SLUG=volo
APP_DIR=/opt/apps/volo
ENV_FILE=/etc/apps/volo/volo.env
DB_NAME=volo
```

Notes:

- local development remains outside OracleHost
- production lives on OracleHost and Cloudflare
- frontend and backend should stay separable
- final domains, ports, hostnames, and deploy usernames belong in private infra notes or GitHub environment values
- if you later need a second hosted environment, use the standard OracleHost `volo-dev` slot then

## What Volo Already Has

Useful current pieces already present in this repo:

- Go backend with configurable runtime port
- PostgreSQL support
- Redis support
- Prometheus `/metrics`
- Sentry configuration support
- Dockerfiles for the API
- documented runtime env vars

Useful files:

- `volo-api/Dockerfile`
- `volo-api/cmd/server/main.go`
- `volo-api/internal/config/config.go`
- `volo-api/docs/DEPLOYMENT.md`
- `TODO.md`

## Required Prep Work

These are the concrete tasks another agent should complete before final OracleHost wiring.

### 1. Add OracleHost-compatible health endpoint

OracleHost deploy verification checks:

```text
/api/health
```

Current Volo backend exposes:

```text
/health
```

Current code:

- `volo-api/cmd/server/main.go` registers `r.Get("/health", handlers.Health)`

Required action:

- add `GET /api/health`
- keep `/health` if you want backwards compatibility

Safest path:

- expose both `/health` and `/api/health`

### 2. Add OracleHost-specific GitHub Actions workflow

OracleHost expects:

- a GitHub Actions workflow
- Linux backend image build
- push to GHCR
- deploy over SSH through `deploy-oracle-app`

Required action:

- add `.github/workflows/deploy-backend-oracle.yml`
- use a production GitHub environment for the hosted deploy
- make the workflow backend-only
- deploy through the narrow deploy user, not a broad shell user

### 3. Stop assuming the standalone `/opt/volo` layout

The current generic deployment doc assumes:

- `/opt/volo`
- app-owned Caddy
- app-owned Postgres and Redis containers
- manual `git pull`
- manual `scp` of landing assets

That is not the OracleHost model.

Required action:

- do not build OracleHost integration around `/opt/volo`
- target `/opt/apps/<app-id>` and `/etc/apps/<app-id>`
- let platform provisioning own the app slot, env file, and public route

### 4. Decide the Linux ONNX story

Current state:

- `volo-api/internal/intent/onnx_stub.go` is a stub
- `TODO.md` already notes Linux ONNX deployment work

Required action if server-side ML inference is part of the first deploy:

- implement the Linux ONNX runtime path
- ensure the chosen library works for the target Linux architecture
- ensure model artifacts are present in the deployed runtime

If the first OracleHost deployment does not require ML inference, document that intentionally so the fallback behavior is explicit.

### 5. Finalize model artifact strategy

The API can read model files via `VOLO_MODEL_PATH`, but deployment still needs a clear model artifact strategy for:

- `intent.onnx`
- tokenizer files
- model config files

Required action:

- decide whether model assets live inside the backend image or in mounted storage
- document the chosen path in the deployment workflow and runtime docs

Simple default:

- bake model assets into the backend image for the first deploy

### 6. Choose frontend placement intentionally

Recommended OracleHost pattern:

- frontend on Cloudflare Pages or Workers
- backend API on OracleHost

Required action:

- decide whether `volo-landing` is deployed to Cloudflare or handled separately

Recommended first pass:

- deploy `volo-landing` to Cloudflare Pages
- keep the Go API on OracleHost

### 7. Finalize production CORS inputs

Production CORS will need:

- the real frontend origin
- the real published browser extension ID

Required action:

- determine the final production extension ID
- decide whether development uses a different extension ID
- lock CORS to real origins before public rollout

## Environment Variables Volo Should Support

Core backend:

- `VOLO_PORT`
- `VOLO_DATABASE_URL`
- `VOLO_REDIS_ADDR`
- `VOLO_REDIS_PASSWORD`
- `VOLO_JWT_SECRET`
- `VOLO_CORS_ORIGINS`
- `VOLO_LOG_LEVEL`
- `VOLO_ENVIRONMENT`

Optional but likely needed:

- `VOLO_GOOGLE_CLIENT_ID`
- `VOLO_GOOGLE_CLIENT_SECRET`
- `VOLO_SENTRY_DSN`
- `VOLO_MODEL_PATH`

Recommended shape by environment:

### Production

```text
VOLO_PORT=<production-port>
VOLO_REDIS_ADDR=<redis-address-or-service-name>
VOLO_LOG_LEVEL=info
VOLO_ENVIRONMENT=production
VOLO_CORS_ORIGINS=https://<frontend-domain>,chrome-extension://<prod-extension-id>
VOLO_MODEL_PATH=/app/models
```

Important:

- do not commit real secret values
- do not commit real hostnames, IPs, SSH usernames, or private routing details

## GitHub Environment Shape

For this project, one hosted GitHub environment is enough:

- `production`

Recommended GitHub variables:

### Production

```text
ORACLE_HOST=<private-infra-value>
ORACLE_USER=<deploy-user>
ORACLE_APP_ID=volo
ORACLE_API_DOMAIN=<production-api-domain>
ORACLE_PORT=<production-port>
GHCR_IMAGE=ghcr.io/<owner>/<repo-or-image>
```

Recommended GitHub secrets:

- `ORACLE_SSH_KEY`
- `GHCR_TOKEN` if the image is private

## Domain Shape

Recommended public domain layout:

- frontend: `volo.petricevicsystems.com`
- API: `api.volo.petricevicsystems.com`

This will work even if `petricevicsystems.com` is not currently serving a website, as long as:

- you own the domain
- DNS is configured
- Cloudflare points the frontend and API records to the right targets

Using subdomains is cleaner than putting the app under a path such as `petricevicsystems.com/volo`.

## Suggested Backend Container Requirements

To fit OracleHost cleanly, the backend container should:

- listen on the runtime port env var
- expose `GET /api/health`
- include or mount model files if ML inference is enabled
- start without assuming a public-facing Caddy container inside the same app stack
- avoid exposing database services publicly

## Suggested Implementation Plan For The Other Agent

1. Add `GET /api/health` without removing the current `/health`.
2. Decide whether the first OracleHost deploy includes ONNX inference or intentionally skips it.
3. Add the OracleHost backend deployment workflow for GHCR plus SSH deploy.
4. Add any deployment-specific compose or runtime files needed for app-slot execution.
5. Lock production CORS to the real frontend and extension origins.
6. Keep all final live values in GitHub environments or private infra notes.
7. Hand off the repo once code and CI match the OracleHost contract.

## Private Wiring Checklist

These are the private values that still need to be inserted during final wiring:

- `ORACLE_HOST`
- `ORACLE_USER`
- `ORACLE_APP_ID`
- `ORACLE_API_DOMAIN`
- `ORACLE_PORT`
- `GHCR_IMAGE`
- `ORACLE_SSH_KEY`
- `VOLO_DATABASE_URL`
- `VOLO_REDIS_ADDR`
- `VOLO_REDIS_PASSWORD` if Redis is password-protected
- `VOLO_JWT_SECRET`
- `VOLO_CORS_ORIGINS`
- `VOLO_GOOGLE_CLIENT_ID` if Google OAuth is enabled
- `VOLO_GOOGLE_CLIENT_SECRET` if Google OAuth is enabled
- `VOLO_SENTRY_DSN` if Sentry is enabled

## What Stays Private

Do not publish these values in the repo:

- real server IPs
- real hostnames
- deploy usernames
- live port allocations
- current app inventory
- route inventory
- backup schedules
- state snapshots
- env file contents

Those belong in private infrastructure documentation and are only needed during final wiring.

## Bottom Line

The goal of this guide is not to fully describe a live server.

The goal is to make Volo deployable by contract.

If another agent completes the prep work in this document, final OracleHost wiring can be done later using private infrastructure details without changing the public repo again.
