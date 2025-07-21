package main

import (
	"context"
	"fmt"
	"os"

	concurrently "github.com/benfiola/go-concurrently/pkg"
	"github.com/google/shlex"
	"github.com/urfave/cli/v3"
)

func main() {
	version := concurrently.GetVersion()

	cli.VersionPrinter = func(cmd *cli.Command) {
		fmt.Fprintf(cmd.Root().Writer, "%s\n", cmd.Root().Version)
	}

	err := (&cli.Command{
		Action: func(ctx context.Context, c *cli.Command) error {
			cmdSlices := [][]string{}
			for _, cmdStr := range c.StringArgs("commands") {
				cmdSlice, err := shlex.Split(cmdStr)
				if err != nil {
					return err
				}
				cmdSlices = append(cmdSlices, cmdSlice)
			}
			return concurrently.Run(ctx, cmdSlices...)
		},
		Arguments: []cli.Argument{
			&cli.StringArgs{
				Name: "commands",
				Min:  1,
				Max:  -1,
			},
		},
		Description: "run commands concurrently",
		Version:     version,
	}).Run(context.Background(), os.Args)

	code := 0
	if err != nil {
		code = 1
	}
	os.Exit(code)
}
