.PHONY: run build release pre-release

run:
	go run main.go

build:
	go build main.go -o ./bin/xend

release:
	goreleaser release

pre-release:
	goreleaser release --snapshot --clean
	@echo "Validating release files..."
	goreleaser check