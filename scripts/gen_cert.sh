#!/bin/bash

HOST=$1
PRIVATE_KEY_FILE=$HOST.key
CERT_FILE=$HOST.crt

cat > cert.conf <<EOF
subjectAltName = @alt_names

[alt_names]
DNS.1 = $1

EOF

openssl genrsa -out $PRIVATE_KEY_FILE 2048
openssl req -new -key $PRIVATE_KEY_FILE -subj "/CN=$HOST" -addext "subjectAltName = DNS:$HOST, DNS:www.$HOST" -sha256 |
openssl x509 -req -out $CERT_FILE -days 3650 -CA ca.crt -CAkey ca.key -set_serial "$(date +%Y%m%d%H%M%S)" -extfile cert.conf
