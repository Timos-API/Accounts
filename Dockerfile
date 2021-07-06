FROM golang:latest as builder

RUN apt-get update && apt-get install ca-certificates tzdata -y && update-ca-certificates

ENV GO111MODULE="on"

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o app .

FROM scratch as prod

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app /app
COPY --from=builder /auth.html /auth.html

ENTRYPOINT ["/app"]
