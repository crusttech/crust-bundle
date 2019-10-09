.PHONY: build push dep

GO          = go
GOGET       = $(GO) get -u
DEP         = $(GOPATH)/bin/dep
REPOSITORY  = crusttech/crust-bundle
BRANCH     ?= $(shell git rev-parse --abbrev-ref HEAD)
IMAGE_TAG  ?= $(if BRANCH=master,latest,$BRANCH)

build:
	docker build --no-cache --rm -t $(REPOSITORY):latest .

cdeps:
	$(GO) get github.com/cortezaproject/corteza-server
