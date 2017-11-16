NAME := mcc
VERSION := v0.9.6
CONFIG_SCHEMA_VERSION := v1.1.0
SRCS      := $(shell find . -name '*.go' -type f)
LDFLAGS   := -ldflags "-X github.com/qmu/mcc/controller.Version=$(VERSION) -X github.com/qmu/mcc/controller.ConfigSchemaVersion=$(CONFIG_SCHEMA_VERSION)"
GH_UPLOAD := github-release upload --user qmu --repo $(NAME) --tag $(VERSION)

.PHONY: version
version:
	go run $(LDFLAGS) *.go -v

.PHONY: run
run:
	go run $(LDFLAGS) *.go -c _example/example.yml

.PHONY: erd
erd:
	go-erd -path ./widget |dot -Tsvg > ./_build/widget_erd.svg
	go-erd -path ./model |dot -Tsvg > ./_build/model_erd.svg
	go-erd -path ./model/vector |dot -Tsvg > ./_build/vector_erd.svg

.PHONY: fmt
fmt:
	gofmt -s -w ./

.PHONY: clean
clean:
	rm -rf _build/ release/

.PHONY: build
build:
	glide install
	mkdir -p _build
	gox $(LDFLAGS) -osarch="linux/amd64 darwin/amd64 linux/386 darwin/386" -output="_build/${NAME}_${VERSION}_{{.OS}}_{{.Arch}}/{{.Dir}}"

# test > textfile > cat > rm... this is necessary because screen would be flush during tests
.PHONY: test
test:
	go test github.com/qmu/mcc/... -cover > _build/test.txt && cat _build/test.txt
	@rm _build/test.txt

# same reason above
.PHONY: bench
bench:
	go test github.com/qmu/mcc/... -bench . -benchmem > _build/bench.txt && cat _build/bench.txt
	@rm _build/bench.txt

.PHONY: lines
lines:
	@echo "=== implements =========================="
	@wc -l $(shell find . -name "*.go" | grep -v /vendor/ | grep -v _test.go)
	@echo "--- without line break & comment line"
	@find . -name "*.go" | grep -v /vendor/ | grep -v _test.go | xargs grep -h "^\s*[^\/\/]" | wc -l
	@echo "=== test code ==========================="
	@wc -l $(shell find . -name "*_test.go" | grep -v /vendor/)
	@echo "--- without line break & comment line"
	@find . -name "*_test.go" | grep -v /vendor/ | xargs grep -h "^\s*[^\/\/]" | wc -l

.PHONY: release
release:
	rm release
	@make clean
	@make build
	mkdir release
	go get github.com/aktau/github-release/...
	cp -R _build/* release
	github-release release \
		--user qmu \
		--repo $(NAME) \
		--tag $(VERSION) \
		--name $(VERSION)

	cd release/ \
		&& $(GH_UPLOAD) --name darwin_386_mcc --file ${NAME}_${VERSION}_darwin_386/mcc \
		&& $(GH_UPLOAD) --name darwin_amd64_mcc --file ${NAME}_${VERSION}_darwin_amd64/mcc \
		&& $(GH_UPLOAD) --name linux_386_mcc --file ${NAME}_${VERSION}_linux_386/mcc \
		&& $(GH_UPLOAD) --name linux_amd64_mcc --file ${NAME}_${VERSION}_linux_amd64/mcc \
		&& tar czvf ${NAME}_${VERSION}_darwin_amd64.tar.gz ${NAME}_${VERSION}_darwin_amd64/ \
		&& $(GH_UPLOAD) --name ${NAME}_${VERSION}_darwin_amd64.tar.gz --file ${NAME}_${VERSION}_darwin_amd64.tar.gz \
		&& echo openssl dgst -sha256 ${NAME}_${VERSION}_darwin_amd64.tar.gz

	git fetch --tags
