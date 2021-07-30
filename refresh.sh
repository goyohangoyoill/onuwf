#!/bin/bash
while [ 1 ]
do
	pid=`ps -ef | grep "ONUWF" | grep -v 'grep' | awk '{print $2}'`
	if [ -z "$pid" ]; then
		echo "ONUWF refreshed"
		./ONUWF &
	fi
	sleep 2
done
