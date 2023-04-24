#!/bin/bash -ex

# Prevent core dumps
ulimit -c 0

# Run as ark user
if [ "$(id -u)" = '0' ]; then
    set -- su-exec ark "$@"
fi

# Pulumi login to blob storage
#pulumi login azblob://pulumi 

# Run what's passed to entry point
exec "$@"