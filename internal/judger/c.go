package judger

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/pkg/errors"
)

type c struct{}

func (c *c) Start() error {
	err := c.prepare()
	if err != nil {
		return errors.Wrap(err, "Failed to prepare for judging C code.")
	}
	fmt.Println("compiling user's code now.")
	err = c.compile()

	if err != nil {
		return errors.Wrap(err, "Failed to compile user's code.")
	}
	fmt.Println("Success.")

	fmt.Println("running user's code now.")
	err = c.run()
	if err != nil {
		return errors.Wrap(err, "Failed to run C code.")
	}

	return nil
}

func (c *c) prepare() error {
	fmt.Println("Preparing judger containers now...")

	imgPulled := false
	conCreated := false

	// Create new docker client instant
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	if err != nil {
		return errors.Wrap(err, "Failed to create cli instance.")
	}

	// Specify the image for judging
	imageName := "miata/goj-judger-c-cpp-img:latest"

	// Check if the image already exist
	imgListSummary, err := cli.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		return errors.Wrap(err, "Failed to get image list.")
	}

	for _, item := range imgListSummary {
		if len(item.RepoTags) > 0 {
			for _, tag := range item.RepoTags {
				if tag == "miata/goj-judger-c-cpp-img:latest" {
					fmt.Println("image already exist:", tag)
					imgPulled = true
					break
				}
			}
		}
	}

	// Pull the image if not exist
	if !imgPulled {
		fmt.Println("Pulling the image.")
		out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("Failed to pull %s image.", imageName))
		}

		defer out.Close()

		io.Copy(os.Stdout, out)
	}

	// Check if the container already exist
	conListSummary, err := cli.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		errors.Wrap(err, "Failed to get container list.")
	}

	for _, item := range conListSummary {
		if item.Image == "miata/goj-judger-c-cpp-img:latest" {
			fmt.Println("container already created:")
			fmt.Println("container name:", item.Names[0])
			fmt.Println("container ID:", item.ID)
			conCreated = true
			// restart the container if stopped
			if item.State == "exited" {
				fmt.Println("restarting the container.")
				err = cli.ContainerStart(ctx, item.ID, types.ContainerStartOptions{})
				if err != nil {
					return errors.Wrap(err, "Failed to restart exited container.")
				}
			}
			break
		}
	}

	// Create the container if not exist
	if !conCreated {
		fmt.Println("Creating new container.")
		resp, err := cli.ContainerCreate(
			ctx,
			&container.Config{
				Image: "miata/goj-judger-c-cpp-img:latest",
				Tty:   false,
				Cmd:   strslice.StrSlice{"/bin/bash", "-c", "while true; do sleep 1; done"},
			},
			&container.HostConfig{
				Mounts: []mount.Mount{
					{
						Type:   mount.TypeBind,
						Source: "/home/lsxph/volume/c-cpp",
						Target: "/compile/c-cpp",
					},
				},
			}, nil, nil, "goj-judger-c-cpp-container")

		if err != nil {
			return errors.Wrap(err, "Error occurred when creating container.")
		}

		if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
			return errors.Wrap(err, "Failed to start container.")
		}

		fmt.Println("Container started.")
	}

	fmt.Println("Done.")

	return nil
}

func (c *c) compile() error {
	// Create new docker client instant
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	if err != nil {
		return errors.Wrap(err, "Failed to create cli instance.")
	}

	exec, err := cli.ContainerExecCreate(ctx, "goj-judger-c-cpp-container", types.ExecConfig{
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          []string{"/bin/bash", "-c", "gcc -o main main.c && echo 'pass' || echo 'fail'"},
	})

	if err != nil {
		return errors.Wrap(err, "Failed to create exec job in container.")
	}

	resp, err := cli.ContainerExecAttach(ctx, exec.ID, types.ExecStartCheck{})
	if err != nil {
		return errors.Wrap(err, "Failed to attach to exec.")
	}

	defer resp.Close()

	var outBuf, errBuf bytes.Buffer

	outputDone := make(chan error)

	go func() {
		_, err := stdcopy.StdCopy(&outBuf, &errBuf, resp.Reader)
		outputDone <- err
	}()

	err = <-outputDone
	if err != nil {
		return errors.Wrap(err, "Failed to read output from exec.")
	}

	res, err := ioutil.ReadAll(&outBuf)
	if err != nil {
		return errors.Wrap(err, "Failed to read from outBuf.")
	}

	if string(res) == "fail" {
		errors.Wrap(err, "Failed to compile source code.")
	}

	return nil
}

func (c *c) run() error {
	// Create new docker client instant
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	if err != nil {
		return errors.Wrap(err, "Failed to create cli instance.")
	}

	exec, err := cli.ContainerExecCreate(ctx, "goj-judger-c-cpp-container", types.ExecConfig{
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          []string{"/bin/bash", "-c", "./main < input"},
	})

	if err != nil {
		return errors.Wrap(err, "Failed to create exec job in container.")
	}

	// err = cli.ContainerExecStart(ctx, exec.ID, types.ExecStartCheck{})
	// if err != nil {
	// 	return errors.Wrap(err, "Failed to run exec in container.")
	// }

	resp, err := cli.ContainerExecAttach(ctx, exec.ID, types.ExecStartCheck{})
	if err != nil {
		return errors.Wrap(err, "Failed to attach to exec.")
	}

	defer resp.Close()

	var outBuf, errBuf bytes.Buffer

	outputDone := make(chan error)

	go func() {
		_, err := stdcopy.StdCopy(&outBuf, &errBuf, resp.Reader)
		outputDone <- err
	}()

	err = <-outputDone
	if err != nil {
		return errors.Wrap(err, "Failed to read output from exec.")
	}

	res, err := ioutil.ReadAll(&outBuf)
	if err != nil {
		return errors.Wrap(err, "Failed to read from outBuf.")
	}

	if compare(res) {
		fmt.Println("PASSED!")
	} else {
		fmt.Println("WRONG ANSWER.")
		fmt.Println(string(res))
		// Signal the back end
		os.Exit(1)
	}

	return nil
}
