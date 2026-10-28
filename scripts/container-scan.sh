#!/bin/bash 

echo "Building..."
docker build -t instay-be:local . 
echo "Scanning..." 
docker run --rm -v /var/run/docker.sock:/var/run/docker.sock -v ~/.cache:/root/.cache aquasec/trivy:0.69.3 image instay-be:local
