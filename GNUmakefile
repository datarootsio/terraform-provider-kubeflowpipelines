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