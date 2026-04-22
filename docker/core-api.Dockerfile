FROM golang:1.26.2 AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main ./cmd/api

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=build /app/main .
COPY --from=build /app/db ./db
COPY --from=build /app/geoip_ipv4.bin ./geoip_ipv4.bin
COPY --from=build /app/geoip_ipv6.bin ./geoip_ipv6.bin
EXPOSE 7171
CMD ["./main"]
