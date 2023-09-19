#!/bin/bash

HOST=$1
PRIVATE_KEY_FILE=$HOST.key
CERT_FILE=$HOST.crt

openssl genrsa -out $PRIVATE_KEY_FILE 2048
cat v3.ext | sed s/%%DOMAIN%%/"$HOST"/g > /tmp/__v3.ext
openssl req -new -key $PRIVATE_KEY_FILE -subj "/C=CA/ST=None/L=NB/O=None/CN=$HOST" -sha256 -config /tmp/__v3.ext |\
openssl x509 -req -out $CERT_FILE -days 3650 -CA ca.crt -CAkey ca.key -CAcreateserial -extfile /tmp/__v3.ext
