# Vercel Function Migration Summary

## Overview

This document summarizes the changes made to deploy Memos to Vercel using serverless functions while maintaining all existing business logic.

## Changes Made

### 1. New Files Created

#### Backend Go Code
- `api/handler.go` - Main Vercel function entry point
- `server/vercel/proxy.go` - Echo to Vercel HTTP adapter
- `server/vercel/db_init.go` - Database initialization for serverless environment

#### Configuration Files
- `vercel.json` - Vercel deployment configuration
- `.env.vercel.example` - Environment variables template
- `VERCEL_DEPLOYMENT.md` - Comprehensive deployment guide
- `deploy-vercel.sh` - Linux/Mac deployment script
- `deploy-vercel.ps1` - Windows deployment script

### 2. Architecture Changes

**Before:**
- Single Go server listening on port 8081
- Echo HTTP server with Connect RPC and gRPC-Gateway
- Long-running background processes (S3 presign runner)

**After:**
- Vercel Serverless Functions handling API requests
- Single entry point (`api/handler.go`) routes all requests
- Echo server wrapped for Vercel compatibility
- Database connections use environment variables for cloud databases

### 3. Key Features Maintained

✅ All API endpoints (REST and Connect RPC)
✅ Authentication (JWT V2, Personal Access Tokens)
✅ Database migrations
✅ Request/response handling
✅ CORS configuration
✅ Error handling and logging

### 4. API Endpoints

All endpoints route through the same handler:

**REST API (gRPC-Gateway):**
- `/api/v1/auth/*` - Authentication
- `/api/v1/users/*` - User management
- `/api/v1/memos/*` - Memo operations
- `/api/v1/attachments/*` - File attachments
- `/api/v1/shortcuts/*` - Shortcuts
- `/api/v1/activities/*` - Activity tracking
- `/api/v1/idp/*` - Identity providers

**Connect RPC:**
- `/memos.api.v1.*` - All services via Connect RPC

**Other:**
- `/healthz` - Health check
- `/file/*` - Static file serving
- `/rss/*` - RSS feeds

### 5. Database Support

Three database drivers supported via environment variables:

1. **PostgreSQL** (Recommended for Vercel)
   ```
   MEMOS_DRIVER=postgres
   MEMOS_DSN=postgres://user:pass@host:port/db?sslmode=require
   ```

2. **MySQL**
   ```
   MEMOS_DRIVER=mysql
   MEMOS_DSN=user:pass@tcp(host:port)/db?parseTime=true
   ```

3. **SQLite**
   ```
   MEMOS_DRIVER=sqlite
   MEMOS_DSN=file:/tmp/memos_prod.db
   ```

### 6. Frontend Compatibility

The frontend code (`web/src/connect.ts`) already uses `window.location.origin` as the base URL, so it automatically adapts to any deployment domain including Vercel. No changes were needed.

### 7. Deployment Process

**Quick Start:**
```bash
# Install Vercel CLI
npm i -g vercel

# Deploy
./deploy-vercel.sh  # or deploy-vercel.ps1 on Windows
```

**Manual Steps:**
1. Build frontend: `cd web && pnpm install && pnpm build`
2. Run: `vercel`
3. Configure environment variables in Vercel Dashboard
4. Redeploy: `vercel --prod`

## Limitations & Considerations

1. **No Local Storage**: Vercel serverless functions are stateless. Use external databases only.

2. **60s Timeout**: Vercel functions have a 60-second timeout limit. Optimize long-running queries.

3. **Background Processes**: The S3 presign runner needs to be replaced with Vercel Cron Jobs or similar.

4. **File Storage**: File uploads need external storage (Vercel Blob, AWS S3, etc.) instead of local filesystem.

5. **Cold Starts**: First request may take longer due to function cold start.

6. **Connection Pooling**: Database connections are created per function invocation, consider connection pooling libraries.

## Testing

**Local Testing:**
```bash
# Set environment variables
export MEMOS_MODE=dev
export MEMOS_DRIVER=sqlite
export MEMOS_DSN=./memos_dev.db

# Run locally (standard Go server)
go run ./cmd/memos --mode dev --port 8081

# Or test the Vercel handler
cd api
go run handler.go
```

**Vercel Preview:**
```bash
vercel
```

**Production Deployment:**
```bash
vercel --prod
```

## Environment Variables

Required for Vercel:

```bash
MEMOS_MODE=prod
MEMOS_DRIVER=postgres          # or mysql/sqlite
MEMOS_DSN=postgres://...      # Your database connection
MEMOS_INSTANCE_URL=https://... # Optional, auto-configured
```

## Rollback

If you need to revert to the original setup:

1. Delete new files:
   - `api/` directory
   - `server/vercel/` directory
   - `vercel.json`
   - `deploy-vercel.*` files
   - `.env.vercel.example`
   - `VERCEL_DEPLOYMENT.md`

2. Use the original entry point:
   ```bash
   go run ./cmd/memos --mode prod --port 8081
   ```

## Support

For detailed deployment instructions, see `VERCEL_DEPLOYMENT.md`.

For database-specific setup:
- **Vercel Postgres**: https://vercel.com/docs/storage/vercel-postgres
- **Neon**: https://neon.tech/docs
- **PlanetScale**: https://docs.planetscale.com
- **Turso**: https://docs.turso.tech
