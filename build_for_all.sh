#!/bin/bash

export GOPATH=$GOPATH:~/project/golang/go_crawler

echo $GOPATH

rm -f 163_comment 163_comment.exe

echo -n "Building for ubuntu..."
go build -o 163_comment src/crawler/main.go
echo "Done"

echo -n "Building for macOS..."
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64
go build -o 163_comment src/crawler/main.go
echo "Done"

echo -n "Building for Windows..."
CGO_ENABLED=0 GOOS=windows GOARCH=amd64
go build -o 163_comment.exe src/crawler/main.go
echo "Done"

