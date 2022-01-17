build:
	mkdir -p bin
	go mod tidy
	go build -o bin/elastiSync main.go

clean:
	go clean
	rm -rf bin