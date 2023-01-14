#!/bin/bash

SPACES='                      '
for PACK in utils configuration model database apis
do
    MSG=`echo "$PACK:$SPACES"|cut -c1-20`
    echo -n "$MSG" && go test $PACK/* -outputdir testData -covermode=count -cover -coverprofile cover_"$PACK"_test.out
done

if [ -n "$ENCOON_KEEP_ALIVE_AFTER_SECONDS" ]; then
    sleep $ENCOON_KEEP_ALIVE_AFTER_SECONDS
fi
