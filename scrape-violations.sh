#!/usr/bin/env bash

DOTNUM=$1
# downlaod site html
curl -s https://ai.fmcsa.dot.gov/SMS/Carrier/${DOTNUM}/BASIC/UnsafeDriving.aspx > scrape-violations.html
# find starting line number
STARTLN=`grep -o ' <th class="violcode" data-sort="Code">Violations' -n scrape-violations.html | awk -F ":" '{print $1}'`
let STARTLN-=1
echo $STARTLN
# find ending line number
IFS=' ' read -a ARR <<< $(grep -o '</tr>' -n scrape-violations.html | awk -F ":" '{print $1}')

for N in "${ARR[@]}"; do
    if [ $N -gt $STARTLN ]; then
        ENDLN=$N
        break
    fi
done

echo $ENDLN

#extract data between line numbers

HTML=$(sed -n ${STARTLN},${ENDLN}p scrape-violations.html)

echo "$HTML"

STARTLN=`grep -o ' <tr class="violSummary ">' -n scrape-violations.html | awk -F ":" '{print $1}'`
echo $STARTLN

# find ending line number
IFS=' ' read -a ARR <<< $(grep -o '</tr>' -n scrape-violations.html | awk -F ":" '{print $1}')

for N in "${ARR[@]}"; do
    if [ $N -gt $STARTLN ]; then
        ENDLN=$N
        break
    fi
done

echo $ENDLN

HTML2=$(sed -n ${STARTLN},${ENDLN}p scrape-violations.html)

echo "$HTML2"
