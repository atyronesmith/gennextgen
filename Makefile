# SPDX-License-Identifier: Apache-2.0
# Copyright (c) 2024 Robin Jarry

version = $(shell git describe --long --abbrev=12 --tags --dirty 2>/dev/null || echo 0.1)
src = $(shell find * -type f -name '*.go') go.mod go.sum
go_ldflags :=
go_ldflags += -X main.version=$(version)
bin_dir = ./bin

.PHONY: all
all: env

env: $(src)  cmd/env/main.go
	go build -trimpath -ldflags='$(go_ldflags)' -o $(bin_dir)/$@ cmd/env/main.go

.PHONY: debug
debug: env

env.debug: $(src) cmd/env/main.go
	go build -ldflags='$(go_ldflags)' -gcflags=all="-N -l" -o $(bin_dir)/$@ cmd/env/main.go

.PHONY: format
format:
	gofmt -w .

.PHONY: lint
lint:
	@gofmt -d . | grep ^ \
	&& echo The above files need to be formatted. Please run make format. && exit 1 \
	|| echo All files formated.

.PHONY: run
run: env
	OVS_NODE_EXPORTER_CONFIG=etc/dev.conf ./$<