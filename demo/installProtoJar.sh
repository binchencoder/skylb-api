#!/bin/sh

# Run this script to prepare proto jars for java demo.

cd `dirname $0`
cd ../..

bazel build skylb-api/cmd/demo/proto/greeting
bazel build janus/gateway/options/options

VER=1.0.0-SNAPSHOT

set -x

mvn install:install-file \
-DgroupId=com.jingoal.skylb \
-DartifactId=demo-proto \
-Dversion=$VER \
-Dfile=bazel-bin/skylb-api/cmd/demo/proto/libgreeting.jar \
-Dpackaging=jar \
-DgeneratePom=true

mvn install:install-file \
-DgroupId=com.jingoal.api \
-DartifactId=janus-gateway-option \
-Dversion=$VER \
-Dfile=bazel-bin/janus/gateway/options/liboptions.jar \
-Dpackaging=jar \
-DgeneratePom=true

