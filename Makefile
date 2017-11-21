SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

all: nachos test bin
.PHONY: all test clean bin

bin:
	$(MAKE) -C bin
test: bin
	$(MAKE) -C test

nachos: $(SOURCES)
	go build nachos.go

clean:
	$(MAKE) -C test clean
	$(MAKE) -C bin clean
	rm nachos
