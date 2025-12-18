# Tools Hub — Go Mastery Roadmap

## Mission Statement

**Primary Goal**: Achieve mastery of the Go programming language through deliberate practice of its core concepts, patterns, and idioms.

**Secondary Goal**: Build a useful repo-aware planning agent as a capstone project.

This roadmap is optimized for **learning Go deeply** rather than shipping features quickly. Every phase is designed to teach fundamental Go concepts that compound into mastery.

---

## Go Mastery Learning Philosophy

### Core Principles
1. **Concurrency First** - Go's differentiating feature deserves early, deep practice
2. **Interfaces Over Concrete Types** - Think in abstractions
3. **Test Everything** - TDD is how you internalize patterns
4. **Profile and Optimize** - Fast code requires measurement
5. **Read the Standard Library** - Best Go code you'll ever see
6. **Build Real Things** - Theory without practice is worthless

### Mastery Metrics
By the end of this roadmap, you should be able to:
- Design and implement concurrent systems with goroutines and channels
- Create elegant APIs using interface composition
- Write idiomatic, production-ready Go code
- Profile and optimize performance bottlenecks
- Handle errors gracefully across complex systems
- Build scalable web services with proper architecture
- Read and understand any Go codebase in the wild

---

## Phase 1 — Foundation & Web Basics (Days 1-15) ✅

### Days 1-9 — HTTP Server Foundation (Completed)
You've already completed:
- HTTP server + middleware
- Template system (base + partials + pages)
- File upload handling
- Routing patterns
- Template data contracts
- Two working tools (slugify, file converter)
- Basic unit tests for each tool (table-driven)
- One handler test per tool route

**Key Learning**: Basic HTTP patterns, template rendering, form handling

### Days 10-12 — JSON Formatter Tool
- `/tools/json` GET/POST
- Internal `jsonutil` package
- Pretty-print and minify modes
- Table-driven tests for `jsonutil`
- Handler test for `/tools/json`
- **Key Learning**: JSON marshaling, error handling patterns, form state preservation

### Days 13-15 — Base64 Encoder/Decoder
- `/tools/base64` GET/POST
- Internal `encodingutil` package
- Encode/decode modes
- Table-driven tests for `encodingutil`
- Handler test for `/tools/base64`
- **Key Learning**: Binary data handling, encoding packages, io.Reader patterns

---

## Phase 2 — Concurrency Foundations (Days 16-30) 🔥

**⚠️ MOST CRITICAL PHASE FOR GO MASTERY**

### Days 16-18 — Goroutines & WaitGroups
**Complexity: Medium** | **Priority: CRITICAL**

Build: **Concurrent File Processor**
- Upload multiple files (3-5 at once)
- Process each file in a separate goroutine
- Use `sync.WaitGroup` to wait for all to complete
- Return all results together

**Core Concepts**:
- Starting goroutines with `go` keyword
- Anonymous functions and closures
- Race conditions (introduce them, then fix them)
- `sync.WaitGroup` for coordination
- Memory visibility issues

**Exercises**:
1. Process files sequentially (baseline)
2. Process files concurrently (see the speedup)
3. Introduce a bug (shared variable race)
4. Fix with proper synchronization
5. Run with `-race` detector

**Key Learning**: How goroutines work, when concurrency helps, how to coordinate

**Testing Thread** (starts here and continues daily):
- Add unit tests for new helpers/util packages the day you create them
- Add at least one handler test for any new tool route
- Use `-race` during concurrency exercises

### Days 19-21 — Channels Fundamentals
**Complexity: Medium** | **Priority: CRITICAL**

Build: **Concurrent Hash Calculator**
- Upload multiple files
- Hash each in a goroutine
- Send results back via channels
- Display all hashes when complete

**Core Concepts**:
- Unbuffered channels (blocking semantics)
- Buffered channels (capacity and blocking)
- Channel direction (`<-chan`, `chan<-`)
- Closing channels
- Range over channels
- Select statement for multiplexing

**Exercises**:
1. One channel per goroutine (messy)
2. One shared channel (cleaner)
3. Bidirectional channel patterns
4. Channel closing and range
5. Select with timeout

**Key Learning**: Channels as communication mechanism, CSP model

