FROM node:22.16.0-alpine AS node-dependencies
WORKDIR /workspace
COPY package.json package-lock.json ./
COPY apps/builder/package.json ./apps/builder/package.json
RUN npm ci
COPY apps/builder ./apps/builder

FROM golang:1.23-alpine AS build
WORKDIR /src
COPY apps/api/go.mod apps/api/go.sum ./
RUN go mod download
COPY apps/api .
RUN go build -o /out/api ./cmd/api

FROM node:22.16.0-alpine
WORKDIR /app
RUN apk add --no-cache ca-certificates git
COPY --from=build /out/api /app/api
COPY --from=node-dependencies /workspace/node_modules /app/node_modules
COPY --from=node-dependencies /workspace/apps/builder /app/builder
COPY --from=build /src/migrations /app/migrations
RUN mkdir -p /var/lib/cms-builder/builds
ENV BUILDER_DIRECTORY=/app/builder
ENV BUILD_OUTPUT_ROOT=/var/lib/cms-builder/builds
ENV WRANGLER_COMMAND=/app/node_modules/.bin/wrangler
EXPOSE 8080
CMD ["/app/api"]
