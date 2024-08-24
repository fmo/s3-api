S3_API_BINARY=s3ApiApp

s3_api_players:
	@echo "Building binary..."
	go build -o ${S3_API_BINARY} ./cmd/api/
	@echo "Done!"

s3_api_players_amd:
	@echo "Building binary..."
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ${S3_API_BINARY} ./cmd/api/
	@echo "Done!"
