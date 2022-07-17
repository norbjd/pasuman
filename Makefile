GIT_TAG = $(shell git describe --tags 2>/dev/null)
GIT_COMMIT = $(shell git rev-parse --short HEAD)
GIT_TREE_STATE = $(shell (git status --porcelain | grep -q .) && echo "-dirty" || echo "")

EXECUTABLE_OUTPUT = pasuman

.PHONY: build
build:
	@CGO_ENABLED=0 go build -ldflags " \
		-X github.com/norbjd/pasuman/cmd.GitVersion=$(GIT_TAG) \
		-X github.com/norbjd/pasuman/cmd.GitCommit=$(GIT_COMMIT)$(GIT_TREE_STATE) \
		-X github.com/norbjd/pasuman/cmd.BuildDate=`date --utc +%Y-%m-%dT%H:%M:%SZ`" \
		-o $(EXECUTABLE_OUTPUT)

.PHONY: build_all_platforms
build_all_platforms:
	$(MAKE) build GOOS=linux GOARCH=amd64 EXECUTABLE_OUTPUT=pasuman-linux-amd64
	$(MAKE) build GOOS=darwin GOARCH=amd64 EXECUTABLE_OUTPUT=pasuman-darwin-amd64
	$(MAKE) build GOOS=darwin GOARCH=arm64 EXECUTABLE_OUTPUT=pasuman-darwin-arm64

.PHONY: test
test:
	@go test -coverprofile=coverage.out ./...

.PHONY: coverage
coverage: test
	@go tool cover -func=coverage.out

.PHONY: lint
lint:
	@docker run --rm -v `pwd`:/app -w /app golangci/golangci-lint:v1.46.2 \
		golangci-lint run -v

.PHONY: gofumpt
gofumpt:
	@docker run --rm -v `pwd`:/app -w /app golang:1.18-alpine \
		sh -c 'go install mvdan.cc/gofumpt@latest && gofumpt -l -w .'

.PHONY: license_fix
license_fix:
	@docker run --rm -v `pwd`:/app -w /app golang:1.18-alpine \
		sh -c 'apk add git && \
			go install github.com/apache/skywalking-eyes/cmd/license-eye@v0.3.0 && \
			license-eye -c skywalking-eyes-config.yaml header fix'

.PHONY: license_check
license_check:
	@docker run --rm -v `pwd`:/app -w /app golang:1.18-alpine \
		sh -c 'apk add git && \
			go install github.com/apache/skywalking-eyes/cmd/license-eye@v0.3.0 && \
			license-eye -c skywalking-eyes-config.yaml header check && \
			license-eye -c skywalking-eyes-config.yaml dep check'

# count lines of Go code (excluding tests)
# remove all import lines
# remove all blank lines
# remove all comment lines
# remove all lines before first ")" (end of imports)
.PHONY: sloc
sloc:
	@find . -name '*.go' ! -name '*_test.go' | \
		grep -v 'pasumantest' | \
		xargs -I file awk ' \
			BEGIN		{ start=0; } \
			/^import/ 	{ next; } \
			/^\s*\/\// 	{ next; } \
			/^$$/      	{ next; } \
			/^)$$/		{ start=1; next; } \
						{ if (start==1) print; }' file | \
		wc -l
