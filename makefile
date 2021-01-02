build:
	go build

install:
	go install

test:
	$(MAKE) -C ./cmd test
	$(MAKE) -C ./nanopool test

testreport:
	$(MAKE) -C ./cmd testreport
	$(MAKE) -C ./nanopool testreport