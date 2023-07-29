default: build

clean:
	@rm -rf dist

build: clean
	go build -o dist/

run: build
	./dist/limacity-dns-update

docker: build
	docker build . -t ghcr.io/axelrindle/limacity-dns-update:latest
	docker compose up
