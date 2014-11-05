#!/bin/bash
pushd java
protoc --java_out=. jp/ne/sakura/ogiekako/fm3/rpc/proto/rpc.proto
popd