### Days 22-24 — Worker Pool Pattern
**Complexity: Medium-High** | **Priority: CRITICAL**

Build: **Concurrent Image Processor**
- Upload 10+ images
- Process with a pool of 3 workers
- Workers pull from job queue
- Progress tracking via channel

**Core Concepts**:
- Bounded concurrency (semaphore pattern)
- Job queue with buffered channel
- Worker goroutines
- Result collection patterns
- Graceful shutdown

**Pattern Template**:
```go
type Job struct {
    ID   int
    Data []byte
}

type Result struct {
    JobID int
    Output []byte
    Error error
}

// jobs channel feeds workers
// results channel collects outputs
// done channel signals completion
```

**Key Learning**: Production concurrency patterns, resource limiting

### Days 25-27 — Context Package
**Complexity: Medium** | **Priority: CRITICAL**

Build: **Cancellable Operations Tool**
- Long-running operation (e.g., large file processing)
- Cancel button that stops operation mid-stream
- Timeout after 30 seconds
- Progress reporting

**Core Concepts**:
- `context.Context` interface
- `context.WithCancel`
- `context.WithTimeout`
- `context.WithDeadline`
- `context.WithValue` (use sparingly)
- Passing context through call chains
- Checking `ctx.Done()` in loops

**Exercises**:
1. Operation without cancellation
2. Add manual cancellation
3. Add automatic timeout
4. Propagate context through helper functions
5. Add request-scoped tracing (context values)

**Key Learning**: Cancellation propagation, request lifecycle management

### Days 28-30 — Pipeline Pattern
**Complexity: High** | **Priority: HIGH**

Build: **CSV Processing Pipeline**
- Stage 1: Read and parse CSV
- Stage 2: Validate rows
- Stage 3: Transform data
- Stage 4: Filter results
- Stage 5: Output as JSON

**Core Concepts**:
- Fan-out, fan-in patterns
- Pipeline stages connected by channels
- Error handling in pipelines
- Backpressure management
- Pipeline shutdown

**Pattern Template**:
```go
func stage1(in <-chan Data) <-chan Result1 {
    out := make(chan Result1)
    go func() {
        defer close(out)
        for data := range in {
            // process
            out <- result
        }
    }()
    return out
}
```

**Key Learning**: Compositional concurrency, stream processing patterns

---

## Phase 3 — Testing & Quality (Days 31-42)

### Days 31-33 — Unit Testing Mastery
**Complexity: Medium** | **Priority: CRITICAL**

**Day 31**: Table-Driven Tests
- Refactor any remaining tests to table-driven style
- Learn idiomatic test structure
- Test naming conventions (`TestFunction_Scenario_ExpectedOutcome`)

**Day 32**: Test Helpers and Fixtures
- Create reusable test helpers
- Golden files for complex outputs
- Test data generation
- Subtests with `t.Run()`

**Day 33**: Testing Concurrent Code
- Test goroutines and channels
- Timeout patterns in tests
- Race detector in CI
- Deadlock detection

**Key Learning**: Idiomatic Go testing, reproducible tests

### Days 34-36 — Benchmarking & Profiling
**Complexity: Medium-High** | **Priority: HIGH**

**Day 34**: Benchmarking Basics
- Write benchmarks for all tools
- `testing.B` and `b.N` loops
- Benchmark comparisons
- Avoid common pitfalls (compiler optimizations)

**Day 35**: CPU Profiling
- Profile your concurrent file processor
- Use `pprof` to find bottlenecks
- Flame graphs
- Optimize hot paths

**Day 36**: Memory Profiling
- Memory allocation profiling
- Escape analysis (`go build -gcflags="-m"`)
- Reduce allocations
- Buffer pooling with `sync.Pool`

**Key Learning**: Performance measurement, optimization techniques

### Days 37-39 — Integration & HTTP Testing
**Complexity: Medium** | **Priority: HIGH**

**Day 37**: HTTP Handler Testing
- Use `httptest.ResponseRecorder`
- Test full request/response cycle
- Test middleware chains
- Session and cookie testing

**Day 38**: Integration Tests
- Test entire tool workflows end-to-end
- Database integration tests (coming soon)
- External service mocking
- Test containers (optional)

