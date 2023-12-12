build:
	@go build -o bin/art-digital cmd/api/main.go

run: build
	@./bin/art-digital

seed: 
	@go run cmd/seed/main.go