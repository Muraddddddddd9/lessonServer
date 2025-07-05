FROM golang:1.24-alpine AS builder

WORKDIR /

COPY . .

RUN GOOS=linux

RUN go build  -o main .

EXPOSE 8080

CMD ["./main"]