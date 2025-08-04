set dotenv-load := true

# Run script in development mode
dev:
    go run ./src

# Build the source code into a single caddy-gitops binary
build-bin:
    go build -o caddy-gitops ./src

