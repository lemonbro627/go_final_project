FROM golang:1.22.2 as builder

WORKDIR /app

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64 TODO_PORT=7540

COPY . .

RUN go mod download

RUN go build -o /service cmd/main.go

EXPOSE $TODO_PORT

CMD ["/service"]