FROM golang:1.24.1

WORKDIR /app

RUN apt-get update && apt-get install -y wait-for-it && rm -rf /var/lib/apt/lists/*

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal
COPY .env .

RUN go build -o my-app ./cmd/main.go

CMD ["sh", "-c", "wait-for-it tarantool:3301 --timeout=30 --strict -- ./my-app"]
