package main

import (
	"flag"

	docker "github.com/fsouza/go-dockerclient"
)

var (
	endpoint = "unix:///tmp/docker.sock"
)

func main() {
	var containerName = flag.String("container-name", "", "the name of this docker container")
	flag.Parse()

	client, err := docker.NewClient(endpoint)
	if err != nil {
		panic(err)
	}

	currentNetworks := getJoinedNetworks(client, *containerName)
	bridgeNetworks := getBridgeNetworks(client)

	for _, id := range bridgeNetworks {
		if !currentNetworks[id] {
			err := client.ConnectNetwork(id, docker.NetworkConnectionOptions{
				Container: *containerName,
			})
			if err != nil {
				panic(err)
			}
		}
	}
}

func getJoinedNetworks(client *docker.Client, containerName string) (networks map[string]bool) {
	networks = make(map[string]bool)

	container, err := client.InspectContainer(containerName)
	if err != nil {
		panic(err)
	}

	for _, net := range container.NetworkSettings.Networks {
		networks[net.NetworkID] = true
	}

	return
}

func getBridgeNetworks(client *docker.Client) (ids []string) {
	networks, err := client.ListNetworks()
	if err != nil {
		panic(err)
	}

	for _, net := range networks {
		if net.Driver == "bridge" {
			ids = append(ids, net.ID)
		}
	}
	return
}
