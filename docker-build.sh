#!/bin/bash
# Build Docker image for code runner

docker build -t golang-code-runner:latest -f internal/code/Dockerfile .

