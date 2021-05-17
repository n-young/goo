.PHONY: all clean

all: build

build:
	go build ./cmd/goo

install:
	go install ./cmd/goo

clean:
	rm goo
