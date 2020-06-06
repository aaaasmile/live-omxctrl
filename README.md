# live-omxctrl
live-omxctrl is a web interface for the omx-player on Raspberry P4.
It uses a golang backend with an embedded Vue.js frontened. 
The control of the omx-player is done using the dbus inteface in golang. 
The project is in a very early stage, but still I can turn on the radio swiss classic on Pi4
using an intranet browser on my iphone.

## Development
I am using the golang env direct on the target (https://dl.google.com/go/go1.13.12.linux-armv6l.tar.gz)
because the dbus and the omx-player are on Pi4. Cross compiling the dbus library was not working as expected in windows.
I am using Visual Code with the ssh remote dev extension and it is working very well.

