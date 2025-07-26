FROM golang:1.25.0-alpine3.22 AS builder

RUN apk --no-cache add curl

# Variables that should be provided by the user
ARG WORKDIR_PATH
ARG BUILD_VERSION
# TODO: set default paths to ./local/cache
ENV GOMODCACHE=/root/go/pkg/mod
ENV GOCACHE=/root/.cache/go-build
ENV RESTAPI_PORT=8080

LABEL build_version=${BUILD_VERSION}

WORKDIR ${WORKDIR_PATH}
EXPOSE ${RESTAPI_PORT}

COPY ./src ./src
COPY go.* .

RUN --mount=type=cache,target=${GOMODCACHE} \
	go mod download

RUN --mount=type=cache,target=${GOMODCACHE} \
	--mount=type=cache,target=${GOCACHE} \
	go build \ 
	-ldflags="-X 'giggler-golang/src/shared/app.buildVersion=${BUILD_VERSION}'" \
	-o /usr/bin/giggler-golang-api \
	./src/cmd/restapi

FROM alpine:3.22 AS runner

RUN apk --no-cache add curl

# Variables that should be provided by the user
ARG WORKDIR_PATH
ARG BUILD_VERSION

WORKDIR ${WORKDIR_PATH}

LABEL build_version=${BUILD_VERSION}

COPY --from=builder /usr/bin/giggler-golang-api .

CMD ["./giggler-golang-api"]
