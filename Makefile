build:
	go mod tidy
	go build -o bin/

clean:
	go clean