# Beaker v3 Executable Checklist

Target version: Go 1.27.1.

Branch target: `v3dev`.

## Scope

This checklist tracks the security, quality, performance, and CLI behavior work for the v3 development line. MySQL schema changes are explicitly out of scope. Redis cache can be cleared and rebuilt.

## Completed in this change

- [x] Lock Go version in `go.mod` with `go 1.27` and `toolchain go1.27.1`.
- [x] Upgrade vulnerable dependencies to current available versions: Gin `v1.12.0`, `golang.org/x/net v0.59.0`, `golang.org/x/crypto v0.57.0`, `golang.org/x/sys v0.48.0`, `golang.org/x/text v0.42.0`, and protobuf `v1.36.12`.
- [x] Add stable server-to-CLI error codes: `BEAKER-400`, `BEAKER-401`, `BEAKER-404`, `BEAKER-500-DB`, `BEAKER-500-CACHE`, `BEAKER-500`, `BEAKER-400-DECODE`, `BEAKER-CLI-NETWORK`, `BEAKER-CLI-HTTP`, and `BEAKER-CLI-RESPONSE`.
- [x] Make CLI failures exit with a non-zero status.
- [x] Make CLI errors directly print `Error Code`, `Description`, `Detail`, and `Request ID` in the terminal.
- [x] Make CLI errors explain the likely problem and next check, including MySQL, Redis, memory, network, login, decode, and response-format failures.
- [x] Add request ids to JSON responses so CLI output can be matched with server logs.
- [x] Replace the old global admin login key with per-login session tokens. New CLI requests must send `X-Beaker-Session`.
- [x] Add HTTP server timeouts, max header size, request id middleware, panic recovery, and request body limits.
- [x] Replace Redis one-connection-per-operation with a Redis connection pool.
- [x] Change cache clearing from `FLUSHALL` to deleting only keys under `REDIS_PREFIX`.
- [x] Cache Jet template sets instead of rebuilding them for every render.
- [x] Add image upload validation for filename, suffix, detected MIME type, max 5 MB size, and atomic write.
- [x] Render article/page Markdown with raw HTML/style skipped and safe-link filtering enabled.
- [x] Fix tweet empty-list pagination so page zero does not underflow database offsets.
- [x] Validate admin JSON decode errors consistently instead of silently writing zero values.
- [x] Avoid crashes on invalid 3DES ciphertext and use cryptographic random bytes for session keys.

## Must Run Before Release

- [x] `gofmt -w cmd internal`
- [x] `git diff --check`
- [x] `GOCACHE=/tmp/beaker-gocache go test ./...`
- [x] `go mod tidy` after final dependency upgrades.
- [x] `GOCACHE=/tmp/beaker-gocache go test -race ./...`.
- [x] `go vet ./...`.
- [ ] `govulncheck ./...` was not run because `govulncheck` is not installed on this machine.
- [ ] Re-check GitHub Dependabot after pushing `v3dev`.

## CLI Error Output Acceptance Example

```text
Error Code: BEAKER-500-DB
Description: database operation failed. Check MySQL address, credentials, network, connection limits, and whether MySQL is out of disk or memory.
Detail: dial tcp 127.0.0.1:3306: connect: connection refused
Request ID: f1c2...
```

## Remaining Follow-Up Work

- [ ] Replace the legacy MD5 password check with a password hash migration plan. This likely needs a MySQL-compatible strategy, so it was not done in this change.
- [ ] Replace 3DES request encryption with TLS-only session auth or an AEAD protocol in a later breaking release.
- [ ] Add integration tests for login, expired session, bad ciphertext, database error, Redis error, and malformed server response.
- [ ] Add CI for `go test`, `go vet`, `govulncheck`, and Dependabot.

## Data Compatibility

- MySQL schema: unchanged.
- Redis: safe to clear. The new clear behavior deletes only keys matching `REDIS_PREFIX + "*"`.
- Server TOML: unchanged in this change.
- Local CLI config: new `SessionToken` is written under each website entry after login.
