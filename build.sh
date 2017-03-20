#!/bin/bash

REPO_NAME="Compel-Prediction"

echo "Building $REPO_NAME"

if go build -o $GOPATH/bin/$REPO_NAME -i github.com/adhuri/$REPO_NAME ;then
echo "+Successful"
else echo "-Failed"
fi


