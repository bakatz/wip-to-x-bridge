rm lambda-handler.zip 2>/dev/null
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bootstrap ./cmd/lambda/main.go
zip lambda-handler.zip bootstrap
rm bootstrap