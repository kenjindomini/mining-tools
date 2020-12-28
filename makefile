build:
	go build

install:
	go install

test:
	cd miningtools; go test -coverprofile coverage.out

testreport:
	go tool cover -HTML coverage.out