# ── Node Build stage ──────────────────────────────────────────────────────────
FROM node:22-alpine AS frontend-builder

WORKDIR /app

RUN corepack enable pnpm

COPY package.json pnpm-workspace.yaml pnpm-lock.yaml ./
COPY apps/web/package.json ./apps/web/
RUN pnpm install --frozen-lockfile

COPY . .
RUN pnpm build

# ── Go Build stage ────────────────────────────────────────────────────────────
FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /stories-go ./cmd/api

# ── Final stage ───────────────────────────────────────────────────────────────
FROM gcr.io/distroless/static:nonroot

COPY --from=builder /stories-go /stories-go
COPY --from=frontend-builder /app/apps/web/dist /dist

EXPOSE 8080

ENTRYPOINT ["/stories-go"]
