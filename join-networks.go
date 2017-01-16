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
	bridgeNetworks := getActiveBridgeNetworks(client)

	toJoin := getNetworksToJoin(currentNetworks, bridgeNetworks)
	toLeave := getNetworksToLeave(currentNetworks, bridgeNetworks)

	for _, id := range toLeave {
		err := client.DisconnectNetwork(id, docker.NetworkConnectionOptions{
			Container: *containerName,
		})
		if err != nil {
			panic(err)
		}
	}

	for _, id := range toJoin {
		err := client.ConnectNetwork(id, docker.NetworkConnectionOptions{
			Container: *containerName,
		})
		if err != nil {
			panic(err)
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

func getActiveBridgeNetworks(client *docker.Client) (networks map[string]bool) {
	networks = make(map[string]bool)

	allNetworks, err := client.ListNetworks()
	if err != nil {
		panic(err)
	}

	for _, net := range allNetworks {
		if net.Driver == "bridge" && len(net.Containers) > 0 {
			networks[net.ID] = true
		}
	}

	return
}

func getNetworksToJoin(currentNetworks map[string]bool, bridgeNetworks map[string]bool) (ids []string) {
	for id := range bridgeNetworks {
		if !currentNetworks[id] {
			ids = append(ids, id)
		}
	}

	return
}

func getNetworksToLeave(currentNetworks map[string]bool, bridgeNetworks map[string]bool) (ids []string) {
	for id := range currentNetworks {
		if !bridgeNetworks[id] {
			ids = append(ids, id)
		}
	}

	return
}
