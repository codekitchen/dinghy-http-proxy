package main

import (
	"flag"
	"fmt"

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

	// get actual container ID, in case name is passed in
	container, err := client.InspectContainer(*containerName)
	if err != nil {
		panic(err)
	}
	containerID := container.ID

	currentNetworks := getJoinedNetworks(client, containerID)
	bridgeNetworks := getActiveBridgeNetworks(client, containerID)

	toJoin := getNetworksToJoin(currentNetworks, bridgeNetworks)
	toLeave := getNetworksToLeave(currentNetworks, bridgeNetworks)

	fmt.Printf("currently in %d networks, found %d bridge networks, %d to join, %d to leave\n", len(currentNetworks), len(bridgeNetworks), len(toJoin), len(toLeave))

	for _, id := range toLeave {
		fmt.Printf("leaving network %s\n", id)
		err := client.DisconnectNetwork(id, docker.NetworkConnectionOptions{
			Container: *containerName,
		})
		if err != nil {
			panic(err)
		}
	}

	for _, id := range toJoin {
		fmt.Printf("joining network %s\n", id)
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

func getActiveBridgeNetworks(client *docker.Client, containerName string) (networks map[string]bool) {
	networks = make(map[string]bool)

	allNetworks, err := client.ListNetworks()
	if err != nil {
		panic(err)
	}

	for _, netOverview := range allNetworks {
		if netOverview.Driver == "bridge" {
			// grab the network details (including the list of containers)
			net, err := client.NetworkInfo(netOverview.ID)
			if err != nil {
				panic(err)
			}
			_, containsSelf := net.Containers[containerName]
			if net.Options["com.docker.network.bridge.default_bridge"] == "true" ||
				len(net.Containers) > 1 ||
				(len(net.Containers) == 1 && !containsSelf) {
				networks[net.ID] = true
			}
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
