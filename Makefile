.PHONY: all
all: build run

.PHONY: build
build:
	go fmt
	go build -o server

.PHONY: clean
clean:
	rm server

.PHONY: run
run:
	./server

.PHONY: kill
kill:
	if pgrep server; then pkill server; fi

.PHONY: test
test: kill build
	./server &
	go test
	pkill server
