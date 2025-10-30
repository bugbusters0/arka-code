.PHONY: run build install clean migrate

run:
	go run main.go

build:
	go build -o bin/app main.go

install:
	go mod download
	go mod tidy

clean:
	rm -rf bin/

migrate:
	mysql -u root -p < database.sql

dev:
	air 