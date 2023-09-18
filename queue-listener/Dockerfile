FROM golang:1.21-alpine as builder

RUN mkdir -p /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o queueListener .

RUN chmod +x /app/queueListener


FROM alpine:latest

RUN mkdir -p /app

COPY --from=builder /app/queueListener /app

CMD ["/app/queueListener"]