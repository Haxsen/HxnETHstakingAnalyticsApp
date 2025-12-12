# Simplified Backend Refactor Plan

## Current Problems
- **400+ line main.go** with everything mixed together
- **Global variables** (`coingeckoClient`)
- **Handlers doing too much** (validation + business logic + HTTP)
- **Repeated error handling** patterns
- **Empty `internal/api/`** directory

## Simplified Solution (3 Phases)

### Phase 1: Extract Handlers (2 hours)
**Goal**: Move HTTP handlers out of main.go

**Files to create**:
- `internal/api/handlers.go` - All HTTP handlers
- `internal/api/responses.go` - Common response helpers

**Benefits**:
- main.go shrinks from 400+ to ~50 lines
- Handlers grouped logically
- Easier to test handlers separately

### Phase 2: Add Server Struct (1-2 hours)
**Goal**: Replace global variables with dependency injection

**Files to create**:
- `internal/server/server.go` - Server struct with injected dependencies

**Benefits**:
- No more global state
- Easier testing with mocks
- Clear dependency management

### Phase 3: Extract Business Logic (2-3 hours)
**Goal**: Separate business logic from HTTP concerns

**Files to create**:
- `internal/services/token_service.go` - Token business logic
- `internal/services/valuation_service.go` - Valuation calculations

**Benefits**:
- Business logic reusable
- Thinner HTTP handlers
- Logic can be tested independently

## Final Structure
```
backend/
├── main.go                    # Entry point (~50 lines)
├── internal/
│   ├── api/
│   │   ├── handlers.go        # HTTP handlers
│   │   └── responses.go       # Response helpers
│   ├── server/
│   │   └── server.go          # Server with DI
│   └── services/              # Business logic
│       ├── token_service.go
│       └── valuation_service.go
```

## Why This Approach?
- **80% of benefits** with 50% of effort
- **Backward compatible** - no API changes
- **Incremental** - can stop after any phase
- **Focused** - solves the biggest pain points first

## Timeline: 5-7 hours total
- Much faster than the original 11-16 hour plan
- Each phase delivers immediate value
- Can be done in one sitting per phase

## Ready to start with Phase 1?
This simplified plan focuses on the highest-impact changes while keeping things manageable.
