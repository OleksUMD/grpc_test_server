# Example of using a standalone gRPC server
This is a simple gRPC server with examples how to use it in unary/stream mode.

## Prerequisites
### All the following tools should be installed before trying to run the server
* [Go](https://go.dev/)
* [Docker](https://docs.docker.com/engine/install/ubuntu/) and [docker compose](https://docs.docker.com/compose/install/)
* [OpenSSL](https://github.com/openssl/openssl#build-and-install)
* [GNU make](https://github.com/mirror/make)

## Steps to run the project
> N.B. Tested only on Ubuntu 24.04 with Go 1.26.5
* Clone the repository:
`$ git clone https://github.com/OleksUMD/grpc_test_server.git`
* Go to the project's directory:
`$ cd grpc_test_server`
* Create .env file:
`$ touch .env` or `cp .env.example .env`
* Set up variables inside the .env file
* Run Makefile:
`$ make execute`

## Testing
There are two main ways how the server can be tested.
The first one is to use .proto file from this [repo](https://github.com/OleksUMD/proto-example) and import them into any appropriate desktop client(Postman, Insomnia, Yaak etc.)
The second one is to download test gRPC client from this [repo](https://github.com/OleksUMD) and follow the corresponding instructions there.
