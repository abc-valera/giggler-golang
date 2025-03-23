FROM golang:1.24-alpine

# Variables that should be provided by the user
ENV OPENAPI_PORT=8081

# Env variables for Go modules and cache
ENV GOMODCACHE=/root/go/pkg/mod
ENV GOCACHE=/root/.cache/go-build

EXPOSE ${OPENAPI_PORT}

WORKDIR /workspace

COPY ./external/openapi/ .

RUN --mount=type=cache,target=${GOMODCACHE} \
    --mount=type=cache,target=${GOCACHE} \
    go build -o bin ./cmd/main.go

CMD ["./bin"]
