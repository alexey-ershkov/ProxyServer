#!/bin/sh

openssl genrsa -out ../src/ca.key 2048
openssl req -new -x509 -days 3650 -key ../src/ca.key -out ../src/ca.crt -subj "/CN=alersh proxy CA"
openssl genrsa -out ../src/cert.key 2048
mkdir ../certs/
