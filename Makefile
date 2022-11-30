init:
	sudo rm -rf ./images/name ./images/name.tar; tree .
test:
	sudo go run main.go build . -t name; tree .
build:
	sudo rm ./handlers/build.go; vim ./handlers/build.go

.PHONY: init test build