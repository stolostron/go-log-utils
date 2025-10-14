# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

go-log-utils is a Go library that provides utilities for setting up zap loggers in a consistent way across Open Cluster Management projects. The library is focused on integration with controller-runtime and klog, providing command-line configuration and standardized logging formats.

## Project Structure

The main code is in the `zaputil` package:
- `zaputil/zaputil.go` - Core implementation of logger configuration and builder functions
- `zaputil/doc.go` - Package documentation with usage examples

## Development Commands

### Formatting
```shell
make fmt
```
This runs multiple formatters in sequence:
- `gofmt -s` for standard Go formatting
- `gci` for import organization (groups by standard, default, then project prefix)
- `gofumpt` for stricter formatting rules

### Linting
```shell
make lint
```
Runs golangci-lint with the configuration in `.golangci.yml`. The project uses `enable-all: true` with specific linters disabled (see `.golangci.yml` for the full list).

### Security Scanning
```shell
make gosec-scan
```
Runs gosec security scanner and outputs results in sonarqube format to `gosec.json`.

### Testing
```shell
go test -v ./...
```
Note: Currently there are no test files in the repository.

### Pre-commit Checks
Before submitting a PR, run all three checks:
```shell
make fmt
make lint
make test
```

## Architecture

### FlagConfig
The `FlagConfig` struct is the primary interface for configuring zap loggers. It:
- Binds to command-line flags for log level (`--zap-log-level`) and encoder (`--zap-encoder`)
- Defaults to `info` level and `console` encoding
- Provides `GetConfig()` to generate a `zap.Config` based on the production preset with ISO-8601 timestamps

### Custom Level Encoding
The library implements custom level encoding (zaputil.go:53-60):
- Known levels (≥ -1) are formatted as text ("info", "error", etc.)
- Lower levels (< -2) are formatted as "lvl-N" to match klog verbosity conventions

### Integration Points

**controller-runtime**: Use `BuildForCtrl()` to create a logger, then set it with `ctrl.SetLogger(zapr.NewLogger(ctrlZap))`. Important: always call `WithName()` or `WithValues()` in your functions to get accurate caller information.

**klog**: Use `BuildForKlog()` with a FlagSet containing the klog flags. This function:
- Reads the `-v` flag value from klog's FlagSet
- Converts klog verbosity (positive integers) to zap levels (negative integers)
- Skips line endings in console mode since klog already adds them

**glog/klog interop**: Use `SyncWithGlogFlags()` to copy command-line flag values from glog to klog FlagSet, with special handling for the `log_backtrace_at` flag which is intentionally skipped.

## Code Standards

- All commits must be signed off with `git commit --signoff` (DCO requirement)
- Apache 2.0 license headers required on all Go files
- Line length limit: 120 characters
- Tab width: 4 spaces (for linting)
- Import organization: standard library, third-party, then project imports

## Dependencies

Core dependency: `go.uber.org/zap` v1.27.0

The library is designed to work with:
- `sigs.k8s.io/controller-runtime` (ctrl.SetLogger)
- `k8s.io/klog/v2` (for klog integration)
- `github.com/go-logr/zapr` (for logr adapter)

## Go Version

Requires Go 1.24.0 or later (as specified in go.mod).
