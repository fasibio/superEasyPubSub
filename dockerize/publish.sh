#!/bin/bash
# $1 = version
docker login -u $dockerhubuser -p $dockerhubpassword
docker push fasibio/supereasypubsub:$1