# ── Build stage ───────────────────────────────────────────────────────────────
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /stories-go ./cmd/api

# ── Final stage ───────────────────────────────────────────────────────────────
FROM gcr.io/distroless/static:nonroot

COPY --from=builder /stories-go /stories-go

EXPOSE 8080

ENTRYPOINT ["/stories-go"]
