FROM joukojo/golang:latest AS builder 

WORKDIR /app
COPY . .
RUN go build -buildvcs=false -o mail-testserver ./cmd/mail-testserver

FROM debian:trixie-slim

WORKDIR /app

RUN apt-get update && apt-get install -y ca-certificates && \
useradd -m app && \
chown -R app:app /app && \
rm -rf /var/lib/apt/lists/*

USER app
COPY --from=builder /app/mail-testserver .
CMD ["./mail-testserver"]