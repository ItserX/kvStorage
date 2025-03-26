FROM golang:1.24.1

WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o my-app ./cmd/main.go

CMD ["./my-app"]