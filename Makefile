build:
	@go build -o ./bin/teletimer
run: build
	@./bin/teletimer