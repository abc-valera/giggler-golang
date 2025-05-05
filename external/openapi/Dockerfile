FROM golang:1.24-alpine
WORKDIR /workspace

COPY ./external/openapi/main.go .
COPY ./external/openapi/bundle.yaml .
COPY ./external/openapi/scalar_docs.html .

RUN go build -o bin ./main.go

ENV OPENAPI_PORT=3100
EXPOSE ${OPENAPI_PORT}

CMD ["./bin"]
