#!/bin/bash

echo "εncooη system test starting."

SPACES='                      '
for PACK in utils configuration model database apis
do
    MSG=`echo "$PACK:$SPACES"|cut -c1-20`
    echo -n "$MSG" && go test $PACK/* -outputdir testData -covermode=count -cover -coverprofile cover_"$PACK"_test.out && go tool cover -html=testData/cover_"$PACK"_test.out  -o testData/cover_"$PACK"_test.html
done

echo "εncooη system test completed."

if [ -n "$ENCOON_KEEP_ALIVE_AFTER_SECONDS" ]; then
    sleep $ENCOON_KEEP_ALIVE_AFTER_SECONDS
fi

echo "εncooη system test ended."
