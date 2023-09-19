FROM golang:1.21-alpine as builder

RUN mkdir -p /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o queueWorker ./api

RUN chmod +x /app/queueWorker


FROM alpine:latest

RUN mkdir -p /app

COPY --from=builder /app/queueWorker /app

CMD ["/app/queueWorker"]