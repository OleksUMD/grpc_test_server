BIN_FILENAME=test-grpc-server

.PHONY: clean build execute docker-up docker-down gen-ssl-certs

gen-ssl-certs:
	# example
	# CA_PASS=1234 SERVER_PASS=1234 make gen-ssl-certs
	# or set them in the .env file
	./openssl.sh

docker-up:
	docker compose --env-file .env up -d --wait postgres

docker-down:
	docker compose --env-file .env down

clean:
	rm -rf ./bin

build: clean
	go build -o ./bin/$(BIN_FILENAME) ./cmd

execute: clean gen-ssl-certs build docker-down docker-up
	# example
	# CA_PASS=1234 SERVER_PASS=1234 make execute
	# or set them in the .env file
	./bin/$(BIN_FILENAME)
