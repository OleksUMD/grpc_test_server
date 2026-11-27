#!/bin/bash
# Output files
# ca.key: Certificate Authority private key file (in real-life : this shouldn't be shared)
# ca.crt: Certificate Authority trust certificate (in real-life : share this with gRPC client)
# server.key: Server private key, password protected (in real-life : this shouldn't be shared)
# server.csr: Server certificate signing request (in real-life : share this with CA owner)
# server.crt: Server certificate to be "installed" on server (in real-life : this shouldn't be shared)
# server.pem: server.key, converted into a format required by gRPC (in real-life : this shouldn't be shared)
#
# Private files (not to be shared): ca.key, server.key, server.pem, server.crt
# Public files (to be shared): ca.crt (share to gRPC client), server.csr (share to CA)

SERVER_CN=localhost
MY_SUBJECT="/CN=${SERVER_CN}"

export $(grep -v '^#' .env | xargs -d '\n')

if [ "$CA_PASS" == "" ] || [ "$SERVER_PASS" == "" ]; then
	echo "The required environment variables(CA_PASS & SERVER_PASS) are not set";
	exit 1;
fi

# create if does not exist otherwise cleanup the corresponding directory
mkdir -p ./certs
rm ./certs/*

# 1: Generate Certificate Authority + Trust Certificate (ca.crt)
openssl genrsa -passout env:CA_PASS -aes256 -out ./certs/ca.key 4096
openssl req -passin env:CA_PASS -new -x509 -sha256 -days 365 -key ./certs/ca.key -out ./certs/ca.crt -subj "/CN=ca"

# 2: Generate server private key (server.key)
openssl genrsa -passout env:SERVER_PASS -aes256 -out ./certs/server.key 4096

# 3: Get a certificate signing request from the CA (server.csr)
openssl req -passin env:SERVER_PASS -new -key ./certs/server.key -out ./certs/server.csr -subj ${MY_SUBJECT}

# 4: Self-sign the certificate with the CA that just created (server.crt)
openssl x509 -req -extfile <(printf "subjectAltName=DNS:${SERVER_CN}") -passin env:CA_PASS -sha256 -days 365 -in ./certs/server.csr -CA ./certs/ca.crt -CAkey ./certs/ca.key -set_serial 01 -out ./certs/server.crt

# 5: Convert the server private key to .pem format that required by gRPC (server.pem)
openssl pkcs8 -topk8 -nocrypt -passin env:SERVER_PASS -in ./certs/server.key -out ./certs/server.pem
