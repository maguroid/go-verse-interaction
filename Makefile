.PHONY: goerli goerli/london goerli/legacy sand/london sand/legacy
goerli/london:
	source .env && \
	go run ./goerli/london
goerli/legacy:
	source .env && \
	go run ./goerli/legacy
sand/london:
	source .env && \
	go run ./sandverse/london
sand/legacy:
	source .env && \
	go run ./sandverse/legacy