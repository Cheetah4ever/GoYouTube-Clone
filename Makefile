.PHONY: test coverage

COVERAGE_PROFILE := coverage.out
COVERAGE_REPORT := coverage.html

test:
	go test ./...

coverage:
	go test -coverprofile=$(COVERAGE_PROFILE) ./...
	go tool cover -func=$(COVERAGE_PROFILE)
	go tool cover -html=$(COVERAGE_PROFILE) -o $(COVERAGE_REPORT)
	@echo "HTML coverage report: $(COVERAGE_REPORT)"
