#!/bin/zsh
openssl req -new -key $3/cert.key -subj "/CN=$1" -sha256 | openssl x509 -req -days 3650 -CA $3/ca.crt -CAkey $3/ca.key -set_serial "$2" > $4$1.crt
