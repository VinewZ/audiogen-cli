sql:
	sqlc generate

goose-up:
	goose sqlite3 ./books.db -dir ./internal/sql/schema up

goose-down:
	goose sqlite3 ./books.db -dir ./internal/sql/schema down

build:
	go build -o ./bin/audiogen-cli ./
