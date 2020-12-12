# live-omxctrl
Live-omxctrl is a web interface for playing music and urls on Raspberry P4.
This application can play files stored on the Pi4 device and also links provided by the tube.  
Simply copy an URl from the tube app and paste it into this Live Music app. Playing files and playlist from the flash card is from yesterday.

At the end it is also another simple Web UI for omxplayer, but with some small 
features more then usually.

Playing a chill out music from the tube controlled by the iphone is pretty easy:

![alt text](https://github.com/aaaasmile/live-omxctrl/blob/master/doc/05-12-_2020_22-23-43.png?raw=true)

## Use Case
Turn on the Pi4 and turn on the HiFi connected to the Pi4 via hdmi cable.
In your Iphone browser goto the Pi4 Ip Address with :5548/omx/#/ suffix and paste a tube URL. Tap the icon _play url_ and enjoy the music. 
Before turning off the Pi4 device, in OS view tap the shutdown icon and wait a little before power turn off.

## OS View
After destroying a couple of sd cards turning off the Pi4 using an
hard power shutdown, I added a software shutdown functionality to this app. 


![alt text](https://github.com/aaaasmile/live-omxctrl/blob/master/doc/05-12-_2020_22-56-38.png?raw=true)

_Reboot_ and a _Service Restart_ are also available. 
For the _Service Restart_ the app should be configured as a service called _live-omxctrl.service_. The service owner
should be able to run sudo commands without asking a password.


## Development
This app uses omxplayer in background with a golang backend, an embedded Vue.js frontend and Websocket for updating the state.   
The omxplayer is controlled using the dbus inteface. 

The golang environment is set up directly on the target (https://dl.google.com/go/go1.13.12.linux-armv6l.tar.gz).
In this way, the dbus and the omx-player are on Pi4. Cross compiling the dbus library was not working as expected in windows.
With Visual Code and the ssh remote dev extension, the development is pretty the same as in the windows manchine. 

Before starting the application with:
```go run main.go```

the config file _config.toml_ should be changed.
```
ServiceURL=<pi4-IP:portapp>
TmpInfo=<temp directory to store the url description>
```

The UI is developed using Vue.js and Vuetify directly on the golang backend without setting up a node.js environment. 

Sound strange? It works.  

The content of files with the .vue extension is copied inside the same
file named with the .js extension, inside the template section.  
This could be done automatically in Visual Code 
using the extension [vuetempltojs](https://github.com/aaaasmile/vuetempltojs). On Pi4 the exetnsion will not works, so the command should be called manually.
For example, after changing the file _playerbar.vue_ the playerbar.js is updated with:

```
igors@pi4:~/app/go/live-omxctrl $ /usr/local/bin/vuetojs.bin -vue ./static/js/vue/components/playerbar.vue 
```

If you don't like it, feel free to setup a separate Vue project inside a node.js
environment that wil run inside a golang backend developed at the same time.

