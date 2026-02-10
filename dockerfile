FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod . go.sum ./
COPY getData.go .
COPY main.go .
COPY . .

RUN go get
RUN go build -o bin .

ENTRYPOINT [ "/app/bin" ]