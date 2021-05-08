#!/bin/sh
/usr/bin/time -f '%Uut %Sst %ert %MkB %C' "$@" < input 1> output 2> result
