default: build

clean:
	@rm -rf dist

build: clean
	go build -o dist/

run: build
	./dist/limacity-dns-update
