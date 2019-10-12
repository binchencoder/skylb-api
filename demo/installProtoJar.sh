#!/bin/sh

# Run this script to prepare proto jars for java demo.

cd `dirname $0`
cd ../..

bazel build skylb-api/cmd/demo/proto/greeting
bazel build ease-gateway/gateway/options/options

VER=1.0.0-SNAPSHOT

set -x

mvn install:install-file \
-DgroupId=com.binchencoder.skylb \
-DartifactId=demo-proto \
-Dversion=$VER \
-Dfile=bazel-bin/skylb-api/cmd/demo/proto/libgreeting.jar \
-Dpackaging=jar \
-DgeneratePom=true

mvn install:install-file \
-DgroupId=com.binchencoder.api \
-DartifactId=ease-gateway-option \
-Dversion=$VER \
-Dfile=bazel-bin/ease-gateway/gateway/options/liboptions.jar \
-Dpackaging=jar \
-DgeneratePom=true