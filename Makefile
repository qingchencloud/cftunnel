VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -s -w -X github.com/qingchencloud/cftunnel/cmd.Version=$(VERSION)

build:
	go build -ldflags "$(LDFLAGS)" -o cftunnel .

clean:
	rm -f cftunnel

.PHONY: build clean
