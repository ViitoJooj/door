FROM golang:1.26.2 AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main ./cmd/api
FROM debian:bookworm-slim
WORKDIR /app
COPY --from=build /app/main .
EXPOSE 8080
CMD ["./main"]