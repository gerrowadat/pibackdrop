#!/bin/bash
#

RELOAD_FILE=/home/$USER/reloadplease

# Check for the reload semaphore file, and reload/delete it if present.
while true
do
	sleep 1
	if [ -f $RELOAD_FILE ]; then
		# hit and release shift-F5
		wtype -M shift -P F5 -p f5
		rm $RELOAD_FILE
	fi
done
