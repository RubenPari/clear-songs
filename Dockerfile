ARG GO_VERSION=1
FROM golang:${GO_VERSION}-bookworm as builder

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY src/ ./src
WORKDIR /usr/src/app/src
RUN go build -v -o /run-app .

FROM debian:bookworm
COPY --from=builder /run-app /run-app
CMD ["/run-app"]
