SINGLETON =
COMMANDS = rew


ifndef GOAMD64
	GOAMD64 = v2
endif

GOOS = $(shell uname -s | tr [A-Z] [a-z])
ifeq ($(GOOS), darwin)
	GOBIN = /usr/local/go/bin/go
	UPXBIN = /usr/local/bin/upx
else
	GOBIN = /usr/local/bin/go
	UPXBIN = /usr/bin/upx
endif

RELEASE = -s -w
GOARGS = GOOS=$(GOOS) GOARCH=amd64 GOAMD64=$(GOAMD64) CGO_ENABLED=1
GOBUILD = $(GOARGS) $(GOBIN) build -ldflags="$(RELEASE)"


.PHONY: all build clean upx upxx $(SINGLETON) $(COMMANDS)

all: clean build

$(SINGLETON):
	@echo "Compile $(SINGLETON) ..."
	$(GOBUILD) -o $(SINGLETON) *.go

$(COMMANDS):
	@echo "Compile $@ ..."
	$(GOBUILD) -o $@ ./cmd/$@

build: $(SINGLETON) $(COMMANDS)
	@echo "Build success."

clean:
	rm -f $(SINGLETON) $(COMMANDS)
	@echo "Remove old files."

upx: clean build
	$(UPXBIN) $(SINGLETON) $(COMMANDS)

upxx: clean build
	$(UPXBIN) --ultra-brute $(SINGLETON) $(COMMANDS)