FROM alpine:latest as certificates
RUN apk update && apk add --no-cache ca-certificates && update-ca-certificates

FROM golang:1.17 as builder
WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/cyber_power main.go

FROM scratch
COPY --from=certificates /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/cyber_power /app/cyber_power
RUN mkdir -p /etc/cyberpower && \
    touch /etc/cyberpower/config.yaml

WORKDIR /app
ENTRYPOINT [ "/app/cyber_power" ]
