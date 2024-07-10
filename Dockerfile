FROM golang:1.22.2 as builder

WORKDIR /app

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64 

COPY . .

RUN go mod tidy

RUN go build -o /service cmd/main.go

EXPOSE 7540

CMD ["/service"]