**Day 39**: Fuzzing
- Add fuzzing tests (Go 1.18+)
- Fuzz your parsers (JSON, CSV, etc.)
- Find edge cases automatically
- Corpus management

**Key Learning**: Testing strategies beyond unit tests

### Days 40-42 — CI/CD & Deployment
**Complexity: Medium** | **Priority: HIGH**

**Day 40**: GitHub Actions
- Automated testing pipeline
- Multiple Go versions
- Coverage reporting
- Linting (golangci-lint)

**Day 41**: Docker & Containers
- Multi-stage Dockerfile
- Minimal runtime image
- Health checks
- Container security

**Day 42**: First Deployment
- Deploy to Fly.io/Render/Railway
- Environment configuration
- Logging and monitoring
- Custom domain (optional)

---

## Phase 4 — Interfaces & Design Patterns (Days 43-55)

### Days 43-46 — Interface Mastery
**Complexity: Medium-High** | **Priority: CRITICAL**

**Day 43**: Interface Basics
- Create a `Tool` interface
- Implement for existing tools
- Interface composition
- Empty interface and type assertions

**Day 44**: io.Reader and io.Writer Patterns
- Deep dive into `io` package
- `io.Reader`, `io.Writer`, `io.Closer`
- `io.Copy`, `io.Pipe`
- Chaining readers/writers
- Build custom Reader/Writer implementations

**Day 45**: Interface Design Principles
- Accept interfaces, return structs
- Small interfaces (1-3 methods)
- Implicit satisfaction
- Interface segregation
- Refactor existing code to better interfaces

**Day 46**: Advanced Interface Patterns
- Type switches and type assertions
- Interface embedding
- Polymorphism in Go
- Duck typing examples

**Key Learning**: Go's approach to abstraction and composition

### Days 47-50 — Error Handling Mastery
**Complexity: Medium** | **Priority: CRITICAL**

**Day 47**: Error Types and Wrapping
- Custom error types
- Error wrapping with `%w`
- `errors.Is()` and `errors.As()`
- Error chains

**Day 48**: Sentinel Errors
- Define package-level sentinel errors
- When to use vs custom types
- Error comparison patterns
- Migration from old error checking

**Day 49**: Error Handling Strategies
- Return errors vs panic
- Error context and stack traces
- Logging errors appropriately
- User-facing vs internal errors

**Day 50**: Panic and Recovery
- When to panic (never in library code)
- Recover in HTTP handlers
- Graceful degradation
- Panic/recover patterns

**Key Learning**: Production-grade error handling

### Days 51-55 — Generics (Go 1.18+)
**Complexity: Medium** | **Priority: MEDIUM**

**Day 51**: Type Parameters Basics
- Generic functions
- Type constraints
- Type inference
- When to use generics

**Day 52**: Generic Data Structures
- Generic slice helpers
- Generic map operations
- Generic cache implementation
- Performance considerations

**Day 53**: Constraints and Interfaces
- Built-in constraints (`any`, `comparable`)
- Custom constraints
- Interface constraints
- Constraint composition

**Day 54**: Refactoring with Generics
- Identify code duplication in your project
- Refactor to use generics where appropriate
- When NOT to use generics

**Day 55**: Advanced Generic Patterns
- Generic result types
- Generic option patterns
- Type lists
- Gotchas and limitations

**Key Learning**: Modern Go features, type-safe abstractions

---

## Phase 5 — Database & Persistence (Days 56-70)

### Days 56-60 — SQL Database Fundamentals
**Complexity: Medium-High** | **Priority: CRITICAL**

**Day 56**: database/sql Package
- Connect to SQLite
- Connection pooling
- Prepared statements
- Query vs Exec vs QueryRow

**Day 57**: CRUD Operations
- Build a tool history tracker
- Store tool usage in database
- Basic CRUD operations
- Transaction management

**Day 58**: Advanced Queries
- JOIN operations
- Subqueries
- Aggregations
- Pagination

**Day 59**: Transactions and Locking
- Begin/Commit/Rollback
- Isolation levels
- Optimistic locking
- Connection lifecycle

**Day 60**: Database Testing
- In-memory SQLite for tests
- Test fixtures and migrations
- Rollback tests
- Integration test patterns

**Key Learning**: Go's database abstraction, SQL in production

