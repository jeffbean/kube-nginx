package main

import (
	"fmt"
	"os"
	"text/template"

	"github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/client"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/fields"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/labels"
)

type ServerPort struct {
	ServerAddress string
	ServicePort   int
}

type ServiceServers struct {
	ServiceName string
	ServerPorts []ServerPort
}

func main() {
	config := &client.Config{
		Host:     "http://localhost:8080",
		Username: "test",
		Password: "password",
	}
	client, err := client.New(config)
	if err != nil {
		// handle error
	}
	//TODO: have namespace a flag for the script
	ns := api.NamespaceDefault
	list, err := client.Services(ns).List(labels.Everything())
	if err != nil {
		fmt.Println(err)
	}
	//TODO: allow selction of nodes by labels or fields
	nodes, err := client.Nodes().List(labels.Everything(), fields.Everything())
	if err != nil {
		fmt.Println(err)
	}
	//TODO: this data structure is a mess. Maybe put in Kubernetes
	// client helper methods for some of these data points
	var templateServiceServers []ServiceServers
	for _, service := range list.Items {
		if service.Name == "kubernetes" {
			// TODO: better way of filtering out the kubernetes service?
			continue
		}
		s := &ServiceServers{ServiceName: service.Name}
		sp := &ServerPort{}
		for _, port := range service.Spec.Ports {
			sp.ServicePort = port.NodePort
		}
		for _, node := range nodes.Items {
			for _, address := range node.Status.Addresses {
				sp.ServerAddress = address.Address
			}
		}
		s.ServerPorts = append(s.ServerPorts, *sp)

		templateServiceServers = append(templateServiceServers, *s)
	}

	// TODO: allow your own template file
	t := template.New("Nginx config template")
	t, err = t.ParseFiles("nginx.tmpl")
	if err != nil {
		fmt.Println(err)
	}
	t.ExecuteTemplate(os.Stdout, "nginx.tmpl", templateServiceServers)
}
