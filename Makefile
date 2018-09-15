.PHONY: test
test:
	go build -race
	go test -race -v ./...

.PHONY: release-deps
release-deps:
	GO111MODULE=off go get github.com/motemen/gobump
	GO111MODULE=off go get github.com/tcnksm/ghr

.PHONY: release
release: release-deps
	_tools/release.bash
