FROM joukojo/golang:latest AS builder 

WORKDIR /app
COPY . .
RUN go build  -buildvcs=false -o mail-testserver ./cmd/mail-testserver

FROM debian:trixie-slim

WORKDIR /app

RUN apt-get update && apt-get install -y ca-certificates && \
useradd -m app && \
chown -R app:app /app && \
rm -rf /var/lib/apt/lists/*

USER app
COPY --from=builder /app/mail-testserver .
CMD ["./mail-testserver"]

EXPOSE 8025 8080

LABEL org.opencontainers.image.title="mailserver" \
      org.opencontainers.image.description="A lightweight SMTP test server with HTTP API for integration testing" \
      org.opencontainers.image.version="1.0.0" \
      org.opencontainers.image.authors="Jouko Johansson" \
      org.opencontainers.image.url="https://github.com/joukojo/go-mail-testserver" \
      org.opencontainers.image.source="https://github.com/joukojo/go-mail-testserver" \
      org.opencontainers.image.documentation="https://github.com/joukojo/go-mail-testserver/blob/main/README.md" \
      org.opencontainers.image.licenses="Apache-2.0" \
      org.opencontainers.image.vendor="joukojo" \
      org.opencontainers.image.base.name="debian:trixie-slim"