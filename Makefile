build:
	@go build -o ./bin/teletimer
	
run: build
	@./bin/teletimer

test:
	@go test -v