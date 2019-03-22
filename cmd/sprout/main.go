package main

import (
	"fmt"

	"github.com/connorwalsh/sprout/upnp"
)

func main() {
	routers, err := upnp.DiscoverRouters()
	if err != nil {
		panic(err)
	}

	for _, router := range routers {
		if router.Err != nil {
			fmt.Println("Error probing device at location: " + router.Location.String())
			continue
		}
		fmt.Println("")
		fmt.Println("Name: " + router.Root.Device.FriendlyName)
		fmt.Println("UDN: " + router.Root.Device.UDN)
		fmt.Println("Device Type: " + router.Root.Device.DeviceType)
		fmt.Println("Manufacturer: " + router.Root.Device.Manufacturer)
		fmt.Println("Model Name: " + router.Root.Device.ModelName)
		fmt.Println("Model Number: " + router.Root.Device.ModelNumber)
		fmt.Println("Serial Number: " + router.Root.Device.SerialNumber)
		fmt.Println("Model Description: " + router.Root.Device.ModelDescription)
		fmt.Println("Location: " + router.Location.String())

		// services, err := router.GetUPnPServices()
		// if err != nil {
		// 	panic(err)
		// }

		// fmt.Println(services)
		// b, _ := json.MarshalIndent(services, "", "    ")
		// fmt.Println(string(b))
	}
}
