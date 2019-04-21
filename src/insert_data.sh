#!/bin/bash

i=0
j=0

while true
do  
        echo "172.0.0.${j} - - [04/Mar/2018:13:49:52 +0000] http \"GET /foo?query=t HTTP/1.0\" 200 ${j} \"-\" \"KeepAliveClient\" \"-\" 1.005 1.854" >> "./access.log"
        let "j++"
        echo $j
done
