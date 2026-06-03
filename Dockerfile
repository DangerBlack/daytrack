FROM node:22-alpine AS frontend-builder
WORKDIR /app
COPY client/package.json client/package-lock.json ./
RUN npm ci
COPY client/ .
RUN npm run build

FROM golang:1.22-alpine AS backend-builder
WORKDIR /app
COPY service/go.mod service/go.sum ./
RUN go mod download
COPY service/ .
RUN CGO_ENABLED=1 go build -o /app/server .

FROM alpine:3.19
RUN apk --no-cache add ca-certificates tzdata
RUN mkdir -p /data /archive
WORKDIR /app
COPY --from=backend-builder /app/server .
COPY --from=frontend-builder /app/dist /app/client/dist
EXPOSE 3000
CMD ["./server"]
