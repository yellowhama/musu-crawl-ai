# syntax=docker/dockerfile:1.7
# Build stage — pure-Go static binary (no CGO, modernc/sqlite is pure-Go).
FROM golang:1.26-bookworm AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /out/musu-crawl .

# Runtime stage — distroless (static), includes ca-certs + tzdata.
FROM gcr.io/distroless/static-debian12:latest
COPY --from=build /out/musu-crawl /usr/local/bin/musu-crawl
WORKDIR /app
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/musu-crawl"]
