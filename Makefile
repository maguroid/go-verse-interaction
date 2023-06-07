.PHONY: goerli sand sand/legacy
goerli:
	source .env && \
	go run ./goerli
sand/london:
	source .env && \
	go run ./sandverse/london
sand/legacy:
	source .env && \
	go run ./sandverse/legacy