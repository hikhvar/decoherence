package commands

import (
	"fmt"
	"github.com/hikhvar/decoherence/pkg/store"
	"github.com/pkg/errors"
	"gopkg.in/urfave/cli.v1"
)

const (
	expectedStoreFlagName = "expected"
	gotStoreFlagName      = "got"
)

func NewCompareCommand() cli.Command {
	return cli.Command{
		Name:  "check",
		Usage: "compares two records",
		Action: func(c *cli.Context) error {
			expected, err := store.NewReadJSON(c.String(expectedStoreFlagName))
			if err != nil {
				return cli.NewExitError(errors.Wrap(err, "failed to read expected store file"), 1)
			}
			got, err := store.NewReadJSON(c.String(gotStoreFlagName))
			if err != nil {
				return cli.NewExitError(errors.Wrap(err, "failed to read got store file"), 1)
			}
			// TODO: Check metadata
			res := store.ComputeDiffs(expected.Files(), got.Files())
			fmt.Printf("%+v\n", res)
			return nil
		},
		Flags: []cli.Flag{
			// TODO: Env Variables and better help
			cli.StringFlag{
				Name: expectedStoreFlagName,
			},
			cli.StringFlag{
				Name: gotStoreFlagName,
			},
		},
	}
}
