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
		cli.BoolFlag{Name: "encrypt,e", Usage: "encrypt a file"},
		cli.BoolFlag{Name: "decrypt,d", Usage: "decrypt a file"},
	}
	app.Before = func(context *cli.Context) error {
		switch {
		case context.GlobalBool("encrypt"):
			app.Action = handler(encrypt)
		case context.GlobalBool("decrypt"):
			app.Action = handler(decrypt)
		}
		return nil
	}
	if err := app.Run(os.Args); err != nil {
		logger.Fatal(err)
	}
}
