default: build

clean:
	@rm -rf dist

OUTPUT ?= dist/
build: clean
	for arch in amd64 arm64 ; do \
		echo "Building for linux/$$arch" … ; \
		GOOS=linux GOARCH="$$arch" go build -ldflags="-s -w" -o "$(OUTPUT)/$$arch/" . ; \
	done

test:
	go test -v -cover -coverprofile=coverage.out ./...

run: build
	./dist/limacity-dns-update

docker: build
	docker build . -t ghcr.io/axelrindle/limacity-dns-update:latest
	docker compose up
