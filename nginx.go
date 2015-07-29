package main

import (
	"fmt"

	"github.com/GoogleCloudPlatform/kubernetes/pkg/client"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/fields"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/labels"
)

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
	list, err := client.Nodes().List(labels.Everything(), fields.Everything())
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("list: %#v", list)
}
