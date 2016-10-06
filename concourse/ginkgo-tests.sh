#!/bin/bash

set -e

mkdir -p $GOPATH/src

echo "Moving project to GOPATH..."
prefix_path=$GOPATH/src/github.com/SAP
mkdir -p $prefix_path
cp -r aker-oauth-plugins $prefix_path
cd $prefix_path/aker-oauth-plugins

echo "Fetching test tools..."
go get github.com/onsi/ginkgo/ginkgo

echo "Running tests..."

# setting the timezone is necessary in order
# to get the correct token expirity date
export TZ='Europe/Sofia'

ginkgo -r
