
FETCH_SOURCE ?= remote:git@github.com:writethesky/basic-proto.git
#FETCH_SOURCE ?= local:../basic-proto
FETCH_SOURCE_TYPE=$(shell echo ${FETCH_SOURCE}|awk -F ':' '{print $$1}')
FETCH_SOURCE_SITE=$(shell echo ${FETCH_SOURCE}|awk -F '${FETCH_SOURCE_TYPE}:' '{print $$2}')
FETCH_SOURCE_TYPE_REMOTE="remote"
FETCH_SOURCE_TYPE_LOCAL="local"
PROTO_DIRECTORY="proto"
PROTO_TARGET_DIRECTORY="pb"
BUF_VERSION:=1.0.0-rc6

install-tools:
	curl -sSL \
    	"https://github.com/bufbuild/buf/releases/download/v${BUF_VERSION}/buf-$(shell uname -s)-$(shell uname -m)" \
    	-o "$(shell go env GOPATH)/bin/buf" && \
    	chmod +x "$(shell go env GOPATH)/bin/buf"
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.6.0
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.6.0
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.27.1
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1.0
	go get github.com/vektra/mockery/v2/.../

clean:
	rm -rf ${PROTO_DIRECTORY} ${PROTO_TARGET_DIRECTORY}
fetch-proto: clean
	mkdir ${PROTO_DIRECTORY}
	@echo {\"type\": \"${FETCH_SOURCE_TYPE}\", \"site\": \"${FETCH_SOURCE_SITE}\"}
ifeq ($(FETCH_SOURCE_TYPE),$(shell echo $(FETCH_SOURCE_TYPE_REMOTE)))
	@echo "Ready to pull the remote code..."
	git clone ${FETCH_SOURCE_SITE} $(PROTO_DIRECTORY)
else ifeq ($(FETCH_SOURCE_TYPE),$(shell echo $(FETCH_SOURCE_TYPE_LOCAL)))
	@echo "Ready to copy local code..."
	cp -r ${FETCH_SOURCE_SITE}/* $(PROTO_DIRECTORY)/
endif
generate: fetch-proto
	buf generate ${PROTO_DIRECTORY}
	mockery --all --output ./mock

run:
	go run main.go

