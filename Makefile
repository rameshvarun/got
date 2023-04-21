
.PHONY: install test clean

install: $(shell find . -type f -and -name '*.go')
	go install

test: install
	./runtests.sh

clean:
	find . -name "*~" -type f -delete
