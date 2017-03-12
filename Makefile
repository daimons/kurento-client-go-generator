.PHONYL: all

all: clean prep build format

clean:
	rm -rf kurento

prep:
	mkdir kurento

build:
	go run main.go
	go get golang.org/x/net/websocket

format:
	goimports -w ./kurento