### Days 61-65 — Repository Pattern & Advanced DB
**Complexity: High** | **Priority: HIGH**

**Day 61**: Repository Pattern
- Interface-based repositories
- Dependency injection
- Mock repositories for testing
- Repository composition

**Day 62**: Query Builders
- Type-safe query building
- Optional: Try SQLC or sqlx
- Dynamic query generation
- SQL injection prevention

**Day 63**: Database Migrations
- Schema versioning
- Up/down migrations
- Migration testing
- Production migration strategies

**Day 64**: Connection Management
- Pool sizing
- Connection health checks
- Retry logic
- Circuit breakers

**Day 65**: Caching with Redis (Optional)
- Redis basics in Go
- Cache-aside pattern
- TTL management
- Cache invalidation

**Key Learning**: Production database patterns, data layer design

### Days 66-70 — Enhanced Tools with Persistence

**Build**:
- **User Preferences System**: Store user tool settings
- **Tool History**: Track all tool usage with timestamps
- **Saved Results**: Save and retrieve tool outputs
- **Statistics Dashboard**: Query and display usage metrics

**Key Learning**: Applying database knowledge to real features

---

## Phase 6 — Advanced HTTP & Web Patterns (Days 71-85)

### Days 71-75 — HTTP Client Mastery
**Complexity: Medium** | **Priority: HIGH**

**Day 71**: HTTP Client Basics
- `http.Client` configuration
- Custom timeouts
- Connection pooling
- Keep-alive

**Day 72**: Retry and Backoff
- Exponential backoff
- Retry policies
- Circuit breaker pattern
- Request context

**Day 73**: External API Integration
- Build a tool that calls external APIs
- Rate limiting client-side
- Response caching
- Error handling

**Day 74**: HTTP/2 and WebSockets
- HTTP/2 benefits
- WebSocket basics
- Build a real-time tool (e.g., live log viewer)
- Server-sent events

**Day 75**: Streaming Responses
- Stream large files
- Chunked transfer encoding
- Streaming JSON responses
- Flush patterns

**Key Learning**: Production HTTP client patterns

### Days 76-80 — Advanced Web Patterns
**Complexity: Medium-High** | **Priority: HIGH**

**Day 76**: Middleware Deep Dive
- Middleware composition patterns
- Request/response modification
- Middleware ordering
- Reusable middleware

**Day 77**: Rate Limiting
- Token bucket algorithm
- Sliding window
- Per-user rate limits
- Distributed rate limiting

**Day 78**: Request Validation
- Input sanitization
- Schema validation
- Error response formatting
- Validation middleware

**Day 79**: Authentication & Sessions
- Session management
- JWT basics
- Cookie security
- CSRF protection

**Day 80**: Security Hardening
- XSS prevention
- SQL injection prevention
- Secure headers
- HTTPS enforcement

**Key Learning**: Production web service patterns

### Days 81-85 — Observability
**Complexity: Medium** | **Priority: HIGH**

**Day 81**: Structured Logging
- `log/slog` package (Go 1.21+)
- Log levels and context
- Structured fields
- Log aggregation

**Day 82**: Metrics and Monitoring
- Prometheus metrics
- Custom metrics
- Health endpoints
- Readiness vs liveness

**Day 83**: Tracing
- Request tracing
- Context propagation
- Distributed tracing basics
- Performance debugging

**Day 84**: Graceful Shutdown
- Signal handling
- Connection draining
- Cleanup patterns
- Zero-downtime deploys

**Day 85**: Production Readiness
- Configuration management
- Environment variables
- Feature flags
- Deployment checklist

**Key Learning**: Operating Go services in production

---

## Phase 7 — Advanced Tools (Days 86-105)

### Days 86-90 — Text Processing Suite
**Complexity: Medium** | **Priority: MEDIUM**

**Day 86-87**: URL Encoder/Decoder
- Component encoding (path, query, fragment)
- Query parameter parsing
- Unicode handling

**Day 88-89**: Hashing Tool
- SHA-256, SHA-512, MD5
- HMAC support
- File hashing with streaming
- Concurrent hashing of multiple files

**Day 90**: Timestamp Converter
- Unix epoch ↔ human readable
- Multiple timezones
- Format parsing
- Duration calculations

### Days 91-95 — Data Transformation Tools
**Complexity: Medium** | **Priority: HIGH**

