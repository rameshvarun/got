
.PHONY: default test clean

default: ${GOPATH}/bin/got

${GOPATH}/bin/got: $(shell find . -type f -and -name '*.go')
	go install

test:
	go test
	echo "No tests yet"

clean:
	rm ${GOPATH}/bin/got