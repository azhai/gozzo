COMMANDS = rew

ifndef GOAMD64
	GOAMD64 = v2
endif
GOOS = $(shell uname -s | tr [A-Z] [a-z])
ifeq ($(GOOS), darwin)
	GOBIN = /usr/local/go/bin/go
	UPXBIN = /usr/local/bin/upx
else
	GOBIN = /usr/bin/go
	UPXBIN = /usr/bin/upx
endif
RELEASE = -s -w
GOARGS = GOOS=$(GOOS) GOARCH=amd64 GOAMD64=$(GOAMD64) CGO_ENABLED=1
GOBUILD = $(GOARGS) $(GOBIN) build -ldflags="$(RELEASE)"


.PHONY: all build clean upx $(COMMANDS)

all: clean build

$(COMMANDS):
	@echo "Compile $@ ..."
	$(GOBUILD) -o ./bin/$@ ./cmd/$@

build: $(COMMANDS)
	@echo "Build success."

clean:
	#$(GOBIN) clean
	rm -f ./bin/*
	@echo "Clean all."

upx: clean build
	$(UPXBIN) ./bin/*
