FROM golang:1.21-alpine as builder

RUN mkdir -p /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o loggerService ./api

RUN chmod +x /app/loggerService


FROM alpine:latest

RUN mkdir -p /app

COPY --from=builder /app/loggerService /app

CMD ["/app/loggerService"]