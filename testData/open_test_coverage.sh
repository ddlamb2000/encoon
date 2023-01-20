#!/bin/bash

for PACK in utils configuration model database apis
do
    go tool cover -html=testData/cover_"$PACK"_test.out
done
