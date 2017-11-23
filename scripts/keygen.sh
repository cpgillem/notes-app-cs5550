#!/bin/bash
PATH=$1

/bin/rm -f "$PATH"/keys/*
/usr/bin/ssh-keygen -t rsa -b 4096 -f "$PATH"/app.rsa
/usr/bin/openssl rsa -in "$PATH"/app.rsa -pubout -outform PEM -out "$PATH"/app.rsa.pub
# https://gist.github.com/ygotthilf/baa58da5c3dd1f69fae9
