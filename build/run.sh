#!/bin/bash

caseNum=$1
tlimit=$2
for((i=1;i<=$caseNum;i++))
do
  timeout $(expr $tlimit / 1000)s ./main < input-$i 1> output-$i 2>> rte
  if [ $? -ne 0 ]; then
    echo 'Runtime Error occured at test case '$i
    exit $i
  fi
done
