#!/bin/bash 

echo "Checking..." 
gosec -exclude-generated ./...
