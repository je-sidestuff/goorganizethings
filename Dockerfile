FROM golang:latest

WORKDIR /app

RUN go mod init github.com/preaje/goorganizethings

COPY . .

RUN go build -o main

ENTRYPOINT ["/app/main"]
