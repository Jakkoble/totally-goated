FROM golang:1.25-alpine AS build
RUN apk add --no-cache gcc musl-dev
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" wasm/wasm_exec.js
RUN GOOS=js GOARCH=wasm go build -o wasm/game.wasm .
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o /server_bin ./server

FROM alpine:latest
RUN apk add --no-cache sqlite-libs
WORKDIR /app
COPY --from=build /server_bin /app/server
COPY --from=build /src/wasm /app/wasm
EXPOSE 8080
CMD ["/app/server"]
