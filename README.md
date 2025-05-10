# pibackdrop

This is a hacky thing to make a raspberry pi display a set of differet background images on its displays based on a control signal I haven't figured out yet.

The use case is for having the pi connected to a projector that'll display different background images during a live music gig. Album covers, song artwork, etc.

Ideas for the control bits is remotely via phone/ipad using a web app, or using MIDI over usb directly. Haven't figured it out yet.


## Install directions

Install raspbian on a rpi. Use `raspi-config` under "Display" to make sure you're running Wayfire/Wayland.

Install extra wayfire plugins, like the one we use to hide the mouse cursor:
```
sudo apt install libglibmm-2.4-dev libglm-dev libxml2-dev libpango1.0-dev libcairo2-dev wayfire-dev libwlroots-dev libwf-config-dev meson ninja-build libvulkan-dev cmake
git clone https://github.com/seffs/wayfire-plugins-extra-raspbian && cd wayfire-plugins-extra-raspbian
meson build --prefix=/usr --buildtype=release
ninja -C build && sudo ninja -C build install
```

Add the following to the bottom of `~/.config/wayfire.ini`

```
[core]
plugins = \
	autostart \
	hide-cursor

[autostart]
panel = wfrespawn wf-panel-pi
background = wfrespawn pcmanfm --desktop --profile LXDE-pi
xdg-autostart = lxsession-xdg-autostart
chromium = chromium-browser http://localhost:8080 --kiosk --noerrdialogs --disable-infobars --no-first-run --ozone-platform=wayland --enable-features=OverlayScrollbar --start-maximized
reloader = /home/pi/pibackdrop-reloader.sh
screensaver = false
dpms = false

```

Now, populate /home/pi/datadir with images. Reboot.


## Usage

Connect the Pi to an external HDMI source - I also have it connected to a waveshare screen that sits on my pedalboard so I can see what's supposed to be projected.

Right now, the only way to control it is to connect to the same wifi and visit http://rpi:8080/a and click the 'set' links.

Setting up the hotspot:

```
nmcli device wifi hotspot ssid myssid password mypass
nmcli connection # Note the uuid of the 'Hotspot' connection
nmcli connection modify my-hotspot-uuid-from-above-commend connection.autoconnect yes connection.autoconnect-priority 100
```


