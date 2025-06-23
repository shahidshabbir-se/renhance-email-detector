# Variables
PROXY_PORT := 8080
TEMPL_DEV_PORT := 7331

.PHONY: dev backend frontend clean build

## Start full-stack dev environment (Templ + Air + Proxy)
dev:
	@echo "🚀 Starting Templ + Air with proxy to :$(PROXY_PORT)"
	templ generate --watch --proxy="http://localhost:$(PROXY_PORT)" --cmd="air"

## Run only backend with air (no templ)
backend:
	@echo "🌀 Starting backend with Air"
	air

## Generate templ files (one-time)
frontend:
	@echo "🎨 Generating templ components..."
	templ generate

## Clean generated templ files
clean:
	@echo "🧹 Cleaning generated .go files from .templ..."
	find . -type f -name '*.templ.go' -delete

## Compile Go backend (without running)
build:
	@echo "🔨 Building Go binary..."
	go build -o bin/app .
