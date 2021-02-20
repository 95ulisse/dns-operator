DOCS_IMG ?= dns-operator-docs

DOCS_RUN_OPTIONS ?= --rm -v $(CURDIR):/mkdocs -p 8000:8000

default: docs-build

docs-image:
	docker build -t $(DOCS_IMG) -f docs.Dockerfile ./

docs-serve: docs-image
	docker run $(DOCS_RUN_OPTIONS) $(DOCS_IMG) mkdocs serve

docs-build: docs-image
	docker run $(DOCS_RUN_OPTIONS) $(DOCS_IMG) sh -c 'mkdocs build && chown -R $(shell id -u):$(shell id -g) .'

.PHONY: default docs-image docs-serve docs-build