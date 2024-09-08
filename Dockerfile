FROM golang:alpine

WORKDIR /app

COPY . .

ENV GOOS=linux
ENV GOARCH=amd64

RUN go mod download

RUN go build -o main main.go

EXPOSE 8080

CMD ["./main"]
