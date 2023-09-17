default: build

clean:
	@rm -rf dist

OUTPUT ?= dist/
build: clean
	go build -ldflags="-s -w" -o $(OUTPUT) .

test:
	go test -v -cover -coverprofile=coverage.out ./...

run: build
	./dist/limacity-dns-update

docker: build
	docker build . -t ghcr.io/axelrindle/limacity-dns-update:latest
	docker compose up
