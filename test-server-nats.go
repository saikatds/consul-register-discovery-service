package main

import (
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/nats-io/gnatsd/server"
	"github.com/nats-io/gnatsd/test"
	"github.com/nats-io/go-nats"
	"log"
	"net/http"
	"os"
	"runtime"
)

//port number of consul registration
const Port  = 4222

// RunDefaultServer will run a server on the default port.
func RunDefaultServer() *server.Server {
	return RunServerOnPort(nats.DefaultPort)
}

// RunServerOnPort will run a server on the given port.
func RunServerOnPort(port int) *server.Server {
	opts := test.DefaultTestOptions
	opts.Port = port
	return RunServerWithOptions(opts)
}

// RunServerWithOptions will run a server with the given options.
func RunServerWithOptions(opts server.Options) *server.Server {
	return test.RunServer(&opts)
}

// Registering service with consul
func registerServiceWithConsul() {
	config := consulapi.DefaultConfig()
	consul, err := consulapi.NewClient(config)
	if err != nil {
		log.Fatalln(err)
	}

	registration := new(consulapi.AgentServiceRegistration)

	registration.ID = "nats-service"
	registration.Name = "nats-service"
	address := hostname()
	registration.Address = address

	registration.Port = Port
	registration.Check = new(consulapi.AgentServiceCheck)
	registration.Check.HTTP = fmt.Sprintf("http://%s:%v/healthcheck", address, Port)
	registration.Check.Interval = "5s"
	registration.Check.Timeout = "3s"
	consul.Agent().ServiceRegister(registration)
}

// return the hostname
func hostname() string {
	hn, err := os.Hostname()
	if err != nil {
		log.Fatalln(err)
	}
	return hn
}


// health check handler
func healthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `nats server's health is ok.'`)
}

func main(){

	//Registering Reply Model
	registerServiceWithConsul()
	http.HandleFunc("/healthcheck", healthCheck)
	fmt.Println("user service is up on : ", hostname() , Port)
	http.ListenAndServe(":4222", nil)

	log.Println("nats service Successfully Registered to Consul")

	server:= RunDefaultServer()

	isStarted := server.Start

	log.Println("Server status : ",isStarted)

	//keep the server alive
	runtime.Goexit()

}
