#!/bin/bash
pushd ~/src/fm3
protoc --java_out=java proto/rpc.proto
popd
