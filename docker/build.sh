#!/bin/bash

# First parameter is the version of the backend build

docker build  -t docker.pkg.github.com/unizar-30226-2021-11/backend/backend:"$1" .
docker push docker.pkg.github.com/unizar-30226-2021-11/backend/backend:"$1" 
