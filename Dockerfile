FROM golang:alpine

WORKDIR /app

COPY src/ .

ENV GOOS=linux
ENV GOARCH=amd64

RUN go mod download

RUN go build -o src main.go

EXPOSE 8080

CMD ["./src"]
