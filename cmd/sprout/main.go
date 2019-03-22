package main

import (
	"fmt"

	"github.com/connorwalsh/sprout/upnp"
)

func main() {
	routers, _, err := upnp.DiscoverRouters()
	if err != nil {
		panic(err)
	}

	for _, router := range routers {
		fmt.Println("")
		fmt.Println("Name: " + router.Device.FriendlyName)
		fmt.Println("UDN: " + router.Device.UDN)
		fmt.Println("Device Type: " + router.Device.DeviceType)
		fmt.Println("Manufacturer: " + router.Device.Manufacturer)
		fmt.Println("Model Name: " + router.Device.ModelName)
		fmt.Println("Model Number: " + router.Device.ModelNumber)
		fmt.Println("Serial Number: " + router.Device.SerialNumber)
		fmt.Println("Model Description: " + router.Device.ModelDescription)
		fmt.Println("Location: " + router.URLBaseStr)

		// services, err := router.GetUPnPServices()
		// if err != nil {
		// 	panic(err)
		// }

		// fmt.Println(services)
		// b, _ := json.MarshalIndent(services, "", "    ")
		// fmt.Println(string(b))
	}

}
