# this docker image, which has go installed,  will build the go app:
FROM golang:1.20-alpine as builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o serviceBroker ./cmd/api

RUN chmod +x /app/serviceBroker

# this tiny, vanilla linux image, will just host and run the app
FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/serviceBroker /app

CMD ["/app/serviceBroker"]
