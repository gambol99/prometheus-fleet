#
#   Author: Rohith
#   Date: 2015-08-21 16:22:01 +0100 (Fri, 21 Aug 2015)
#
#  vim:ts=2:sw=2:et
#

NAME=prometheus-fleet
AUTHOR=gambol99
HARDWARE=$(shell uname -m)
VERSION=$(shell awk '/Version =/ { print $$3 }' version.go | sed 's/"//g')

.PHONY: build docker docker-release release static test full-test clean

default: build

build:
	mkdir -p bin/
	go get
	go build -o bin/${NAME}

static:
	mkdir -p bin
	CGO_ENABLED=0 GOOS=linux go build -a -tags netgo -ldflags '-w' -o bin/${NAME}

docker: clean static
	sudo docker build -t ${AUTHOR}/${NAME} .

docker-release: docker 
	sudo docker tag -f ${AUTHOR}/${NAME} docker.io/${AUTHOR}/${NAME}:${VERSION}
	sudo docker push docker.io/${AUTHOR}/${NAME}:${VERSION}

full-test: build
	go get gopkg.in/yaml.v2
	go get github.com/stretchr/testify
	go test -v

test:
	go test -v

all: clean build docker

release: static
	mkdir -p release
	gzip -c bin/${NAME} > release/${NAME}_${VERSION}_linux_${HARDWARE}.gz
	rm -f release/${NAME}

clean:
	rm -rf ./bin
	rm -rf ./release
	go clean
