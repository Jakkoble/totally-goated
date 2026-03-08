FROM golang:1.25-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" wasm/wasm_exec.js
RUN GOOS=js GOARCH=wasm go build -o wasm/game.wasm .

FROM nginx:alpine
COPY --from=build /src/wasm/ /usr/share/nginx/html/
EXPOSE 80
