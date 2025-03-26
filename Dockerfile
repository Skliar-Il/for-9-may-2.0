FROM golang:1.24-alpine

RUN mkdir /app

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build ./cmd/app/main.go

CMD ./main

