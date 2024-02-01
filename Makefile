.PHONY: generate-mocks
generate-mocks:
	bash mocks_gen.sh services
	bash mocks_gen.sh repository


.PHONY: test
test:
	go test ./services/...
	go test ./repository/...