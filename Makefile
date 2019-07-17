.PHONY: build push dep

GO          = go
GOGET       = $(GO) get -u
DEP         = $(GOPATH)/bin/dep
REPOSITORY  = crusttech/crust-bundle
BRANCH     ?= $(shell git rev-parse --abbrev-ref HEAD)
IMAGE_TAG  ?= $(if BRANCH=master,latest,$BRANCH)

build: dep
	docker build --no-cache --rm -t $(REPOSITORY):$(IMAGE_TAG) .

push:
	docker push $(REPOSITORY):$(IMAGE_TAG)

cdeps: $(DEP)
	$(DEP) ensure -update github.com/cortezaproject/corteza-server

$(DEP):
	$(GOGET) github.com/tools/godep
