# Tools Hub — Long-Term Roadmap

Tools Hub is a long-horizon project designed to evolve from a suite of practical developer utilities into a **repo-aware planning agent** capable of analyzing real codebases and proposing technically sound next steps.

This roadmap documents that progression.

---

## Project North Star

**Ultimate Goal (Months → Years)**  
Build a repo-aware planning agent that can:

- Ingest a GitHub repository
- Understand its structure and architecture
- Detect what is implemented vs missing
- Compare reality against an intended roadmap
- Propose the next concrete engineering steps

Tools Hub begins as a collection of utilities and matures into an engineering analysis and planning system.

---

## Phase 1 — Core Architecture & Tooling Foundations

### Day 1–4 — Bootstrap (Completed)
- HTTP server + middleware
- Template cache (base + partials + pages)
- File upload UI (file converter)
- Panic recovery + logging

### Day 5–7 — First Real Tool: Slugify (Completed)
- `/tools/slugify` GET/POST
- Internal pure function (`textutil`)
- Pattern established:
  - one route → one handler
  - internal method switch
  - pure logic under `internal/tools`

### Day 8 — Routing Normalization (Completed)
- Removed inline route closures
- Unified file converter into single handler
- `Routes()` is declarative wiring only

### Day 9 — Template Data Foundation (Completed)
- Introduced canonical `templateData` envelope
- All templates render with non-nil data
- Tool-specific data lives under `ToolData`
- Shared fields (`Error`, `Flash`, `PageTitle`) standardized

---

## Phase 2 — Structured Text & Encoding Tools

### Day 10–12 — JSON Formatter Tool
- `/tools/json` GET/POST
- Internal `jsonutil` package
- Pretty-print via `json.Indent`
- Minify support
- Graceful malformed JSON errors
- Input/output preserved across renders

### Day 13–14 — Base64 Tool
- Encode / decode modes
- Internal pure encoding logic
- Inline output + optional download
- Clear validation and error messaging

### Day 15–18 — Hashing Tools
- `/tools/hash`
- SHA-256 and MD5 (with security disclaimer)
- Text input + optional file upload
- Downloadable `.hash` output
- Internal `hashutil` package

### Day 19–21 — Timestamp Converter
- Epoch ↔ human-readable
- Timezone selection
- Strict parsing and validation
- Internal `timeutil` package

---

## Phase 3 — Quality, Testing, and Deployment

### Day 22–23 — Continuous Integration
- GitHub Actions pipeline:
  - `go fmt`
  - `go vet`
  - `go test ./...`
  - Build `cmd/web`
- Module caching
- CI required for merge

### Day 24–27 — Tests
- Unit tests for:
  - slugify
  - fileconvert
  - json formatter
  - base64
  - hashing
- HTTP handler tests via `httptest`
- Multipart upload tests for fileconvert

### Day 28–30 — Docker + First Deployment
- Multi-stage Dockerfile
- Minimal runtime image
- First deployment (Fly.io / Render / Railway)
- Optional auto-deploy from `main`

---

## Phase 4 — File & Archive Intelligence

### Day 31–35 — Archive Builder (ZIP)
- Upload multiple files
- Generate downloadable `.zip`
- Safe path normalization
- Size and file-count limits
- Internal `archiveutil` package

### Day 36–40 — Archive Extractor (ZIP)
- Upload `.zip`
- Inspect or repackage contents
- Protect against zip-slip attacks
- Limits on extracted size
- Tests using malicious fixtures

### Day 41–45 — TAR / GZIP Tools
- `.tar.gz` creator
- `.tar.gz` extractor
- Reuse archive patterns
- Consistent safety constraints

---

## Phase 5 — Binary & Audio File Understanding

### Day 46–48 — WAV Metadata Reader
- Upload `.wav`
- Parse header:
  - sample rate
  - channels
  - bit depth
  - duration estimate
- No decoding, metadata only

### Day 49–52 — WAV → MP3 Conversion
- Controlled external tool integration (e.g. `ffmpeg`)
- Strict file validation
- Temporary sandboxing
- Resource limits + timeouts

### Day 53–55 — MP3 Metadata Editor
- Edit title / artist
- Return updated MP3
- Validate tags safely

---

## Phase 6 — Text, Unicode, and Byte Mastery

### Day 56 — ASCII & Unicode Explorer
- Display bytes, runes, code points
- Hex + decimal views

### Day 57 — UTF-8 Analyzer
- Show raw byte sequences
- Mark multibyte boundaries

### Day 58 — Rune Inspector
- Rune count vs byte count
- Indexing differences visualized

### Day 59 — Unicode-Aware Palindrome Tool
- Rune-safe two-pointer logic
- Debug visualization
- Web-based explanation of algorithm

### Day 60 — Byte-Level Debugger
- Upload file
- Hex dump view (like `hexdump -C`)
- Pagination and size limits

---

## Phase 7 — Repo-Aware Planning Agent (Core Vision)

### Day 61–75 — Repo Planner MVP

#### Tool: `/tools/repoplan`

**Input**
- GitHub repository URL
- Optional branch
- Optional roadmap text

**Core Capabilities**
- Fetch repo metadata via GitHub API
- Build normalized file tree
- Detect:
  - routes
  - handlers
  - templates
  - internal tool packages
- Identify missing fundamentals:
  - tests
  - CI
  - Docker
  - architectural inconsistencies

**Output**
- Repo inventory summary
- Completed vs missing milestones
- Suggested next 5–10 tasks
- Grouped by category

No agent loops yet — deterministic, explainable output only.

---

## Phase 8 — Planner v2: Context & History

### Day 76–90
- Accept prior roadmap text
- Compare planned vs detected state
- Detect architectural drift
- Highlight completed milestones
- Suggest rebased next steps

---

## Long-Term Evolution (Open-Ended)

Potential future directions:
- Plan diffs over time
- Repo “health score”
- GitHub issue generation
- PR-sized task chunking
- Multi-language repo analysis
- Architectural contract enforcement

---

## Guiding Principles

- **One route → one handler**
- **Pure logic lives in `internal/tools`**
- **Templates receive a stable data contract**
- **No code execution during repo analysis**
- **Explainable, deterministic reasoning before AI heuristics**

---

## Summary

Tools Hub is not just a collection of utilities.  
It is an evolving system designed to understand, analyze, and plan software projects — including itself.

This roadmap is intentionally ambitious, but each phase produces real, reviewable engineering artifacts.

