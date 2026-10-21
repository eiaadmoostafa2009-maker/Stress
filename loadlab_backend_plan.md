# LoadLab Backend Plan (Cobra CLI)

## 1. Project goal

LoadLab is a local command-line tool for testing HTTP endpoints that I own or am explicitly authorized to test. It runs bounded concurrent HTTP requests, reports latency/throughput/errors, protects the local machine, and stores test summaries in SQLite.

The first version is a CLI application, not a web application. Cobra replaces the Vue interface. Since there is no HTTP API in this version, Gin and JWT are not required yet.

## 2. Layer architecture

```text
cmd/loadlab/
  root.go             # Cobra root command and shared setup
  run.go              # `loadlab run`
  history.go          # `loadlab history`
  version.go          # `loadlab version`

internal/
  config/              # Defaults and configuration validation
  models/              # TestConfig, TestResult, TestRecord
  service/              # Use cases and business rules
  loadtest/             # Worker pool, HTTP client, collector, metrics
  resource/             # CPU/RAM monitoring and safety decisions
  repository/           # Interfaces and persistence implementations
  database/             # SQLite connection and migrations
```

### Responsibilities

- **Cobra / CLI:** Parse flags, print help, display results, and call services. Do not put worker logic here.
- **Service:** Validate the requested test, enforce business rules, coordinate the resource guard and load-test engine, and save the final summary.
- **Load-test engine:** Use a bounded worker pool, context-aware HTTP requests, bounded channels, and a collector. Avoid creating an unbounded goroutine per request.
- **Resource guard:** Sample system and process resource usage, issue warnings, and tell the manager to reduce scheduling, pause, or cancel according to policy.
- **Repository:** Persist and retrieve test history. The engine should not write every individual request to SQLite.
- **Database:** Open SQLite and apply schema migrations.
- **Models:** Hold the data types shared between layers.

## 3. Example command

```bash
loadlab run \
  --url http://localhost:9000/health \
  --concurrency 25 \
  --duration 30s
```

Flow:

1. Cobra parses the flags.
2. The CLI constructs a `TestConfig`.
3. The service validates the URL, duration, and concurrency.
4. The resource guard checks whether the machine is in an acceptable state.
5. The manager starts a bounded worker pool with a cancellable context.
6. Workers send compact results through a bounded channel.
7. A collector computes totals, error counts, throughput, and latency percentiles.
8. The repository saves one aggregate test summary.
9. The CLI prints the report.

## 4. Suggested core types

```go
type TestConfig struct {
    URL         string
    Concurrency int
    Duration    time.Duration
}

type RequestResult struct {
    StatusCode int
    Latency    time.Duration
    Err        error
}

type TestSummary struct {
    TotalRequests  int64
    FailedRequests int64
    RequestsPerSec float64
    P50            time.Duration
    P95            time.Duration
    P99            time.Duration
}
```

Add fields as requirements become clear. Validate configuration at the service boundary, not only in Cobra flags.

## 5. Resource-safety rules

- Start with a conservative default concurrency, such as 25.
- Define a hard concurrency ceiling, but never treat that ceiling as a safe recommendation for every computer.
- Allow only one active test initially.
- Require a finite test duration and per-request timeout.
- Bound worker count, job/result queues, and response-body reading.
- Do not retain entire response bodies or persist every request by default.
- Monitor available system memory and CPU; distinguish system-wide metrics from LoadLab process metrics.
- Use warning thresholds with hysteresis to avoid rapid state changes.
- On critical resource pressure, stop scheduling new jobs or cancel the test, drain results, and save the final summary.
- Validate targets and redirects to reduce accidental access to unrelated or sensitive local/network services.
- Test against a local endpoint first and increase concurrency gradually.

Monitoring is a safety signal, not a guarantee that the OS cannot run out of memory.

## 6. Middleware and JWT decision

Middleware is a function that wraps an HTTP request handler to apply cross-cutting behavior, such as authentication, logging, request IDs, timeouts, or recovery.

A JWT authentication middleware usually:
1. Extracts a token from an expected header or cookie.
2. Verifies its signature and permitted algorithm.
3. Checks claims such as expiration (`exp`) and issuer/audience when configured.
4. Places the verified identity/claims into the request context.
5. Rejects the request with an appropriate status if verification fails; otherwise, it calls the next handler.

A valid JWT proves that the token was issued by a trusted signer and has valid claims. It does **not** automatically prove the user is allowed to access a specific resource. Authorization checks must still verify permissions and resource ownership.

For the local Cobra-only MVP, there is no inbound HTTP request pipeline, so Gin middleware and JWT authentication are unnecessary. If a Gin API is added later, middleware belongs at the HTTP boundary and should call the same service layer as the CLI. Do not add JWT just because the load generator sends HTTP requests to a target: target authentication is a separate feature, configured only when testing an authorized application that requires it.

## 7. Build order

1. Create the Cobra CLI and `version` command.
2. Add configuration types and validation.
3. Add a local test HTTP server for safe development.
4. Implement the bounded worker pool and cancellation.
5. Add result collection and metrics.
6. Add resource monitoring and safety decisions.
7. Add SQLite models, migrations, and repository.
8. Add `history` and persisted summaries.
9. Add repeatable tests, race tests, and resource-limit tests.
10. Consider an optional Gin API only if a real use case needs it.

## 8. Definition of done for the first milestone

- `loadlab --help` works.
- `loadlab version` works.
- Invalid URLs, durations, and concurrency values are rejected.
- A short test against a local endpoint completes and prints a summary.
- Cancellation terminates workers without leaking goroutines.
- Memory use is bounded by design, not merely observed after a crash.
- A completed test summary can be saved and read from SQLite.
