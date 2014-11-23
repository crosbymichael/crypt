package main

import (
	"crypto/cipher"
	"fmt"
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

func getInputFile(context *cli.Context) *os.File {
	if context.GlobalBool("stdin") {
		return os.Stdin
	}
	in, err := os.Open(context.Args().Get(0))
	if err != nil {
		logger.Fatal(err)
	}
	return in
}

func stdoutArgIndex(context *cli.Context) int {
	if !context.GlobalBool("stdin") {
		return 1
	}
	return 0
}

func getOutputFile(context *cli.Context) *os.File {
	if context.GlobalBool("stdout") {
		return os.Stdout
	}
	out, err := os.Create(context.Args().Get(stdoutArgIndex(context)))
	if err != nil {
		logger.Fatal(err)
	}
	return out
}

func action(wrap rwHandler) func(*cli.Context) {
	return func(context *cli.Context) {
		key := getKey(context)
		if len(key) == 0 {
			logger.Fatal("no key provided")
		}
		in, out := getInputFile(context), getOutputFile(context)
		err := process(in, out, key, wrap)
		in.Close()
		out.Close()
		if err != nil {
			if outPath := context.Args().Get(stdoutArgIndex(context)); outPath != "" {
				os.Remove(outPath)
			}
			logger.Fatal(err)
		}
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "crypt"
	app.Version = "1"
	app.Author = "@crosbymichael"
	app.Usage = "encrypt and decrypt files easily"
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "key", Usage: "key to use for the encryption algo"},
		cli.BoolFlag{Name: "encrypt,e", Usage: "encrypt a file"},
		cli.BoolFlag{Name: "decrypt,d", Usage: "decrypt a file"},
		cli.BoolFlag{Name: "stdin,i", Usage: "accept input for STDIN"},
		cli.BoolFlag{Name: "stdout,o", Usage: "return output to STDOUT"},
	}
	app.Before = func(context *cli.Context) error {
		if context.GlobalBool("stdin") && context.GlobalString("key") == "" {
			return fmt.Errorf("--key must be supplied when receiving input via STDIN")
		}
		switch {
		case context.GlobalBool("encrypt"):
			app.Action = action(encrypt)
		case context.GlobalBool("decrypt"):
			app.Action = action(decrypt)
		}
		return nil
	}
	if err := app.Run(os.Args); err != nil {
		logger.Fatal(err)
	}
}
