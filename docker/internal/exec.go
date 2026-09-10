package internal

import (
	"context"
	"fmt"
	"io"
	"time"

	ctr "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

type ExecOptions struct {
	Cmd    []string
	Stdout io.Writer
	Stderr io.Writer
}

func Exec(
	ctx context.Context,
	client client.ContainerAPIClient,
	container string,
	options ExecOptions,
) error {
	id, err := client.ContainerExecCreate(ctx, container, ctr.ExecOptions{
		Cmd:          options.Cmd,
		AttachStdout: options.Stdout != nil,
		AttachStderr: options.Stderr != nil,
	})
	if err != nil {
		return fmt.Errorf("creating exec: %w", err)
	}

	conn, err := client.ContainerExecAttach(ctx, id.ID, ctr.ExecStartOptions{})
	if err != nil {
		return fmt.Errorf("attaching to exec process: %w", err)
	}
	defer conn.Close()

	err = client.ContainerExecStart(ctx, id.ID, ctr.ExecStartOptions{})
	if err != nil {
		return fmt.Errorf("starting exec: %w", err)
	}

	var stat ctr.ExecInspect
	ctx, cancel := context.WithTimeout(ctx, time.Second*15)
	defer cancel()

	for {
		stat, err = client.ContainerExecInspect(ctx, id.ID)
		if err != nil {
			return err
		}
		if !stat.Running {
			break
		}
	}
	if stat.ExitCode != 0 {
		return fmt.Errorf("exec returned non-zero exit code: %d", stat.ExitCode)
	}

	if options.Stdout != nil {
		_, err = stdcopy.StdCopy(options.Stdout, options.Stderr, conn.Reader)
		if err != nil {
			return fmt.Errorf("copying exec output: %w", err)
		}
	}

	return nil
}
