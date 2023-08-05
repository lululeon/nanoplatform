# this docker image, which has go installed,  will build the go app:
FROM golang:1.21rc4-bookworm as builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o auth-service ./cmd/api

RUN chmod +x /app/auth-service

FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/auth-service /app

EXPOSE 7567

HEALTHCHECK --interval=20s CMD curl -f http://localhost:7567/health || exit 1

CMD ["/app/auth-service"]
