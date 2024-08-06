.PHONY: all
all: build run

.PHONY: build
build:
	go fmt
	go build -o server

.PHONY: run
run:
	./server

.PHONY: test
test: build
	./server &
	go test
	pkill server
