FROM golang:1.25.0-alpine3.22 AS builder

RUN apk --no-cache add curl

# Variables that should be provided by the user
ARG WORKDIR_PATH

ENV HOME=/root
ENV GOMODCACHE=$HOME/go/pkg/mod
ENV GOCACHE=$HOME/.cache/go-build
ENV RESTAPI_PORT=8080

WORKDIR ${WORKDIR_PATH}
EXPOSE ${RESTAPI_PORT}

RUN echo ${HOME}

COPY ./src ./src
COPY go.* .

RUN --mount=type=cache,target=${GOMODCACHE} \
	go mod download

RUN --mount=type=cache,target=${GOMODCACHE} \
	--mount=type=cache,target=${GOCACHE} \
	go build -o /usr/bin/giggler-golang-api ./src/cmd/restapi

FROM alpine:3.22 AS runner

RUN apk --no-cache add curl

# Variables that should be provided by the user
ARG WORKDIR_PATH

WORKDIR ${WORKDIR_PATH}

COPY --from=builder /usr/bin/giggler-golang-api .

CMD ["./giggler-golang-api"]
