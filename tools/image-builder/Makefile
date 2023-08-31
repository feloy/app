all: bin

.PHONY: bin
bin:
	go build -o image-builder main.go

.PHONY: install
install: bin
	mv image-builder ~/bin/