.PHONY: build run

export GOPATH:=$(shell pwd)

build:
	go fmt
	go install
