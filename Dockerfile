FROM golang:1.24.1 AS build

ENV GO111MODULE=auto \
    CGO_ENABLED=0 \
    GOOS=linux

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o app

FROM debian:bookworm-slim AS runtime

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=build /app/app .

CMD ["./app"]

