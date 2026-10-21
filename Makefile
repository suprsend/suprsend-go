.PHONY: help lint

help:
	@echo "Usage:"
	@echo "  make lint                    golangci-lint run ./..."
	@echo "Install golangci-lint v2.15.0 (see AGENTS.md):"
	@echo "  curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b \"\$$(go env GOPATH)/bin\" v2.15.0"
	@echo "See AGENTS.md."

lint:
	golangci-lint run ./...
