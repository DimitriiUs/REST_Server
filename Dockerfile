FROM golang:latest
WORKDIR /app
RUN go mod init github.com/DmitriiUs/REST_Server
COPY . .
WORKDIR /app/cmd/REST_Server/
RUN go build -o api
EXPOSE 3030

CMD ["./api"]
