FROM golang:1.22 AS builder

WORKDIR /app

COPY . .

RUN go build -o url-shortener /app/cmd/url-shoertener

FROM ubuntu:22.04

RUN apt-get update && \
    apt-get install -y curl

WORKDIR /app

COPY --from=builder /app/url-shortener .
COPY config/local.yaml /app/

EXPOSE 8083

CMD ["./url-shortener"]
