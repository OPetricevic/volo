# ─────────────────────────────────────────────────────────────
# Volo — One-Command Development Setup (Windows)
#
# Run: .\setup.ps1
#
# This script:
#   1. Installs all npm dependencies (extension, desktop, landing)
#   2. Starts Postgres, Redis, API, Ollama via Docker Compose
#   3. Runs the database migration
#   4. Builds the extension
#   5. Prints instructions for final steps
# ─────────────────────────────────────────────────────────────

Write-Host ""
Write-Host "═══════════════════════════════════════════" -ForegroundColor Cyan
Write-Host "  Volo — Development Environment Setup" -ForegroundColor Cyan
Write-Host "═══════════════════════════════════════════" -ForegroundColor Cyan
Write-Host ""

# Check prerequisites
$missing = @()
if (-not (Get-Command docker -ErrorAction SilentlyContinue)) { $missing += "Docker" }
if (-not (Get-Command node -ErrorAction SilentlyContinue)) { $missing += "Node.js" }
if (-not (Get-Command go -ErrorAction SilentlyContinue)) { $missing += "Go" }

if ($missing.Count -gt 0) {
    Write-Host "Missing prerequisites: $($missing -join ', ')" -ForegroundColor Red
    Write-Host "Install them and re-run this script." -ForegroundColor Red
    exit 1
}

Write-Host "[1/5] Installing npm dependencies..." -ForegroundColor Yellow

Set-Location volo-extension
npm install --silent 2>$null
Write-Host "  ✓ Extension" -ForegroundColor Green

Set-Location ..\volo-desktop
npm install --silent 2>$null
Write-Host "  ✓ Desktop App" -ForegroundColor Green

Set-Location ..\volo-landing
npm install --silent 2>$null
Write-Host "  ✓ Landing Page" -ForegroundColor Green

Set-Location ..

Write-Host ""
Write-Host "[2/5] Starting Docker services (Postgres, Redis, API, Ollama)..." -ForegroundColor Yellow
docker compose -f docker-compose.dev.yml up -d --build 2>$null
Write-Host "  ✓ Services started" -ForegroundColor Green

Write-Host ""
Write-Host "[3/5] Waiting for database..." -ForegroundColor Yellow
Start-Sleep -Seconds 5

# Run migration for macros table
Write-Host "[4/5] Running database migrations..." -ForegroundColor Yellow
docker exec -i $(docker compose -f docker-compose.dev.yml ps -q postgres) psql -U volo -d volo -f /docker-entrypoint-initdb.d/001_initial.sql 2>$null
# Apply macros migration
Get-Content volo-api\migrations\002_macros.sql | docker exec -i $(docker compose -f docker-compose.dev.yml ps -q postgres) psql -U volo -d volo 2>$null
Write-Host "  ✓ Migrations applied" -ForegroundColor Green

Write-Host ""
Write-Host "[5/5] Building extension..." -ForegroundColor Yellow
Set-Location volo-extension
npm run build 2>$null
Set-Location ..
Write-Host "  ✓ Extension built (dist/)" -ForegroundColor Green

Write-Host ""
Write-Host "═══════════════════════════════════════════" -ForegroundColor Green
Write-Host "  Setup complete!" -ForegroundColor Green
Write-Host "═══════════════════════════════════════════" -ForegroundColor Green
Write-Host ""
Write-Host "Services running:" -ForegroundColor White
Write-Host "  • API:      http://localhost:8080/health" -ForegroundColor Gray
Write-Host "  • Postgres:  localhost:5432 (user: volo, pass: volo)" -ForegroundColor Gray
Write-Host "  • Redis:     localhost:6379" -ForegroundColor Gray
Write-Host "  • Ollama:    http://localhost:11434" -ForegroundColor Gray
Write-Host ""
Write-Host "Next steps:" -ForegroundColor White
Write-Host "  1. Load extension: chrome://extensions → Load unpacked → volo-extension\dist" -ForegroundColor Gray
Write-Host "  2. Desktop app:    cd volo-desktop; npx electron-vite dev" -ForegroundColor Gray
Write-Host "  3. Landing page:   cd volo-landing; npm run dev" -ForegroundColor Gray
Write-Host "  4. Pull AI model:  docker exec -it <ollama-container> ollama pull gemma2:2b" -ForegroundColor Gray
Write-Host ""
