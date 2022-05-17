package main

import (
	"fmt"
	"log"

	"github.com/icepie/miio.go"
)

func main() {
	client := miio.New("192.168.1.20").SetToken("547d929b571f36ec2a639e520563530d")
	// client.SetToken() // will try to use token from handshake if not set

	// https://home.miot-spec.com/spec/mmgg.pet_waterer.s1
	payload := []map[string]interface{}{
		{"id": 18, "did": "switch", "siid": 5, "piid": 1, "value": 1},
	}

	resp, err := client.GetProperties(payload)
	if err != nil {
		log.Println(client.Token(), err)
	}
	fmt.Printf("%s\n", resp)

	resp, err = client.GetProperties(payload)
	if err != nil {
		log.Println(client.Token(), err)
	}
	fmt.Printf("%s\n", resp)

	// //switch:toggle
	// resp, err = client.Action(2, 1, []interface{}{})
	// if err != nil {
	// 	log.Println(client.Token(), err)
	// }

}
