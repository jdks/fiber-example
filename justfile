set dotenv-load
set dotenv-filename := ".env.development"
alias b := build

build:
    go build -o server ./cmd/server/main.go
    go build -o seed ./cmd/seed/main.go

serve:
    ./server

seed:
    ./seed --users 3 --events 1000
