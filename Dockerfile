FROM golang:1.18.3-alpine

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /app

COPY . .

RUN go build -o myapp

EXPOSE 8080

ENTRYPOINT ["./myapp"]