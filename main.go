package main

import (
	"fmt"
	"github.com/hikhvar/decoherence/pkg/commands"
	"log"
	"os"

	"gopkg.in/urfave/cli.v1"
)

func main() {
	app := cli.NewApp()
	app.Name = "decoherence"
	app.Usage = "checks the difference between two directory trees"
	app.Commands = []cli.Command{
		commands.NewRecordCommand(),
		commands.NewCompareCommand(),
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Run completed")
	if err != nil {
		fmt.Println("failed to fetch data", err)
	}
}
