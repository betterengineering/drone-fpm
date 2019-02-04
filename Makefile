test:
	go test -cover ./...

build: clean
	go build -o out/generate-entrypoint github.com/lodge93/drone-fpm/cmd/generate-entrypoint

clean:
	rm -rf out/

generate-entrypoint:
	./out/generate-entrypoint
