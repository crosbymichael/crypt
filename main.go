package main

import (
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

var (
	key    []byte
	logger = logrus.New()
)

func inputFile(context *cli.Context) (*os.File, error) {
	if context.GlobalBool("stdin") {
		return os.Stdin, nil
	}
	in, err := os.Open(context.Args().Get(0))
	if err != nil {
		return nil, err
	}
	return in, nil
}

func stdoutArgIndex(context *cli.Context) int {
	if !context.GlobalBool("stdin") {
		return 1
	}
	return 0
}

func outputFile(context *cli.Context) (*os.File, error) {
	if context.GlobalBool("stdout") {
		return os.Stdout, nil
	}
	out, err := os.Create(context.Args().Get(stdoutArgIndex(context)))
	if err != nil {
		return nil, err
	}
	return out, nil
}

func do(context *cli.Context, key []byte, a Action) error {
	in, err := inputFile(context)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := outputFile(context)
	if err != nil {
		return err
	}
	defer out.Close()
	p, err := newProcessor(in, out, key, a)
	if err != nil {
		return err
	}
	return p.Run()
}

func getAction(context *cli.Context) Action {
	switch {
	case context.GlobalBool("encrypt"):
		return Encrypt
	case context.GlobalBool("decrypt"):
		return Decrypt
	}
	return 0
}

func main() {
	app := cli.NewApp()
	app.Name = "crypt"
	app.Version = "alpha"
	app.Author = "@crosbymichael"
	app.Usage = `
encrypt and decrypt files easily

NOTE!: While the version is alpha things may break between commits.  
Do not expect compatibility between builds until the version goes to 1.
`

	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "key", Usage: "key to use for the encryption algo"},
		cli.BoolFlag{Name: "encrypt,e", Usage: "encrypt a file"},
		cli.BoolFlag{Name: "decrypt,d", Usage: "decrypt a file"},
		cli.BoolFlag{Name: "stdin,i", Usage: "accept input for STDIN"},
		cli.BoolFlag{Name: "stdout,o", Usage: "return output to STDOUT"},
	}
	app.Before = func(context *cli.Context) error {
		if !context.GlobalBool("encrypt") && !context.GlobalBool("decrypt") {
			return nil
		}
		app.Action = func(context *cli.Context) {
			a := getAction(context)
			if err := do(context, key, a); err != nil {
				logger.Fatal(err)
			}
		}
		if context.GlobalBool("stdin") && context.GlobalString("key") == "" {
			return fmt.Errorf("--key must be supplied when receiving input via STDIN")
		}
		key = getKey(context)
		if len(key) == 0 {
			return fmt.Errorf("no key provided via --key or STDIN")
		}
		return nil
	}
	if err := app.Run(os.Args); err != nil {
		logger.Fatal(err)
	}
}
