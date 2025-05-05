# BUILD STAGE
FROM golang:1.24-alpine AS builder
WORKDIR /builder-workspace

COPY go.* .
RUN go mod download

COPY ./src ./src

RUN go build -o bin ./src/cmd/api


# RUN STAGE
FROM alpine
WORKDIR /runner-workspace

COPY --from=builder /builder-workspace/bin .

RUN apk --no-cache add curl

ENV API_PORT=3000
EXPOSE ${API_PORT}

CMD ["/runner-workspace/bin"]
