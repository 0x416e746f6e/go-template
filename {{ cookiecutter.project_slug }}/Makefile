VERSION := $(shell git describe --tags --always --dirty="-dev" --match "v*.*.*" || echo "development" )
VERSION := $(VERSION:v%=%)

.PHONY: build
build:
	@CGO_ENABLED=0 go build \
			-ldflags "-X main.version=${VERSION}" \
			-o ./bin/{{ cookiecutter.project_slug }} \
		github.com/{{ cookiecutter.github_org }}/{{ cookiecutter.project_slug }}/cmd

.PHONY: snapshot
snapshot:
	@goreleaser release --snapshot --clean

.PHONY: help
help:
	@go run github.com/{{ cookiecutter.github_org }}/{{ cookiecutter.project_slug }}/cmd serve --help
