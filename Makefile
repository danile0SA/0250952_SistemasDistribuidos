SHELL := powershell.exe
.SHELLFLAGS := -Command

# Variables
PROTOC_GEN_GO_PATH := $(shell go env GOPATH)\bin\protoc-gen-go.exe
PROTO_FILES := $(shell Get-ChildItem -Path api/v1 -Filter *.proto)
GO_TEST := go test -race ./...

compile:
	@if (Test-Path $(PROTOC_GEN_GO_PATH)) { \
		protoc api/v1/*.proto --go_out=. --go_opt=paths=source_relative --proto_path=.; \
	} else { \
		Write-Host "protoc-gen-go not found, installing..."; \
		go install google.golang.org/protobuf/cmd/protoc-gen-go@latest; \
		protoc api/v1/*.proto --go_out=. --go_opt=paths=source_relative --proto_path=.; \
	}

test:
	$(GO_TEST)

