package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

func main() {
	fmt.Println("Welcome to docker client")
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()
out:
	for {
		showOptions()
		option, err := readOption()
		if err != nil {
			fmt.Println("Invalid option")
			continue
		}

		switch option {
		case 0:
			fmt.Println("Leaving docker client")
			break out
		case 1:
			listContainers(cli)
		case 2:
			listImages(cli)
		case 3:
			fmt.Println("Enter container id")
			containerId, err := readOptionUsingFmt()
			if err != nil {
				fmt.Println("Invalid option")
				continue
			}
			stopContainer(containerId, cli)
		case 4:
			fmt.Println("Enter container id")
			containerId, err := readOptionUsingFmt()
			if err != nil {
				fmt.Println("Invalid option")
				continue
			}
			startContainer(containerId, cli)
		case 5:
			fmt.Println("Enter container id")
			containerId, err := readOptionUsingFmt()
			if err != nil {
				fmt.Println("Invalid option")
				continue
			}
			removeContainer(containerId, cli)
		case 6:
			fmt.Println("Enter image id")
			imageId, err := readOptionUsingFmt()
			if err != nil {
				fmt.Println("Invalid option")
				continue
			}
			removeImage(imageId, cli)
		default:
			fmt.Println("Invalid option")
		}
	}

}

func showOptions() {
	fmt.Println("Choose the option you want:")
	fmt.Println("1 - list containers")
	fmt.Println("2 - list images")
	fmt.Println("3 - stop container")
	fmt.Println("4 - start container")
	fmt.Println("5 - remove container")
	fmt.Println("6 - remove image")
	fmt.Println("0 - leave")
}

func readOption() (int, error) {
	scanner := bufio.NewReader(os.Stdin)
	input, err := scanner.ReadString('\n')
	if err != nil {
		return 0, err
	}
	input = strings.Replace(input, "\n", "", -1)
	option, err := strconv.Atoi(input)
	if err != nil {
		return 0, err
	}
	return option, nil
}

func readOptionUsingFmt() (string, error) {
	var input string
	_, err := fmt.Scan(&input)
	if err != nil {
		return "", err
	}
	return input, nil
}

func listContainers(client *client.Client) {
	containers, err := client.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		fmt.Println("Error looking for docker containers")
		return
	}

	if len(containers) == 0 {
		fmt.Println("There are no containers running!")
		return
	}

	for _, ctr := range containers {
		fmt.Printf("%s %s (status: %s)\n", ctr.ID, ctr.Image, ctr.Status)
	}
}

func stopContainer(containerId string, client *client.Client) {
	err := client.ContainerStop(context.Background(), containerId, container.StopOptions{})
	if err != nil {
		fmt.Printf("Error stopping container %s\n", containerId)
		fmt.Printf("Error: %s\n", err)
		return
	}
	fmt.Printf("Stopped container %s successfuly\n", containerId)
}

func startContainer(containerId string, client *client.Client) {
	err := client.ContainerStart(context.Background(), containerId, container.StartOptions{})
	if err != nil {
		fmt.Printf("Error starting container %s\n", containerId)
		fmt.Printf("Error: %s\n", err)
		return
	}
	fmt.Printf("Start container %s successfuly\n", containerId)
}

func removeContainer(containerId string, client *client.Client) {
	err := client.ContainerRemove(context.Background(), containerId, container.RemoveOptions{})
	if err != nil {
		fmt.Printf("Error removing container %s\n", containerId)
		fmt.Printf("Error: %s\n", err)
		return
	}
	fmt.Printf("Container %s removed successfuly\n", containerId)
}

func listImages(client *client.Client) {
	images, err := client.ImageList(context.Background(), image.ListOptions{All: true})
	if err != nil {
		fmt.Println("Error looking for docker images")
		return
	}
	if len(images) == 0 {
		fmt.Println("There are no images in this machine!")
		return
	}

	for _, img := range images {
		fmt.Printf("ID %s with tags %s\n", img.ID, img.RepoTags)
	}
}

func removeImage(imageId string, client *client.Client) {
	_, err := client.ImageRemove(context.Background(), imageId, image.RemoveOptions{})
	if err != nil {
		fmt.Printf("Error removing image %s\n", imageId)
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Image %s removed\n", imageId)
}
