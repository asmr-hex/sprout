package upnp

import (
	"fmt"
	"net/url"

	"github.com/huin/goupnp"
	"github.com/huin/goupnp/dcps/internetgateway1"
	"github.com/huin/goupnp/dcps/internetgateway2"
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

// An existing device on the network which is encountered an error while
// probing.
type UnreachableDevice struct {
	Location *url.URL
	Err      error
}

type WANService interface {
	GetExternalIPAddress() (string, error)
	AddPortMapping(string, uint16, string, uint16, string, bool, string, uint32) error
	DeletePortMapping(string, uint16, string) error
	GetGenericPortMappingEntry(uint16) (string, uint16, string, uint16, string, bool, string, uint32, error)
}

type Router struct {
	Name                 string
	Manufacturer         string
	UDN                  string // unique device name
	Type                 string // device type (InternetGatewayDevice:1?2)
	ModelNumber          string
	SerialNumber         string
	Description          string
	PortMappingPermitted bool

	root    *goupnp.RootDevice
	service WANService
}

func (r *Router) discoverWANService() error {
	switch r.Type {
	case IGDv1:
		// attempt to find WANIPConnection services
		wanIPConnectionServices, err := internetgateway1.NewWANIPConnection1ClientsFromRootDevice(r.root, &r.root.URLBase)
		if err != nil {
			// TODO (cw|3.23.2019) looks like this throws if it can't find any devices
			return err
		}

		// if we found a WANIPConnection service, we prefer to just use that
		if len(wanIPConnectionServices) != 0 {
			// TODO (cw|3.23.2019) how many services are we dealing with here??
			r.service = wanIPConnectionServices[0]
			r.PortMappingPermitted = true

			return nil
		}

		// otherwise, lets search for WANPPPConnection services
		wanPPPConnectionServices, err := internetgateway1.NewWANPPPConnection1ClientsFromRootDevice(r.root, &r.root.URLBase)
		if err != nil {
			// TODO (cw|3.23.2019) looks like this throws if it can't find any devices
			return err
		}

		if len(wanPPPConnectionServices) != 0 {
			// TODO (cw|3.23.2019) how many services are we dealing with here??
			r.service = wanPPPConnectionServices[0]
			r.PortMappingPermitted = true

			return nil
		}

		// hmm, i guess we didn't find anything...i guess this thing isn't permitted to
		// perform port mapping...?
		// for now, let's not return an error, but just keep the PortMappingPermitted
		// field as false.
	case IGDv2:
		// attempt to find WANIPConnection services
		wanIPConnectionServices, err := internetgateway2.NewWANIPConnection2ClientsFromRootDevice(r.root, &r.root.URLBase)
		if err != nil {
			// TODO (cw|3.23.2019) looks like this throws if it can't find any devices
			return err
		}

		// if we found a WANIPConnection service, we prefer to just use that
		if len(wanIPConnectionServices) != 0 {
			// TODO (cw|3.23.2019) how many services are we dealing with here??
			r.service = wanIPConnectionServices[0]
			r.PortMappingPermitted = true

			return nil
		}

		// otherwise, lets search for WANPPPConnection services
		wanPPPConnectionServices, err := internetgateway2.NewWANPPPConnection1ClientsFromRootDevice(r.root, &r.root.URLBase)
		if err != nil {
			// TODO (cw|3.23.2019) looks like this throws if it can't find any devices
			return err
		}

		if len(wanPPPConnectionServices) != 0 {
			// TODO (cw|3.23.2019) how many services are we dealing with here??
			r.service = wanPPPConnectionServices[0]
			r.PortMappingPermitted = true

			return nil
		}

		// hmm, i guess we didn't find anything...i guess this thing isn't permitted to
		// perform port mapping...?
		// for now, let's not return an error, but just keep the PortMappingPermitted
		// field as false.
	default:
	}

	return nil
}

func (r *Router) GetExternalIPAddress() (string, error) {
	return "", nil
}

func (r *Router) AddPortMapping() error {
	return nil
}

func (r *Router) RemovePortMapping() error {
	return nil
}

// performs a multicast discovery of all Internet Gateway Devices (v1 or v2) and
// returns a Unique Device Name (string) -> Device map, a slice of unreachable
// devices, and an error.
//
func DiscoverRouters() ([]*Router, []*UnreachableDevice, error) {
	var (
		targetDeviceTypes  = []string{IGDv1, IGDv2}
		routers            = []*Router{}
		unreachable        = []*UnreachableDevice{}
		encounteredDevices = map[string]bool{}
	)

	// discover all internet gateway devices
	for _, targetDeviceType := range targetDeviceTypes {
		devices, err := goupnp.DiscoverDevices(targetDeviceType)
		if err != nil {
			return nil, nil, err
		}

		// include each device into the slice of unique routers
		for _, device := range devices {
			// if the device exists, but encountered an error, it is unreachable
			if device.Err != nil {
				unreachable = append(
					unreachable,
					&UnreachableDevice{
						Location: device.Location,
						Err:      device.Err,
					},
				)

				continue
			}

			if _, ok := encounteredDevices[device.Root.Device.UDN]; ok {
				// we've already encountered this device, skip ahead
				continue
			}

			encounteredDevices[device.Root.Device.UDN] = true

			router := &Router{
				Name:         device.Root.Device.FriendlyName,
				Manufacturer: device.Root.Device.Manufacturer,
				UDN:          device.Root.Device.UDN,
				Type:         device.Root.Device.DeviceType,
				ModelNumber:  device.Root.Device.ModelNumber,
				SerialNumber: device.Root.Device.SerialNumber,
				Description:  device.Root.Device.ModelDescription,
				root:         device.Root,
			}

			// get service
			err := router.discoverWANService()
			if err != nil {
				return nil, nil, err
			}

			// remoteHost, externalPort, _, internalPort, internalClient, enabled, description, duration, err := router.service.GetGenericPortMappingEntry(2)
			// if err != nil {
			// 	return nil, nil, err
			// }
			// fmt.Println(remoteHost)
			// fmt.Println(externalPort)
			// fmt.Println(internalPort)
			// fmt.Println(internalClient)
			// fmt.Println(enabled)
			// fmt.Println(description)
			// fmt.Println(duration)

			err = router.service.AddPortMapping(
				"",
				8080,
				"TCP",
				8080,
				"jackrabbitspal.us",
				true,
				"idk",
				996881901,
			)
			if err != nil {
				fmt.Println("COULTN'T PORT FORWARD")
				return nil, nil, err
			}
			// fmt.Println(remoteHost)
			// fmt.Println(externalPort)
			// fmt.Println(internalPort)
			// fmt.Println(internalClient)
			// fmt.Println(enabled)
			// fmt.Println(description)
			// fmt.Println(duration)

			routers = append(routers, router)
		}
	}

	return routers, unreachable, nil
}
