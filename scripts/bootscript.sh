#!/bin/bash

# run nginx
/usr/sbin/nginx

air  --build.cmd "go build -o bin/agent cmd/agent.go"  --build.bin "./bin/agent"