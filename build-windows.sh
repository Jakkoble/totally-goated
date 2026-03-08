#!/bin/bash
GOOS=windows GOARCH=amd64 go build -ldflags="-H windowsgui" -o build/totally-goated.exe