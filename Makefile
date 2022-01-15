build:
	mkdir -p bin
	go mod tidy
	go build -o bin/elasticSync main.go

clean:
	go clean
	rm -rf bin