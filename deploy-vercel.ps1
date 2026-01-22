# Deployment script for Vercel (PowerShell)

Write-Host "Deploying Memos to Vercel..." -ForegroundColor Green

# Build frontend
Write-Host "Building frontend..." -ForegroundColor Yellow
Set-Location web
pnpm install
pnpm build
Set-Location ..

# Deploy to Vercel
Write-Host "Deploying to Vercel..." -ForegroundColor Yellow
vercel

Write-Host "Deployment complete!" -ForegroundColor Green
Write-Host "Set up environment variables in Vercel Dashboard:" -ForegroundColor Cyan
Write-Host "  - MEMOS_MODE=prod" -ForegroundColor White
Write-Host "  - MEMOS_DRIVER=postgres (or mysql/sqlite)" -ForegroundColor White
Write-Host "  - MEMOS_DSN=your-database-connection-string" -ForegroundColor White
Write-Host "  - MEMOS_INSTANCE_URL=https://your-app.vercel.app" -ForegroundColor White
