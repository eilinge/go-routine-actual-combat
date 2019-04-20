#!/bin/bash

i=0
j=0

while (( $i < 1000 ))
do  
    
    while (( $j < 250 ))
    do
        echo "172.0.0.${j} - - [04/Mar/2018:13:49:52 +0000] http \"GET /foo?query=t HTTP/1.0\ 200 ${i} \"-\" \"KeepAliveClient\" \"-\" 1.005 1.854" >> "./access.log"
        let "j++"
        # echo $j
    done
    let "i++"
done