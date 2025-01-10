# BUILD STAGE
FROM golang:1.24-alpine AS builder
WORKDIR /src

COPY go.* .
RUN go mod download

COPY . .
RUN go build -o build/app-bin main.go

# RUN STAGE
FROM alpine
WORKDIR /src

COPY --from=builder /src/build/app-bin build/

ENV DB_GORM_PATH=:memory:
ENV PORT=3000

EXPOSE ${PORT}

CMD ["/src/build/app-bin"]
