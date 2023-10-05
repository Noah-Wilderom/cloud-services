FROM golang:1.21-alpine as builder

RUN mkdir -p /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o projectService .

RUN chmod +x /app/projectService


FROM alpine:latest

RUN mkdir -p /app

COPY --from=builder /app/projectService /app
COPY --from=builder /app/templates /templates

CMD ["/app/projectService"]