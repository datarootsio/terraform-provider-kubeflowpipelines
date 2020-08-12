WEBSITE_REPO=github.com/hashicorp/terraform-website
PKG_NAME=kubeflowpipelines
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)

default: build test

build:
	go install

test:
	go test -v ./...  

testacc:
	TF_ACC=1 KUBEFLOWPIPELINES_HOST=http://localhost:8080 go test -v ./... -timeout 120m -covermode=count -coverprofile=coverage.out

fmt:
	gofmt -w $(GOFMT_FILES)

website:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

website-test:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider-test PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)