**Day 91-92**: CSV to JSON Converter
- Robust CSV parsing
- Header detection
- Type inference
- Streaming for large files

**Day 93-94**: Diff Tool
- Line-by-line comparison
- Side-by-side view
- File diff support
- Syntax highlighting (optional)

**Day 95**: Regex Tester
- Pattern testing
- Match highlighting
- Capture groups
- Common pattern library

### Days 96-100 — Advanced Text Tools
**Complexity: Medium** | **Priority: MEDIUM**

**Day 96-97**: Unicode Text Analyzer
- Bytes, runes, code points
- UTF-8 visualization
- Normalization forms
- Character frequency

**Day 98-99**: String Manipulation Suite
- Case transformations (camel, snake, kebab)
- Rune-aware operations
- Word count and analysis
- Palindrome detection

**Day 100**: Markdown Preview
- Markdown to HTML
- CommonMark + GFM
- Syntax highlighting
- XSS prevention

### Days 101-105 — File Operations
**Complexity: Medium** | **Priority: MEDIUM**

**Day 101-102**: Archive Creator (ZIP)
- Multiple file upload
- Safe path handling
- Concurrent archiving
- Size limits

**Day 103-104**: Archive Extractor
- ZIP inspection and extraction
- Zip-slip protection
- Selective extraction
- Preview contents

**Day 105**: File Inspector
- MIME type detection
- Metadata extraction
- Hex dump viewer
- File analysis

---

## Phase 8 — Repo Analysis Agent (Days 106-130) 🎯

**Prerequisites**: You now have mastered:
- ✅ Concurrency and goroutines
- ✅ Interfaces and design patterns
- ✅ Database and persistence
- ✅ HTTP clients and servers
- ✅ Testing and profiling
- ✅ Production patterns

### Days 106-110 — GitHub API Integration
**Complexity: High** | **Priority: CRITICAL**

**Day 106-107**: API Client Foundation
- GitHub REST API client
- OAuth implementation
- Personal Access Token support
- Rate limiting with context

**Day 108**: API Data Models
- Repository metadata
- File tree structure
- Commit history
- Issue/PR data

**Day 109**: Caching Strategy
- Response caching with Redis/SQLite
- Cache invalidation
- ETag support
- Offline mode

**Day 110**: Error Handling & Retry
- Network error handling
- Rate limit backoff
- Pagination handling
- Partial failure recovery

**Key Learning**: External API integration at scale

### Days 111-120 — File Tree Analysis Engine
**Complexity: High** | **Priority: CRITICAL**

**Day 111-113**: Tree Parsing
- Parse recursive tree API response
- Build in-memory file structure
- Concurrent file fetching
- Path normalization

**Day 114-116**: Pattern Detection
- Detect Go project structure (`cmd/`, `internal/`, `pkg/`)
- Identify frameworks (MVC, Clean Architecture, etc.)
- Route-to-handler mapping
- Test file detection

**Day 117-118**: Language Analysis
- Multi-language support (Go, Python, Node.js)
- Framework detection per language
- Convention analysis
- Best practice checking

**Day 119-120**: Performance Optimization
- Benchmark analysis engine
- Optimize tree traversal
- Concurrent analysis
- Memory optimization

**Key Learning**: AST-like analysis, pattern matching

### Days 121-125 — Roadmap Parser & Comparison
**Complexity: High** | **Priority: CRITICAL**

**Day 121-122**: Markdown Parser
- Parse structured roadmap documents
- Phase/milestone extraction
- Task identification
- Dependency detection

**Day 123**: Comparison Engine
- Match roadmap items to code
- Detect completed features
- Identify missing features
- Architectural drift detection

**Day 124**: Recommendation Engine
- Prioritize next tasks
- Dependency ordering
- Effort estimation
- Risk assessment

**Day 125**: Report Generation
- Markdown report output
- Summary statistics
- Progress visualization
- Actionable recommendations

**Key Learning**: Complex data processing and analysis

### Days 126-130 — UI & Polish
**Complexity: Medium-High** | **Priority: HIGH**

**Day 126-127**: Repo Planner UI
- Input form (URL, branch, roadmap)
- Progress indicator for analysis
- Results display with sections
- Interactive file tree

