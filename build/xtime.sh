#!/bin/sh
rm -rf result > /dev/null
rm -rf rte > /dev/null
rm -rf build_result > /dev/null
rm -rf output-* > /dev/null
rm -rf wa_* > /dev/null

./time -f '%Uut %Sst %ert %MkB %C' "./run.sh" $1 $2 2> result
