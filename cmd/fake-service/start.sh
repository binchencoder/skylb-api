#!/bin/sh

# Assuming binary fake-service is under the same dir as start.sh.

cd `dirname $0`

# skylb endpoints
export SKYLB=skylbserver:1900,skylbserver:1900

# Inform skylb of the address of real services.
./fake-service --skylb-endpoints=$SKYLB \
  -debug-svc-endpoint=vexillary-service=192.168.10.41:4100 "$@" \
  -v=3 -alsologtostderr \
  &

