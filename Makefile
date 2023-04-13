OUT := out/tests
SOURCE_FILES := $(shell find . -name '*.go')

$(OUT): $(SOURCE_FILES)
	go build -o $@

.DEFAULT_GOAL := build
.PHONY: build
build: $(OUT)

.PHONY: help
help: build
	@echo "\n*** to pass arguments to the tests, call $(OUT) directly ***\n\n"
	@$(OUT) --help

.PHONY: run
run: build
	$(OUT)