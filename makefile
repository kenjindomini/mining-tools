build:
	go build

install:
	go install

test:
	$(MAKE) -C ./cmd test

testreport:
	$(MAKE) -C ./cmd testreport