FROM golang:1.21-alpine as builder

RUN mkdir -p /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o brokerService ./api

RUN chmod +x /app/brokerService


FROM alpine:latest

RUN mkdir -p /app

COPY --from=builder /app/brokerService /app

CMD ["/app/brokerService"]
