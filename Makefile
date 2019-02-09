test:
	go test -cover ./...

build: clean
	go build -o out/generate-entrypoint github.com/lodge93/drone-fpm/cmd/generate-entrypoint

clean:
	rm -rf out/

generate-entrypoint:
	./out/generate-entrypoint

docker: build generate-entrypoint
	docker build -f build/Dockerfile -t quay.io/lodge93/drone-fpm:$(shell git log -n 1 --pretty=format:"%H") .

end-to-end-test: docker
	docker run --rm -v $(PWD):/workdir -w /workdir -e PLUGIN_DEB_SYSTEMD=/workdir/test/generate-entrypoint.service -e PLUGIN_NAME=generate-entrypoint -e PLUGIN_VERSION=snapshot-$(shell git log -n 1 --pretty=format:"%H") -e PLUGIN_INPUT_TYPE=dir -e PLUGIN_OUTPUT_TYPE=deb -e PLUGIN_PACKAGE=/workdir/out/generate-entrypoint-snapshot-$(shell git log -n 1 --pretty=format:"%H").deb -e PLUGIN_COMMAND_ARGUMENTS=/workdir/out/generate-entrypoint=/usr/local/bin/ quay.io/lodge93/drone-fpm:$(shell git log -n 1 --pretty=format:"%H")