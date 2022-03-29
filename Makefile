BINS = goo

.PHONY: all $(BINS) install clean test format

all: $(BINS)

$(BINS):
	go build ./cmd/$@

install:
	go install ./cmd/goo

clean:
	rm -f $(BINS)

test:
	go test ./test/* -v -race

format:
	gofmt -s -w .
