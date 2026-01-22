# Memos Codebase Guide for AI Agents

## Project Overview

Memos is a self-hosted knowledge management platform: Go 1.25 (gRPC + Connect RPC), React 18.3 + TypeScript + Vite 7, SQLite/MySQL/PostgreSQL, Protocol Buffers v2.

## Build/Lint/Test Commands

### Backend (Go)
```bash
# Run dev server
go run ./cmd/memos --mode dev --port 8081

# Run all tests
go test ./...

# Run tests for specific package
go test ./store/...
go test ./server/router/api/v1/test/...

# Run single test function
go test -run TestFunctionName ./path/to/package

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
go test -cover ./...

# Lint (golangci-lint)
golangci-lint run

# Format imports
goimports -w .
```

### Frontend (TypeScript/React)
```bash
cd web && pnpm install

# Dev server (proxies API to localhost:8081)
pnpm dev

# Type checking
pnpm lint

# Auto-fix lint issues
pnpm lint:fix

# Format code
pnpm format

# Build for production
pnpm build

# Build and copy to backend
pnpm release
```

### Protocol Buffers
```bash
cd proto && buf generate    # Regenerate Go & TypeScript
cd proto && buf lint        # Lint proto files
cd proto && buf breaking --against .git#main
```

## Code Style Guidelines

### Go

**Error Handling:**
- Wrap errors with `github.com/pkg/errors`: `errors.Wrap(err, "context")` or `errors.Wrapf(err, "format %s", arg)`
- Return gRPC errors: `status.Errorf(codes.NotFound, "message")`
- Forbidden: `fmt.Errorf` - use errors.Wrap/Wrapf/Errorf instead

**Naming:**
- Package names: lowercase, single word (e.g., `store`, `server`)
- Interfaces: `Driver`, `Store`, `Service` - PascalCase
- Methods: PascalCase for exported, camelCase for internal
- Constants: PascalCase

**Comments:**
- Public exported functions must have godoc comments (enforced by godot)
- Single-line: `// Comment`, Multi-line: `/* Comment */`

**Imports:**
- Group: stdlib → third-party → local (separated by blank lines)
- Sorted alphabetically within groups
- Format with `goimports -w .`
- Forbidden: `ioutil.ReadDir` - use `os.ReadDir` instead

**Linter:** golangci-lint with revive, govet, staticcheck, misspell, gocritic, sqlclosecheck, rowserrcheck, nilerr, godot, forbidigo, mirror, bodyclose

### TypeScript/React

**Components:**
- Functional components with hooks only
- Props interface: `interface Props { ... }`
- Use `useMemo`, `useCallback` for optimization

**State Management:**
- Server state: React Query hooks (`web/src/hooks/`)
- Client state: React Context (`AuthContext`, `ViewContext`, `MemoFilterContext`)
- Avoid `useState` for server data

**Types:**
- Use explicit interfaces for complex objects
- No `any` - use `unknown` or specific types
- Use `as const` for literal types

**Formatting (Biome):**
- Line width: 140 characters
- Semicolons: always
- Quotes: double
- Indent: 2 spaces
- Trailing commas: all

**Imports:**
- Absolute imports with `@/` alias
- Auto-organized by Biome

**Styling:**
- Tailwind CSS v4 via `@tailwindcss/vite`
- Use `clsx` and `tailwind-merge` for conditional classes

## Key Architecture Patterns

**API Layer:** Dual protocol - Connect RPC (browsers, `/memos.api.v1.*`) + gRPC-Gateway (REST, `/api/v1/*`). Same service implementations for both.

**Authentication:** JWT Access Tokens (V2, 15-min) and Personal Access Tokens (PAT, long-lived). Both use `Authorization: Bearer <token>` header.

**Store Layer:** All database operations through `Driver` interface. Three implementations: SQLite, MySQL, PostgreSQL. In-memory caching for instance settings, users, user settings.

**Frontend State:** React Query v5 for server state (staleTime: 30s, gcTime: 5min), React Context for client state.

**Migrations:** Schema version in `instance_setting` table (key: `bb.general.version`). Migration files: `store/migration/{driver}/{version}/NN__description.sql`

## Common Tasks

**Add API endpoint:** Edit `proto/api/v1/*_service.proto` → `cd proto && buf generate` → implement in `server/router/api/v1/*_service.go` → if public, add to `acl_config.go`

**DB schema change:** Create migration files for all drivers → update `LATEST.sql` → update store interface if needed → test with `go test ./store/test/...`

**Run tests against specific DB:**
```bash
DRIVER=sqlite go test ./...
DRIVER=mysql DSN="user:pass@tcp(localhost:3306)/memos" go test ./...
DRIVER=postgres DSN="postgres://user:pass@localhost:5432/memos" go test ./...
```

**Test database:** `store/test/store.go` provides `NewTestingStore()` for isolated DB, `resetTestingDB()` to clean tables

## Important Files

| File | Purpose |
|------|---------|
| `cmd/memos/main.go` | Server entry point |
| `server/router/api/v1/v1.go` | Service registration |
| `store/driver.go` | Database driver interface |
| `store/migrator.go` | Migration logic |
| `web/src/lib/query-client.ts` | React Query config |
| `web/src/hooks/useMemoQueries.ts` | Memo queries/mutations |
| `web/src/contexts/` | React Context providers |

## Environment Variables

Backend: `MEMOS_MODE` (dev/prod/demo), `MEMOS_PORT` (8081), `MEMOS_DRIVER` (sqlite/mysql/postgres), `MEMOS_DATA` (~/.memos)

Frontend: `DEV_PROXY_SERVER` (http://localhost:8081)
