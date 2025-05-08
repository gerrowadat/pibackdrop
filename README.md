# pibackdrop

This is a hacky thing to make a raspberry pi display a set of differet background images on its displays based on a control signal I haven't figured out yet.

The use case is for having the pi connected to a projector that'll display different background images during a live music gig. Album covers, song artwork, etc.

Ideas for the control bits is remotely via phone/ipad using a web app, or using MIDI over usb directly. Haven't figured it out yet.


## Install directions

Install raspbian on a rpi. 

Set kiosk mode (as of May 2025, the officiasl kiosk mode rpi guide doesn't work).

Add this to `/home/pi/.config/labwc/autostart`

```
/home/pi/kiosk.sh
```

The `kiosk.sh` file looks like:

```
sleep 4
/bin/chromium-browser  --kiosk --ozone-platform=wayland --start-maximized --noerrdialogs --disable-infobars --enable-features=OverlayScrollbar  http://localhost:8080/ &
```
