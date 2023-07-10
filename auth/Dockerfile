# this docker image, which has go installed,  will build the go app:
FROM golang:1.20-alpine as builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o auth-service ./cmd/api

RUN chmod +x /app/auth-service

FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/auth-service /app

CMD ["/app/auth-service"]
