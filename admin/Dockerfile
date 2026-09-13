FROM node:20-alpine AS web-builder
WORKDIR /app/web
COPY admin/web/package*.json ./
RUN npm install
COPY admin/web/ ./
RUN npm run build

FROM golang:1.22-alpine AS go-builder
WORKDIR /app/server
COPY admin/server/go.mod ./
RUN go mod download
COPY admin/server/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o fitguide-admin .

FROM python:3.11-alpine
RUN apk add --no-cache bash ca-certificates
WORKDIR /app

COPY --from=go-builder /app/server/fitguide-admin /app/fitguide-admin
COPY --from=web-builder /app/web/dist /app/admin/web/dist

ENV PORT=8080
ENV FITGUIDE_BASE_DIR=/app
ENV FITGUIDE_STATIC_DIR=/app/admin/web/dist

EXPOSE 8080

CMD ["/app/fitguide-admin"]
