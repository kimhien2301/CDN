GOSOURCE = $(shell ls **/*.go)
BINDIR = bin/
SOURCEDIR = src/
TARGET = mirage
GO = go
GIT = git

all: format $(BINDIR)$(TARGET)

format:
	./tools/go_source/recursive-gofmt.sh

golint:
	./tools/go_source/golint.sh

go-vet:
	./tools/go_source/go-vet.sh

$(BINDIR)plugin.so:
	if [ -e src/mirageplugin ]; then $(GO) build -ldflags="-s -w" -buildmode=plugin -o $@ mirageplugin; fi

plugin: $(BINDIR)plugin.so

$(BINDIR)$(TARGET): plugin
	$(GO) build -ldflags="-s -w" -o $@ mirage

run: $(SOURCEDIR)$(TARGET).go
	$(GO) run $<

clean:
	rm -f $(BINDIR)$(TARGET) $(BINDIR)plugin.so
	rm -f bin/*
	rm -f mirage.dat

git-push: format
	$(GIT) push

