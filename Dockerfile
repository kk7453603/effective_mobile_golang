FROM golang:1.22.4

WORKDIR /app

COPY . /app

RUN go build -o ./bin/main ./cmd/main.go

EXPOSE ${Docker_Port}

CMD ["/app/bin/main"]

