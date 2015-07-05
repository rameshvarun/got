
.PHONY: default test clean

default: ${GOPATH}/bin/got

${GOPATH}/bin/got: $(shell find . -type f -and -name '*.go')
	go install

test:
	./runtests.sh

clean:
	rm ${GOPATH}/bin/got
	find . -name "*~" -type f -delete
