#!/bin/bash
#

RELOAD_FILE=/home/pi/reloadplease
export DISPLAY=:0.0

/home/pi/pibackdrop --datadir=/home/pi/datadir/ --reloadfile=$RELOAD_FILE 2>&1 > /home/pi/pibackdrop.log &

# Start chromium
sleep 4
/bin/chromium-browser  --kiosk --start-maximized --noerrdialogs --disable-infobars --enable-features=OverlayScrollbar  http://localhost:8080/ &

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
