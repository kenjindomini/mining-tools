build:
	go build

install:
	go install

test:
	$(MAKE) -C ./miningtools test

testreport:
	$(MAKE) -C ./miningtools testreport