.PHONY: build run test clean

build:
	go fmt ./...
	go build -o helloworld .

run:
	./helloworld

test:
	go test ./...

clean:
	rm -f helloworld
