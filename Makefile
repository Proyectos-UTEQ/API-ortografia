build:
	@go build -o bin/api-ortografia cmd/api/main.go

run: build
	@./bin/api-ortografia

seed: 
	@go run cmd/seed/main.go