#!/usr/bin/env sh
# Build GoKeyMux with version info injected via -ldflags (Linux/macOS).
# Usage: ./build.sh
set -eu

Version=$(git describe --tags --always --dirty 2>/dev/null || printf 'dev')
Commit=$(git rev-parse --short HEAD 2>/dev/null || printf 'unknown')
BuildTime=$(date -u '+%Y-%m-%dT%H:%M:%SZ')

ServerLdFlags="-s -w -X GoKeyMux/version.Version=$Version -X GoKeyMux/version.Commit=$Commit -X GoKeyMux/version.BuildTime=$BuildTime"

go build -trimpath -ldflags "$ServerLdFlags" -o GoKeyMux .
go build -trimpath -ldflags "-s -w" -o keycycle ./cmd/keycycle
go build -trimpath -ldflags "-s -w" -o loadtest ./cmd/loadtest

printf 'Built GoKeyMux %s (commit=%s built=%s)\n' "$Version" "$Commit" "$BuildTime"
