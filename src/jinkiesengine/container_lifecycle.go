package jinkiesengine

import (
	"context"
	"fmt"
	"github.com/docker/go-connections/nat"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type ContainerInfo struct {
	AutoRemove    bool
	ImageName     string
	ContainerName string
	ContainerPort nat.Port
	HostIp        string
	HostPort      string
}

func RunRunRun(jinkies ContainerInfo) container.ContainerCreateCreatedBody {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	imageName := jinkies.ImageName
	out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}

	defer out.Close()
	io.Copy(os.Stdout, out) // write to stdout

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
	}, &container.HostConfig{
		AutoRemove:   jinkies.AutoRemove,
		PortBindings: nat.PortMap{jinkies.ContainerPort: {{HostIP: jinkies.HostIp, HostPort: jinkies.HostPort}}},
	}, nil, nil, jinkies.ContainerName)
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	fmt.Println(resp.ID)

	return resp
}

func StopGirl(jinkies ContainerInfo) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	if stopErr := cli.ContainerStop(ctx, jinkies.ContainerName, nil); err != nil {
		panic(stopErr)
	}
}
