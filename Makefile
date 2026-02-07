# Makefile
.PHONY: help generate clean-generate

# Variables
SWAGGER_FILE=api/swagger.json
GENERATED_DIR=internal/api/generated
GENERATED_FILE=$(GENERATED_DIR)/api.gen.go

generate:
	go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest \
			-generate types,chi-server \
			-package api \
			-o $(GENERATED_FILE) \
			$(SWAGGER_FILE)