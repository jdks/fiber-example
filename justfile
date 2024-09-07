set dotenv-load
set dotenv-filename := ".env.development"
alias b := build

build:
    go build -o server ./cmd/server/main.go

serve:
    ./server