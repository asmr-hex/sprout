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
		fmt.Println("Name: " + router.Name)
		fmt.Println("UDN: " + router.UDN)
		fmt.Println("Device Type: " + router.Type)
		fmt.Println("Manufacturer: " + router.Manufacturer)
		fmt.Println("Model Number: " + router.ModelNumber)
		fmt.Println("Serial Number: " + router.SerialNumber)
		fmt.Println("Model Description: " + router.Description)
		fmt.Println(router.PortMappingPermitted)

		// services, err := router.GetUPnPServices()
		// if err != nil {
		// 	panic(err)
		// }

		// fmt.Println(services)
		// b, _ := json.MarshalIndent(services, "", "    ")
		// fmt.Println(string(b))
	}

}
