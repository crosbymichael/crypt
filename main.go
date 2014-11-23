package main

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

var logger = logrus.New()

func main() {
	app := cli.NewApp()
	app.Name = "crypt"
	app.Version = "1"
	app.Author = "@crosbymichael"
	app.Commands = []cli.Command{
		encryptCommand,
		decryptCommand,
	}

	if err := app.Run(os.Args); err != nil {
		logger.Fatal(err)
	}
}
