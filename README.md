# sprout
:seedling: host your own home-grown servers! :seedling:

### Why?
a step in the right direction for a more "distributed" internet
empower people to grow it themselves!
promote privacy by taking the opportunity to sell personal data away from cloud providers

### features
* easy to use, plugable,  but extensible...(all skill levels)
* tutorial (empower users)
* monitoring system (visualizations)
* suggestions for what to use it for
  * "cloud" storage (you own your own data)
  * media server
  * host static websites for free
  * host simple webservers for free
  * learn how it all works
  * private chat service for your friends
  * tld server
* show money saver analysis?? lol

### supported LAN network topologies (constellations)
* home - contains reverse proxy to all satellites, contains certs, contains ddclient
* satellite - contains reverse proxy to anything here too
#### types of proxy configurations
*

### deliverables
* cli tool (curses and non-curses)
* desktop app

### components
* control center (controls all components below, unites the APIs)
* OS flasher (sd card or usb drive)
* monitoring system
* provisioner (using ansible)
* service discovery (when a new device is attached to the network)
* types (type defs for servers, services, devices, etc.)
* security hardening?
* gateway port forwarding

### forseeable challenges
* handle configuring routers (make http requests simulating a browser...)

### Tests
#### Ansible
``` shell
$ pip install molecule
$ pip install 'molecule[docker]'
```

### Resources
* [Go Change Gateway with UPnP](https://github.com/NebulousLabs/go-upnp)
* [Go UPnP](https://github.com/huin/goupnp)
