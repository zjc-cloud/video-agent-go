
FROM golang:1.22-alpine

WORKDIR /app

COPY . .

RUN go mod init video-agent-go && \
    go mod tidy && \
    go build -o videoagent cmd/main.go

EXPOSE 8080

CMD ["./videoagent"]
