package upnp

import (
	"net/url"

	"github.com/huin/goupnp"
	"github.com/huin/goupnp/dcps/internetgateway1"
)

// UPnP (Universal Plug and Play) is a collection of networking protocols which
// allow for devices on a network to interact with one another in a standardized
// way! This is way cool-- it allows us to programmatically interact with devices
// on the same network without needing to worry much about the make/model of the
// device!
//
// The UPnP architecture is described as follows:
//   Control Points - devices which interact with UPnP devices by searching for
//                    UPnP devices and issuing other UPnP requests (see below).
//   Devices        - UPnP capable devices attached to the network. Devices can
//                    have embedded devices within them. Devices may provide
//                    a variety of UPnP services.
//   Services       - functionality to perform certain actions. different devices
//                    have different services they provide.
//
// UPnP is composed of the following protocols which are implemented as SOAP APIs:
//   Addressing   - facilitates IP address assignment of a connected UPnP device
//                  on the network. All UPnP devices must implement a DHCP client
//                  in order to negotiate with the network DHCP server for IP
//                  assignment. If no DHCP server is present, IP is self assigned.
//   Discovery    - allows a control device to discover UPnP devices on a network.
//                  This uses the Simple Service Discovery Protocol (SSDP).
//   Description  - allows control devices to obtain descriptions of UPnP devices
//                  on the network.
//   Control      - allows control devices to use available UPnP services to perform
//                  actions on a UPnP device.
//   Eventing     - allows control devices to listen for state changes on a UPnP device.
//   Presentation - allows a control point to retrieve a webpage which can be loaded
//                  in a web browser. This gives all UPnP devices the opportunity to
//                  expose a user-friendly interface for humans who don't natively
//                  speak SOAP.
//
// The code here is specifically focused on interacting with internet routers (described
// as InternetGatewayDevices (IGD) in the UPnP device schema) on a home network to enable
// forwarding of specific ports to another device on the local network.

const (
	// UPnP search target (ST) for Internet Gateway Devices v1/v2
	IGDv1 = `urn:schemas-upnp-org:device:InternetGatewayDevice:1`
	IGDv2 = `urn:schemas-upnp-org:device:InternetGatewayDevice:2`
)

// Do i need this?
type Device struct {
	goupnp.RootDevice
	Location *url.URL
	Err      error
}

// performs a multicast discovery of all Internet Gateway Devices (v1 or v2) and
// returns a slice of these root devices.
//
func DiscoverRouters() ([]goupnp.MaybeRootDevice, error) {
	var (
		targets = []string{IGDv1, IGDv2}
		routers = []goupnp.MaybeRootDevice{}
	)

	// discover all internet gateway devices
	for _, target := range targets {
		d, err := goupnp.DiscoverDevices(target)
		if err != nil {
			return nil, err
		}

		routers = append(routers, d...)
	}

	return routers, nil
}

func (d *Device) GetUPnPServices() ([]*internetgateway1.WANIPConnection1, error) {
	return internetgateway1.NewWANIPConnection1ClientsFromRootDevice(&d.RootDevice, d.Location)
}
