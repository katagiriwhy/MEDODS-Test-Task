FROM golang:1.24.2
LABEL authors="katagiri"

WORKDIR /Medods

COPY . .

RUN go mod download

RUN go mod tidy

RUN go build -o api ./cmd/service/main.go

EXPOSE 8080

CMD ["./api"]