init:
	sudo rm -rf ./images/name ./images/name.tar ./containers/*; tree .
test:
	sudo go run main.go build . -t name; tree .; sudo go run main.go run --name test name
vim-build:
	sudo rm ./handlers/build.go; vim ./handlers/build.go
vim-run:
	sudo rm ./handlers/run.go; vim ./handlers/run.go

.PHONY: init test vim-build vim-run