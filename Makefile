.PHONY: run build release pre-release

run:
	go run main.go

build:
	go build main.go -o ./bin/xend

release:

pre-release:
	goreleaser release --snapshot --clean