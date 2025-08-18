set dotenv-load := true

# Run script in development mode
dev:
    go run ./src

# Build the source code into a single caddy-gitops binary
build-bin:
    go build -o caddy-gitops ./src

build-docker:
    docker build -t caddy-gitops:latest .

run-docker:
    docker run -p 2020:2020 --env-file .env caddy-gitops:latest