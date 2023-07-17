#!/bin/sh
export GOPRIVATE="github.com/FRINXio/krakend-websocket"
git config --global url."https://${GH_TOKEN}:x-oauth-basic@github.com/FRINXio".insteadOf "https://github.com/FRINXio"
go mod tidy
make build
go test -v ./tests
