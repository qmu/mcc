NAME = mcc
VERSION = $(shell cat VERSION)

clean:
	rm -rf _build/ release/

build:
	glide install
	mkdir -p _build
	gox -osarch="linux/amd64 darwin/amd64 linux/386 darwin/386" -output="_build/{{.OS}}_{{.Arch}}_{{.Dir}}"

release:
	mkdir release
	go get github.com/progrium/gh-release/...
	cp _build/* release
	gh-release create qmu/$(NAME) $(VERSION)

.PHONY: build
