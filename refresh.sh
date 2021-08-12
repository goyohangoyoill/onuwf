#!/bin/bash
while [ 1 ]
do
	pid=`ps -ef | grep "onuwf" | grep -v 'grep' | awk '{print $2}'`
	if [ -z "$pid" ]; then
		echo "onuwf refreshed"
		./onuwf &
	fi
	sleep 2
done
