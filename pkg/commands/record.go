package commands

import (
	"fmt"
	"github.com/hikhvar/decoherence/pkg/checker"
	"github.com/hikhvar/decoherence/pkg/prospector"
	"github.com/hikhvar/decoherence/pkg/store"
	"github.com/pkg/errors"
	"gopkg.in/urfave/cli.v1"
	"os"
	"runtime"
)

const (
	pathFlagName        = "path"
	storeFlagName       = "store"
	numParallelFlagName = "parallel"
)

func NewRecordCommand() cli.Command {
	return cli.Command{
		Name:  "record",
		Usage: "record the file tree into a store file for later compare",
		Action: func(c *cli.Context) error {
			ex := []prospector.Extractor{&checker.Checksum{}, &checker.FileInfo{}}
			s, err := store.NewWriteJSON(c.String(storeFlagName), store.Meta{
				Version:  1,
				Checkers: enabledCheckers(ex),
			})
			defer ensureErrorHandled(s)
			if err != nil {
				return cli.NewExitError(errors.Wrap(err, "Failed to create new JSON store"), 1)
			}

			p := prospector.NewProspector(c.String(pathFlagName), c.Int(numParallelFlagName), s, ex)
			err = p.Run()
			if err != nil {
				return cli.NewExitError(err, 1)
			}
			return nil
		},
		Flags: []cli.Flag{
			// TODO: Env Variables and better help
			cli.StringFlag{
				Name:  pathFlagName,
				Value: mustWorkingDirectory(),
			},
			cli.StringFlag{
				Name:  storeFlagName,
				Value: "result.json",
				Usage: "The target file",
			},
			cli.IntFlag{
				Name:  numParallelFlagName,
				Value: runtime.NumCPU(),
				Usage: "number of parallel file information extractors. Set to > 1 for maximal read performance on non HDD devices",
			},
		},
	}
}

type Closer interface {
	Close() error
}

func enabledCheckers(extractors []prospector.Extractor) []string {
	var ret []string
	for _, ex := range extractors {
		ret = append(ret, ex.Name())
	}
	return ret
}

func ensureErrorHandled(c Closer) {
	err := c.Close()
	if err != nil {
		fmt.Println(err)
	}
}

func mustWorkingDirectory() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(errors.Wrap(err, "failed to determine working directory"))
	}
	return wd
}
