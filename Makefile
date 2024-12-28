SINGLETON =
COMMANDS  = rew


ifndef GOAMD64
	GOAMD64 = v2
endif

GOBIN    = go
UPXBIN   = upx
#GOOS    = $(shell uname -s | tr [A-Z] [a-z])
#GOARCH  = $(shell uname -m | tr [A-Z] [a-z])
GOARGS   = GOAMD64=$(GOAMD64) CGO_ENABLED=1
RELEASE  = "-s -w"
GOBUILD  = $(GOARGS) $(GOBIN) build -ldflags=$(RELEASE)
BINFILES = $(SINGLETON) $(COMMANDS)


.PHONY: all build clean upx upxx $(BINFILES)

all: clean build

$(SINGLETON):
	@echo "Compile $@ ..."
	$(GOBUILD) -o ./bin/$@ *.go

$(COMMANDS):
	@echo "Compile $@ ..."
	$(GOBUILD) -o ./bin/$@ ./cmd/$@

build: $(BINFILES)
	@echo "Build success."

clean:
	rm -f $(BINFILES:%=./bin/%)
	@echo "Remove old files."

upx: clean build
	$(UPXBIN) $(BINFILES:%=./bin/%)

upxx: clean build
	$(UPXBIN) --ultra-brute $(BINFILES:%=./bin/%)
