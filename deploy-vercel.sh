#!/bin/bash
set -e

echo "Deploying Memos to Vercel..."

# Build frontend
echo "Building frontend..."
cd web
pnpm install
pnpm build
cd ..

# Deploy to Vercel
echo "Deploying to Vercel..."
vercel

echo "Deployment complete!"
echo "Set up environment variables in Vercel Dashboard:"
echo "  - MEMOS_MODE=prod"
echo "  - MEMOS_DRIVER=postgres (or mysql/sqlite)"
echo "  - MEMOS_DSN=your-database-connection-string"
echo "  - MEMOS_INSTANCE_URL=https://your-app.vercel.app"
