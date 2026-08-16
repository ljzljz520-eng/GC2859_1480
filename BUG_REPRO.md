# BUG_REPRO

The following failures were observed while validating the initial project state.
Each section records what failed, how to reproduce it, and the complete command output.
They are preserved intentionally; only failing build gates are omitted from the generated Dockerfile.

## Failure 1: Go test (.)

- Observed problem: `Go test (.)` failed in the initial project state.
- Working directory: `.`
- Command: `cd /app && GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off go test -count=1 ./...`
- Exit status: `1`

```text
?   	wirelab.local/cabling/cmd/wirelab	[no test files]
--- FAIL: TestSearchEndpointReturnsJSONArrayForNoMatch (0.00s)
    handler_test.go:41: expected [] in response, got null
--- FAIL: TestSearchNoMatchReturnsEmptyArrayValue (0.00s)
    service_test.go:26: expected a non-nil empty result
FAIL
FAIL	wirelab.local/cabling/internal/articles	0.005s
FAIL
```

## Architecture reproduction

### linux/amd64
- Go toolchain version: exit `0`
- Node.js version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/wirelab): exit `0`
- Frontend build (web): exit `0`
### linux/arm64
- Go toolchain version: exit `0`
- Node.js version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/wirelab): exit `0`
- Frontend build (web): exit `0`
