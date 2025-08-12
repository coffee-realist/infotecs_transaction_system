FROM golang:1.24.5-alpine3.22 AS builder

WORKDIR /app
RUN apk add --no-cache git gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build \
    -ldflags="-linkmode external -extldflags -static" \
    -o transaction_app ./cmd

FROM scratch

WORKDIR /app

COPY --from=builder /app/transaction_app /app/transaction_app

ENTRYPOINT ["/app/transaction_app"]
