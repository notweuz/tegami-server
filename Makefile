.PHONY: build run clean

build:
	go build -o ./bin/tegami-server ./cmd/tegami-server

run: build
	./bin/tegami-server

clean:
	rm -rf ./bin