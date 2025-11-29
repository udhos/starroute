#!/bin/bash

go install golang.org/x/vuln/cmd/govulncheck@latest
go install golang.org/x/tools/cmd/deadcode@latest
go install github.com/mgechev/revive@latest

gofmt -s -w .

revive ./...

echo gocyclo begin
gocyclo -over 15 .
echo gocyclo end

echo tidy
go mod tidy

echo govulncheck
govulncheck ./...

echo deadcode
deadcode ./cmd/*

echo test
go test -race ./...

echo install
go install ./...
