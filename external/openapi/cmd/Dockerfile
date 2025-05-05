FROM golang:1.24-alpine
WORKDIR /workspace

COPY ./external/openapi/ .

RUN go build -o bin ./cmd/main.go

ENV OPENAPI_PORT=3100
EXPOSE ${OPENAPI_PORT}

CMD ["./bin"]
