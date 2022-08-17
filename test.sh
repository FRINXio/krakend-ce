#!/bin/sh
export GOPRIVATE="github.com/FRINXio"
git config --global url."https://${GH_TOKEN}:x-oauth-basic@github.com/FRINXio".insteadOf "https://github.com/FRINXio"
make build
go test -v ./tests

