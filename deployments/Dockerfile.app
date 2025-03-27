FROM golang:1.24.1

WORKDIR /app

COPY ./go.mod ./go.sum ./   
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal
COPY .env .

RUN go build -o my-app ./cmd/main.go

CMD ["./my-app"]