package judger

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
)

var caseNum int
var tlimit int
var rlimit int
var volumeDir string
var buildCmd string
var runCmd string
var imageName string
var imgVersion string = ":latest"
var containerName string

var ctx context.Context

func init() {
	ctx = context.Background()
}

func prepareImg() error {
	fmt.Println("正在准备编译容器...")

	imgPulled := false
	conCreated := false

	// 创建Docker 客户端实例
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	defer cli.Close()

	if err != nil {
		return errors.Wrap(err, "创建Docker客户端失败")
	}

	// 检查镜像是否存在
	imgListSummary, err := cli.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		return errors.Wrap(err, "获取镜像列表失败")
	}

	for _, item := range imgListSummary {
		if len(item.RepoTags) > 0 {
			for _, tag := range item.RepoTags {
				if tag == imageName+imgVersion {
					fmt.Println("镜像已存在：", tag)
					imgPulled = true
					break
				}
			}
		}
	}

	// Pull the image if not exist
	if !imgPulled {
		fmt.Println("正在拉取判题镜像")

		authConfig := types.AuthConfig{
			Username: "miata",
			Password: "Docker75218644",
		}

		encodedJson, err := json.Marshal(authConfig)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("拉取镜像 %s 失败", imageName+imgVersion))
		}

		authStr := base64.URLEncoding.EncodeToString(encodedJson)

		out, err := cli.ImagePull(ctx, imageName+imgVersion, types.ImagePullOptions{RegistryAuth: authStr})
		defer out.Close()

		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("拉取镜像 %s 失败", imageName+imgVersion))
		}

		io.Copy(os.Stdout, out)
	}

	// 检查判题容器是否存在
	conListSummary, err := cli.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		errors.Wrap(err, "获取容器列表失败")
	}

	for _, item := range conListSummary {
		if item.Image == imageName+imgVersion {
			fmt.Println("容器已经创建")
			fmt.Println("容器名称：", item.Names[0])
			fmt.Println("容器ID：", item.ID)
			conCreated = true
			// restart the container if stopped
			if item.State == "exited" {
				fmt.Println("正在重启容器...")
				err = cli.ContainerStart(ctx, item.ID, types.ContainerStartOptions{})
				if err != nil {
					return errors.Wrap(err, "重启容器失败")
				}
			}
			break
		}
	}

	// Create the container if not exist
	if !conCreated {
		fmt.Println("正在创建新容器...")
		PIDLimit := new(int64)
		*PIDLimit = 100
		resp, err := cli.ContainerCreate(
			ctx,
			&container.Config{
				Image: imageName + imgVersion,
				Tty:   false,
				Cmd:   strslice.StrSlice{"/bin/bash", "-c", "while true; do sleep 1; done"},
			},
			&container.HostConfig{
				Mounts: []mount.Mount{
					{
						Type:   mount.TypeBind,
						Source: volumeDir,
						Target: "/judge",
					},
				},
				AutoRemove:  true,   // 容器退出后自动删除
				NetworkMode: "none", // 禁止网络
				Resources: container.Resources{
					// Ulimits: []*units.Ulimit{{
					// 	Name: "fsize",
					// 	Hard:
					// }},
					Memory:    1024 * 1000 * 1024, // 限制容器内存使用为1024MB
					PidsLimit: PIDLimit,           // 限制进程个数
				},
			}, nil, nil, containerName)

		if err != nil {
			return errors.Wrap(err, "创建容器时发生错误")
		}

		if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
			return errors.Wrap(err, "启动容器失败")
		}

		fmt.Println("成功启动判题容器")
	}

	fmt.Println("准备阶段结束")
	return nil
}

func compile() error {
	fmt.Println("正在编译代码...")

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	defer cli.Close()

	if err != nil {
		return errors.Wrap(err, "创建Docker客户端失败")
	}

	fmt.Println("buildCmd:", buildCmd)
	exec, err := cli.ContainerExecCreate(ctx, containerName, types.ExecConfig{
		AttachStdout: true,
		AttachStderr: true,
		WorkingDir:   "/judge",
		Cmd:          []string{"/bin/bash", "-c", buildCmd},
	})

	if err != nil {
		return errors.Wrap(err, "在Docker容器内创建exec失败")
	}

	err = cli.ContainerExecStart(ctx, exec.ID, types.ExecStartCheck{
		Detach: true,
	})

	if err != nil {
		return errors.Wrap(err, "在Docker容器内执行exec失败")
	}

	inspect, err := cli.ContainerExecInspect(ctx, exec.ID)
	if err != nil {
		return errors.Wrap(err, "获取exec结果失败")
	}

	for inspect.Running {
		inspect, _ = cli.ContainerExecInspect(ctx, exec.ID)
	}

	fmt.Println("ExitCode:", inspect.ExitCode)

	// 编译错误
	if inspect.ExitCode != 0 {
		fmt.Println("编译失败", "ExitCode: ", inspect.ExitCode)
		os.Exit(COMPILE_ERROR)
	}

	return nil
}

func run() error {
	fmt.Println("正在运行代码")
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	defer cli.Close()

	if err != nil {
		return errors.Wrap(err, "创建Docker客户端失败")
	}

	fmt.Println("runCmd:", runCmd)
	exec, err := cli.ContainerExecCreate(ctx, containerName, types.ExecConfig{
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          []string{"/bin/bash", "-c", runCmd},
	})

	if err != nil {
		return errors.Wrap(err, "在Docker容器创建exec失败")
	}

	err = cli.ContainerExecStart(ctx, exec.ID, types.ExecStartCheck{
		Detach: true,
	})

	if err != nil {
		return errors.Wrap(err, "在Docker容器运行exec失败")
	}

	inspect, err := cli.ContainerExecInspect(ctx, exec.ID)

	if err != nil {
		return errors.Wrap(err, "获取exec结果失败")
	}

	for inspect.Running {
		inspect, _ = cli.ContainerExecInspect(ctx, exec.ID)
	}

	if inspect.ExitCode != 0 {
		fmt.Println("RUNTIME_ERROR")
		os.Exit(RUNTIME_ERROR)
	}

	fmt.Println("运行结束")

	return nil
}
