#!/bin/bash

openssl genrsa -out server.key 2048
openssl ecparam -genkey -name secp384r1 -out ./certs/server.key
openssl req -nodes -subj '/CN=localhost' -new -x509 -sha256 -key ./certs/server.key -out ./certs/server.crt -days 365