**Day 128**: Result Visualization
- Completion percentages
- Dependency graphs (optional)
- Comparison tables
- Drill-down details

**Day 129**: Export & Sharing
- Export as markdown
- Shareable permalinks
- PDF generation (optional)
- API endpoint for programmatic access

**Day 130**: Testing & Refinement
- End-to-end testing
- Performance testing with large repos
- Edge case handling
- User documentation

---

## Phase 9 — Advanced Features (Days 131-150) 🚀

### Days 131-140 — Multi-Repo Analysis
**Complexity: High** | **Priority: MEDIUM**

- Compare multiple repositories
- Monorepo support
- Cross-repo dependencies
- Team velocity analysis

### Days 141-150 — AI Enhancement (Optional)
**Complexity: Very High** | **Priority: LOW**

- LLM integration for code understanding
- Natural language task generation
- Code explanation generation
- Intelligent recommendations

---

## Go Mastery Checkpoints

### After Phase 2 (Day 30) ✓
- Can write concurrent code with goroutines and channels
- Understand when concurrency helps/hurts
- Know common concurrency patterns (worker pools, pipelines)
- Can use context for cancellation and timeouts

### After Phase 3 (Day 42) ✓
- Write idiomatic table-driven tests
- Profile and optimize code
- Set up CI/CD pipeline
- Deploy Go applications

### After Phase 4 (Day 55) ✓
- Design with interfaces
- Handle errors like a pro
- Use generics appropriately
- Write clean, idiomatic Go

### After Phase 5 (Day 70) ✓
- Work with SQL databases
- Implement repository pattern
- Handle transactions correctly
- Cache effectively

### After Phase 6 (Day 85) ✓
- Build production HTTP services
- Implement middleware
- Add observability
- Secure web applications

### After Phase 8 (Day 130) ✓
- **You are now a Go master**
- Can build complex, production-ready systems
- Understand Go idioms deeply
- Can read and contribute to any Go codebase

---

## Daily Learning Routine

### Every Day Should Include:
1. **Read Go code** (30 min): Standard library or popular projects
2. **Write code** (90 min): Implement the day's feature
3. **Test code** (30 min): Write tests, run benchmarks
4. **Review code** (15 min): Self-review, refactor, improve

### Weekly:
- Read one blog post about Go internals
- Watch one Go conference talk
- Contribute to an open-source Go project (after Day 60)

### Resources to Use Continuously:
- **Books**:
  - "Concurrency in Go" by Katherine Cox-Buday
  - "100 Go Mistakes" by Teiva Harsanyi
  - "The Go Programming Language" by Donovan & Kernighan
  
- **Websites**:
  - Go by Example (gobyexample.com)
  - Effective Go (official docs)
  - Go blog (blog.golang.org)
  
- **Practice**:
  - Exercism.io Go track
  - LeetCode in Go (for algorithms)
  - Read standard library source

---

## Summary of Changes from Original Plan

### Added (Critical for Mastery):
1. ✅ **Full concurrency phase** (Days 16-30) - THE most important addition
2. ✅ **Interface design phase** (Days 43-46)
3. ✅ **Complete database phase** (Days 56-70)
4. ✅ **Benchmarking and profiling** (Days 34-36)
5. ✅ **Advanced HTTP patterns** (Days 71-85)
6. ✅ **Error handling mastery** (Days 47-50)
7. ✅ **Generics** (Days 51-55)
8. ✅ **Production patterns** (observability, security, etc.)

### Removed (Low Learning Value):
1. ❌ WAV/MP3 audio processing
2. ❌ Color converter
3. ❌ QR code generator
4. ❌ Some redundant archive tools

### Restructured:
1. ⚡ Concurrency moved to Day 16 (was Day 61+)
2. ⚡ Testing moved earlier and expanded
3. ⚡ Repo analysis pushed to end (after mastery achieved)
4. ⚡ Better progression: foundations → concurrency → patterns → databases → advanced

---

## Expected Timeline

- **Part-time (2 hrs/day)**: ~7-8 months to Day 130
- **Full-time (6 hrs/day)**: ~3 months to Day 130
- **Aggressive (8 hrs/day)**: ~2 months to Day 130

**This is realistic** because each day is focused on one concept with hands-on practice.
