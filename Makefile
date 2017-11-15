NAME := mcc
VERSION := v0.9.5
CONFIG_SCHEMA_VERSION := v1.1.0
SRCS      := $(shell find . -name '*.go' -type f)
LDFLAGS   := -ldflags "-X github.com/qmu/mcc/controller.Version=$(VERSION) -X github.com/qmu/mcc/controller.ConfigSchemaVersion=$(CONFIG_SCHEMA_VERSION)"
GH_UPLOAD := github-release upload --user qmu --repo $(NAME) --tag $(VERSION)

version:
	go run $(LDFLAGS) *.go -v

run:
	go run $(LDFLAGS) *.go -c _example/example.yml

erd:
	go-erd -path ./widget |dot -Tsvg > ./_build/widget_erd.svg
	go-erd -path ./model |dot -Tsvg > ./_build/model_erd.svg
	go-erd -path ./model/vector |dot -Tsvg > ./_build/vector_erd.svg

fmt:
	gofmt -s -w ./

clean:
	rm -rf _build/ release/

build:
	glide install
	mkdir -p _build
	gox $(LDFLAGS) -osarch="linux/amd64 darwin/amd64 linux/386 darwin/386" -output="_build/{{.OS}}_{{.Arch}}_{{.Dir}}"

# test > textfile > cat > rm... this is necessary because screen would be flush during tests
test:
	go test github.com/qmu/mcc/... -cover > _build/test.txt && cat _build/test.txt
	@rm _build/test.txt

# same reason above
bench:
	go test github.com/qmu/mcc/... -bench . -benchmem > _build/bench.txt && cat _build/bench.txt
	@rm _build/bench.txt

lines:
	@echo "=== implements =========================="
	@wc -l $(shell find . -name "*.go" | grep -v /vendor/ | grep -v _test.go)
	@echo "--- without line break & comment line"
	@find . -name "*.go" | grep -v /vendor/ | grep -v _test.go | xargs grep -h "^\s*[^\/\/]" | wc -l
	@echo "=== test code ==========================="
	@wc -l $(shell find . -name "*_test.go" | grep -v /vendor/)
	@echo "--- without line break & comment line"
	@find . -name "*_test.go" | grep -v /vendor/ | xargs grep -h "^\s*[^\/\/]" | wc -l


release:
	mkdir release
	go get github.com/aktau/github-release/...
	cp _build/* release
	github-release release \
		--user qmu \
		--repo $(NAME) \
		--tag $(VERSION) \
		--name $(VERSION)

	cd release/ \
		&& $(GH_UPLOAD) --name darwin_386_mcc --file darwin_386_mcc \
		&& $(GH_UPLOAD) --name darwin_amd64_mcc --file darwin_amd64_mcc \
		&& $(GH_UPLOAD) --name linux_386_mcc --file linux_386_mcc \
		&& $(GH_UPLOAD) --name linux_amd64_mcc --file linux_amd64_mcc

	git fetch --tags

.PHONY: build
