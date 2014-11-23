package main

import (
	"crypto/cipher"
	"io"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

var logger = logrus.New()

func encrypt(in, out *os.File, stream cipher.Stream) (io.Reader, io.Writer) {
	return in, &cipher.StreamWriter{S: stream, W: out}
}

func decrypt(in, out *os.File, stream cipher.Stream) (io.Reader, io.Writer) {
	return &cipher.StreamReader{S: stream, R: in}, out
}

func main() {
	app := cli.NewApp()
	app.Name = "crypt"
	app.Version = "1"
	app.Author = "@crosbymichael"
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "key", Usage: "key to use for the encryption algo"},
	}
	app.Commands = []cli.Command{
		cli.Command{
			Name:   "encrypt",
			Usage:  "encript a file",
			Action: handler(encrypt),
		},
		cli.Command{
			Name:   "decrypt",
			Usage:  "decrypt a file",
			Action: handler(decrypt),
		},
	}

	if err := app.Run(os.Args); err != nil {
		logger.Fatal(err)
	}
}
