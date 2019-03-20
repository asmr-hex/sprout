package upnp

import (
	"net/url"

	"github.com/huin/goupnp"
	"github.com/huin/goupnp/dcps/internetgateway1"
)

const (
	// UPnP search target (ST) for Internet Gateway Devices v1
	InternetGatewayDeviceSearchTarget = `urn:schemas-upnp-org:device:InternetGatewayDevice:1`
)

type Device struct {
	goupnp.RootDevice
	Location *url.URL
}

func DiscoverRouters() ([]*Device, error) {
	var (
		routers = []*Device{}
	)

	// discover all internet gateway devices
	maybeRootDevices, err := goupnp.DiscoverDevices(InternetGatewayDeviceSearchTarget)
	if err != nil {
		return nil, err
	}

	for _, maybeRootDevice := range maybeRootDevices {
		routers = append(routers, &Device{
			RootDevice: *maybeRootDevice.Root,
			Location:   maybeRootDevice.Location,
		})
	}

	return routers, nil
}

func (d *Device) GetUPnPServices() ([]*internetgateway1.WANIPConnection1, error) {
	return internetgateway1.NewWANIPConnection1ClientsFromRootDevice(&d.RootDevice, d.Location)
}
