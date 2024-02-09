FROM golang:1.22rc2-alpine3.19 AS build
WORKDIR /app
COPY . .

RUN go mod download

ENV CGO_ENABLED=0 GOOS=linux
RUN go build -o hasura-jwt ./cmd/api-server

FROM scratch
WORKDIR /app
COPY --from=build /app/hasura-jwt hasura-jwt
COPY letsencrypt.pem /etc/ssl/certs/

ENTRYPOINT ["./hasura-jwt"]