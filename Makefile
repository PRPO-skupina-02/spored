docs:
	swag fmt
	swag init --parseDependency -g api/api.go -o api/docs

test:
	godotenv go test ./...

install-cli-tools:
	go install github.com/joho/godotenv/cmd/godotenv@latest
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest