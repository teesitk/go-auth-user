# Use Golang official image to build the Go application
FROM golang:1.24-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .
RUN go build -o main .

FROM alpine:latest  

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]