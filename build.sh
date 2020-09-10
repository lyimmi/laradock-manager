#!/bin/bash
wails build -ldflags="-s -w"
# upx build/laradock-manager
go run ./cmd/package
