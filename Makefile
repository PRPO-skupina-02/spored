docs:
	swag fmt
	swag init --parseDependency -g api/api.go -o api/docs

test:
	godotenv go test ./...