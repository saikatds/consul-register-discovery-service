package main

import (
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/nats-io/go-nats"
	"log"
	"runtime"
	"time"
)


func lookupServiceWithConsul(serviceName string) (string, error) {
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

func main (){

	natsLookupUrl, err := lookupServiceWithConsul("nats-service")
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

	//closing the connection after getting the reply
	defer natConnection.Close()

	sub , request := "saikat" , " 1123456 "

	//Requesting to the nats server and getting the reply
	reply , err := natConnection.Request(sub , []byte (request) , 100*time.Millisecond)

	//Checking if any error coming
	if err != nil {
		log.Fatal("Req-Reply model broke down due to ",err)
	}

	// Request and Reply data
	log.Println(" Request : ", string(request))
	log.Println("Reply : ", string(reply.Data))

	runtime.Goexit()
}

