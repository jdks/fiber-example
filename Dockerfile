FROM golang:1.21-alpine

WORKDIR /app
COPY . .
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server/main.go

EXPOSE 3000
CMD ["./server"]
