FROM golang:1.25.4-alpine3.21 AS builder
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY . .

RUN go mod download

RUN go build -o hasura-jwt ./cmd/api-server

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
USER 1001:1001

WORKDIR /app
COPY --from=builder /app/hasura-jwt hasura-jwt

EXPOSE 3000
ENTRYPOINT ["./hasura-jwt"]
