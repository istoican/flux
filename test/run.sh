#!/bin/bash

while true
do
    key=$(head /dev/urandom | tr -dc A-Za-z0-9 | head -c 13)
    value=$(head /dev/urandom | tr -dc A-Za-z0-9 | head -c 13)
    echo "PUTTING -> $value for key $key"
    curl -s --data $value $1/$key 2> /dev/null > /dev/null
    sleep $2
    curl -s $1/$key 2> /dev/null > /dev/null
done