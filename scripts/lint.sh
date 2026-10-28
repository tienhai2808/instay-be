#!/bin/bash 

echo "Formatting..." 
go fmt ./... 
echo "Linting..." 
golangci-lint run