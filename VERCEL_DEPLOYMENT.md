# Vercel Deployment Guide for Memos

This guide explains how to deploy Memos to Vercel using serverless functions.

## Prerequisites

- Vercel account
- Cloud database (PostgreSQL, MySQL, or SQLite-compatible)
- Node.js 18+ for frontend
- Go 1.25+ for local testing

## Database Options

### Vercel Postgres (Recommended)
```
MEMOS_DRIVER=postgres
MEMOS_DSN=postgres://user:password@host:port/database?sslmode=require
```

### Neon PostgreSQL
```
MEMOS_DRIVER=postgres
MEMOS_DSN=postgresql://user:password@ep-cool-name.us-east-2.aws.neon.tech/neondb?sslmode=require
```

### PlanetScale MySQL
```
MEMOS_DRIVER=mysql
MEMOS_DSN=user:password@tcp(host:port)/database?parseTime=true
```

### Turso SQLite
```
MEMOS_DRIVER=sqlite
MEMOS_DSN=file:/tmp/memos_prod.db
```

## Environment Variables

Required environment variables for Vercel deployment:

```bash
# Server configuration
MEMOS_MODE=prod

# Database configuration
MEMOS_DRIVER=postgres    # sqlite, mysql, or postgres
MEMOS_DSN=postgres://... # Your database connection string

# Optional
MEMOS_INSTANCE_URL=https://your-app.vercel.app
```

## Deployment Steps

1. **Clone the repository**
   ```bash
   git clone https://github.com/your-org/memos.git
   cd memos
   ```

2. **Install frontend dependencies**
   ```bash
   cd web
   pnpm install
   pnpm build
   cd ..
   ```

3. **Deploy to Vercel**
   ```bash
   vercel
   ```

4. **Configure environment variables in Vercel Dashboard**
   - Go to Settings > Environment Variables
   - Add the required variables (see above)
   - Select appropriate environments (Production, Preview, Development)

5. **Redeploy**
   ```bash
   vercel --prod
   ```

## Architecture

### Backend (Go)
- Entry point: `api/handler.go`
- Vercel proxy: `server/vercel/`
- All API endpoints route through a single handler
- Supports both gRPC-Gateway (REST) and Connect RPC

### Frontend (React)
- Built with Vite
- Served from `/` path
- API calls go to `/api/v1/*` (gRPC-Gateway) or `/memos.api.v1.*` (Connect RPC)

### API Endpoints

All backend API endpoints are accessible via:

- **REST API (gRPC-Gateway)**: `https://your-app.vercel.app/api/v1/*`
- **Connect RPC**: `https://your-app.vercel.app/memos.api.v1.*`

Services available:
- Auth Service: `/api/v1/auth/*`
- User Service: `/api/v1/users/*`
- Memo Service: `/api/v1/memos/*`
- Attachment Service: `/api/v1/attachments/*`
- Shortcut Service: `/api/v1/shortcuts/*`
- Activity Service: `/api/v1/activities/*`
- Identity Provider Service: `/api/v1/idp/*`

## Local Development

To test locally:

1. **Set environment variables**
   ```bash
   export MEMOS_MODE=dev
   export MEMOS_DRIVER=sqlite
   export MEMOS_DSN=./memos_dev.db
   ```

2. **Run backend**
   ```bash
   go run ./cmd/memos --mode dev --port 8081
   ```

3. **Run frontend**
   ```bash
   cd web
   pnpm dev
   ```

4. **Access at**: http://localhost:5173

## Limitations

- No persistent local storage in Vercel (use external database)
- Connection pooling is handled by Vercel's Go runtime
- Background runners (S3 presign, etc.) may need adjustment for serverless
- File uploads may need storage service integration (e.g., Vercel Blob, AWS S3)

## Troubleshooting

### Database connection errors
- Verify DSN format for your database type
- Check firewall rules allow Vercel's IP ranges
- Ensure SSL is enabled for cloud databases

### Timeouts
- Vercel Serverless Functions have 60s timeout limit
- Optimize long-running queries
- Consider using Vercel Cron jobs for background tasks

### Cold starts
- First request may take longer
- Consider using Vercel's Edge Runtime for frequently accessed endpoints
- Enable Keep-Alive headers in database connections

## Scaling

Vercel automatically scales based on traffic:
- 0 instances when idle
- Scales up automatically with traffic
- No manual configuration needed
- Pay only for actual usage
