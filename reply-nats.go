package main

import (
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/nats-io/go-nats"
	"log"
	"runtime"

)


func lookupServiceWithConsul1(serviceName string) (string, error) {
	config := consulapi.DefaultConfig()
	consul, err := consulapi.NewClient(config)
	if err != nil {
		return "", err
	}
	services, err := consul.Agent().Services()
	if err != nil {
		return "", err
	}
	srvc := services[serviceName]
	address := srvc.Address
	port := srvc.Port
	return fmt.Sprintf("http://%s:%v", address, port), nil
}

func main() {


	natsLookupUrl, err := lookupServiceWithConsul1("nats-service")
	fmt.Println("URL of NATS in Consul: ", natsLookupUrl)
	if err != nil {
		log.Fatal("Lookup failed : ",err)
	}

	// connecting to nats server
	natConnection  , error := nats.Connect(nats.DefaultURL)

	if error != nil {
		log.Fatal("Could not connect to nats server due to ",error)
	}else{
		log.Println("Successfully connected to nats server with URL ",nats.DefaultURL)
	}

	//declaring subject and reply
	sub := "saikat"

	//Subscribing to the subject
	natConnection.Subscribe(sub , func(msg *nats.Msg) {
		log.Println("Message received  for subject ", sub , " is : " ,string(msg.Data ))
		//giving reply
		natConnection.Publish(msg.Reply , []byte (string(msg.Data)+" EFG") )
	})

	log.Println("Listening to subject ", sub)

	// Keep the connection alive
	runtime.Goexit()
}