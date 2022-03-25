#!/bin/bash -ex

docker run \
    -d \
    --rm \
    --name db \
    --cap-add=sys_nice \
    -e MYSQL_RANDOM_ROOT_PASSWORD=1 \
    -e MYSQL_DATABASE=checks \
    -e MYSQL_USER=checks \
    -e MYSQL_PASSWORD=checks \
    -p 3306:3306 \
    mysql:8.0
