FROM node:22-bookworm AS frontend-builder
WORKDIR /app
COPY client/package.json client/package-lock.json ./
RUN npm install && npm install @rolldown/binding-linux-x64-gnu lightningcss-linux-x64-gnu
COPY client/ .
RUN npm run build

FROM golang:1.22-bookworm AS backend-builder
WORKDIR /app
COPY service/go.mod service/go.sum ./
RUN go mod download
COPY service/ .
RUN CGO_ENABLED=1 go build -o /app/server .

FROM debian:bookworm-slim
RUN mkdir -p /data
WORKDIR /app
COPY --from=backend-builder /app/server .
COPY --from=frontend-builder /app/dist /client/dist
EXPOSE 3000
ENV HTTP_HOST=0.0.0.0
CMD ["./server"]
