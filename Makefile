.PHONY: test
test:
	go test -v -race  -count=1 -cover -covermode=atomic -coverprofile=coverage.out ./...

.PHONY: viewcoverage
viewcoverage:
	go tool cover -html=coverage.out

.PHONY: bench
bench:
	go test $(verbose) -bench=. ./...

