#!/bin/bash
GOOS=js GOARCH=wasm go build -o wasm/game.wasm .