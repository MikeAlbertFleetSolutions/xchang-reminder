package := $(shell basename `pwd`)

.PHONY: default get codetest test fmt lint vet cyclo secure build release

default: fmt codetest

get:
	go get -t -v ./...
	go get -u github.com/tcnksm/ghr
	go get -u golang.org/x/lint/golint
	go get -u github.com/fzipp/gocyclo
	go get -u github.com/securego/gosec/cmd/gosec/...

codetest: lint vet cyclo secure

build:
	mkdir -p target
	rm -f target/$(package)
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -v -o target/$(package)

release:
	mkdir -p target
	rm -f target/*
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -v -o target/$(package)
	zip -j -u target/$(package)_linux_amd64.zip target/$(package)
	rm -f target/$(package)
	ghr -c ${CIRCLE_SHA1} -u ${CIRCLE_PROJECT_USERNAME} -r ${CIRCLE_PROJECT_REPONAME} --replace "v${CIRCLE_BUILD_NUM}" target/

test:
	go test -v
	go test -cover

fmt:
	go fmt ./...

lint:
	@echo golint ./...
	@OUTPUT=`golint ./... 2>&1`; \
	if [ "$$OUTPUT" ]; then \
		echo "golint errors:"; \
		echo "$$OUTPUT"; \
		exit 1; \
	fi

vet:
	go vet -all .

cyclo:
	gocyclo -over 20 .

secure:
	gosec -quiet ./...
