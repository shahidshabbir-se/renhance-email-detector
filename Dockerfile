# ----------------------------
# Stage 1: Base Go Build Image
# ----------------------------
FROM golang:1.24.3-alpine AS base

WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# ----------------------------
# Target: API Builder
# ----------------------------
FROM base AS api-builder
RUN go build -trimpath -ldflags="-s -w" -o api ./cmd/api/main.go

# ----------------------------
# Target: Worker Builder
# ----------------------------
FROM base AS worker-builder
RUN go build -trimpath -ldflags="-s -w" -o worker ./cmd/worker/main.go

# ----------------------------
# Stage 2: Minimal Runtime Image
# ----------------------------
FROM alpine:3.20 AS runtime

WORKDIR /app

RUN apk add --no-cache ca-certificates \
  && adduser -D -u 10001 appuser

COPY templates /app/templates
COPY web /app/web

COPY --from=api-builder /app/api .
COPY --from=worker-builder /app/worker .

USER appuser

# Default to running API, override with CMD in docker-compose or CLI
ENTRYPOINT ["./api"]
