package main

import (
	"encoding/json"
	"fmt"

	"github.com/connorwalsh/sprout/upnp"
)

func main() {
	routers, err := upnp.DiscoverRouters()
	if err != nil {
		panic(err)
	}

	for _, router := range routers {
		services, err := router.GetUPnPServices()
		if err != nil {
			panic(err)
		}

		fmt.Println(services)
		b, _ := json.MarshalIndent(services, "", "    ")
		fmt.Println(string(b))
	}
}
