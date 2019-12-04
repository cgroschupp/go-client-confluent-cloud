# go-client-confluent-cloud


```golang
package main

import (
	"fmt"
	"log"

	"github.com/cgroschupp/go-client-confluent-cloud/confluentcloud"
)

func main() {
	client := confluentcloud.NewClient("<EMAIL>", "<PASSWORD>")
	err := client.Login()
	if err != nil {
		log.Print(err)
		return
	}

	userData, err := client.Me()
	if err != nil {
		log.Print(err)
		return
	}

	clusters, err := client.ListClusters(userData.Account.ID)
	if err != nil {
		log.Print(err)
		return
	}

	for _, cluster := range clusters {
		fmt.Println(cluster.ID)
	}
}
```