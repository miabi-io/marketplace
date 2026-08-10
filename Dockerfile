# syntax=docker/dockerfile:1
#
# Marketplace server image: serves the catalog API + storefront (the catalog is
# embedded, so the binary is stateless). This is what runs the hosted registry at
# marketplace.miabi.io, and what you deploy to self-host a fork or an internal
# mirror; the committed export.json remains the no-server alternative.

FROM node:26-alpine AS web
WORKDIR /web
COPY web/package.json web/package-lock.json ./
RUN npm ci
COPY web/ ./
RUN npm run build

FROM golang:1.26.3 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# The SPA is embedded from internal/web/dist (see internal/web/embed.go).
COPY --from=web /web/dist ./internal/web/dist
ARG VERSION=dev
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o /out/marketplace ./cmd/marketplace

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
COPY --from=build /out/marketplace /usr/local/bin/marketplace
ENV MARKETPLACE_PORT=8088
EXPOSE 8088
ENTRYPOINT ["marketplace"]
CMD ["server"]
