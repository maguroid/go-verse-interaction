.PHONY: goerli goerli/london goerli/legacy sand/london sand/legacy
goerli/london:
	source .env && \
	go run ./cmd/goerli/london
goerli/legacy:
	source .env && \
	go run ./cmd/goerli/legacy
sand/london:
	source .env && \
	go run ./cmd/sandverse/london
sand/legacy:
	source .env && \
	go run ./cmd/sandverse/legacy
build/counter:
	mkdir -p ./gen
	solc --evm-version berlin --bin --abi --optimize --overwrite -o ./gen ./Counter.sol
	abigen --bin=./gen/Counter.bin --abi=./gen/Counter.abi --pkg=counter --out=./lib/counter/counter.go
	rm -rf ./gen