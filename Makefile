# Makefile
.PHONY: build test coverage
BINARY_NAME = intelagent

build:
	make clean
	mkdir releases
	go build -o ./releases/dashclient ./cmd/dashclient/
	GOOS=linux go build -o ./releases/dashserver ./cmd/dashserver/
	GOOS=linux go build -o ./releases/proxy ./cmd/proxy/

update-readme:
	TOTAL_COVERAGE=$$(go tool cover -func=coverage.out | grep total | awk '{print $$3}'); \
	echo "TOTAL_COVERAGE: $${TOTAL_COVERAGE}"; \
	sed "s/### Total Coverage: .*/### Total Coverage: $${TOTAL_COVERAGE}/g" README.md > tREADME.md
	mv tREADME.md README.md

clean:
	@if [ -d ./releases/ ]; then rm -rf ./releases/; fi

test:
	go test -v ./...

coverage:
	go test -race -v -coverprofile=coverage.out ./... | tee coverage.report
	go tool cover -html=coverage.out -o coverage.html