# syntax=docker/dockerfile:1
# Marketplace server image (API + storefront, catalog embedded); the committed export.json is the no-server alternative.

FROM node:26-alpine AS web
WORKDIR /web
COPY web/package.json web/package-lock.json ./
RUN npm ci
COPY web/ ./
RUN npm run build

FROM golang:1.27.3 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# The SPA is embedded from internal/web/dist (see internal/web/embed.go).
COPY --from=web /web/dist ./internal/web/dist
ARG VERSION=dev
ARG COMMIT=unknown
ARG BUILD_DATE=unknown
RUN CGO_ENABLED=0 go build \
    -ldflags "-s -w -X github.com/miabi-io/marketplace/internal/buildinfo.Version=${VERSION} -X github.com/miabi-io/marketplace/internal/buildinfo.CommitID=${COMMIT} -X github.com/miabi-io/marketplace/internal/buildinfo.BuildDate=${BUILD_DATE}" \
    -o /out/marketplace ./cmd/marketplace

FROM alpine:3.20
ARG VERSION=dev
ARG COMMIT=unknown
LABEL org.opencontainers.image.title="Miabi Marketplace" \
      org.opencontainers.image.version="${VERSION}" \
      org.opencontainers.image.revision="${COMMIT}" \
      org.opencontainers.image.source="https://github.com/miabi-io/marketplace" \
      org.opencontainers.image.licenses="Apache-2.0"
RUN apk add --no-cache ca-certificates
COPY --from=build /out/marketplace /usr/local/bin/marketplace
ENV MARKETPLACE_PORT=8088

ENV MARKETPLACE_BASE_URL=https://marketplace.miabi.io
EXPOSE 8088
ENTRYPOINT ["marketplace"]
CMD ["server"]
