#!/usr/bin/env  bash

docker rm -f redis || true > /dev/null
docker run --rm -d -p 6379:6379  --name redis redis:7.2.2
