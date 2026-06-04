# ─────────────────────────────────────────────────────────────
# Volo — Development Commands
#
# Usage:
#   make setup    — First time: install deps + start everything
#   make start    — Start all services (Postgres, Redis, API, Ollama)
#   make stop     — Stop all services
#   make test     — Run all tests across the project
#   make build    — Build everything (API, extension, desktop, landing)
#   make ext      — Build extension only
#   make desktop  — Run desktop app in dev mode
#   make landing  — Run landing page dev server
#   make clean    — Stop services and remove volumes
# ─────────────────────────────────────────────────────────────

.PHONY: setup start stop test build ext desktop landing clean status

# First-time setup: install all deps + start services + build extension
setup:
	@echo ═══ Volo Setup ═══
	@echo [1/4] Installing dependencies...
	cd volo-extension && npm install --silent
	cd volo-desktop && npm install --silent
	cd volo-landing && npm install --silent
	@echo [2/4] Starting services...
	docker compose -f docker-compose.dev.yml up -d --build
	@echo [3/4] Running migrations...
	@timeout /t 5 /nobreak >nul 2>&1 || sleep 5
	docker compose -f docker-compose.dev.yml exec -T postgres psql -U volo -d volo -f /docker-entrypoint-initdb.d/001_initial.sql >nul 2>&1 || true
	type volo-api\migrations\002_macros.sql | docker compose -f docker-compose.dev.yml exec -T postgres psql -U volo -d volo >nul 2>&1 || true
	@echo [4/4] Building extension...
	cd volo-extension && npm run build
	@echo ✓ Done! API at http://localhost:8080/health
	@echo ✓ Load extension: chrome://extensions → Load unpacked → volo-extension\dist
	@echo ✓ Desktop: make desktop
	@echo ✓ Landing: make landing

# Start all backend services
start:
	docker compose -f docker-compose.dev.yml up -d

# Stop all services
stop:
	docker compose -f docker-compose.dev.yml down

# Run all tests
test:
	@echo ═══ API Tests ═══
	cd volo-api && go test -short ./...
	@echo ═══ Extension Tests ═══
	cd volo-extension && npx vitest run
	@echo ═══ Desktop Tests ═══
	cd volo-desktop && npx vitest run

# Build everything
build:
	cd volo-api && go build ./...
	cd volo-extension && npm run build
	cd volo-desktop && npx electron-vite build
	cd volo-landing && npm run build

# Build extension only
ext:
	cd volo-extension && npm run build

# Run desktop app
desktop:
	cd volo-desktop && npx electron-vite dev

# Run landing page dev server
landing:
	cd volo-landing && npm run dev

# Show service status
status:
	docker compose -f docker-compose.dev.yml ps

# Full cleanup (removes data volumes too)
clean:
	docker compose -f docker-compose.dev.yml down -v
	@echo Volumes removed. Next "make setup" starts fresh.
