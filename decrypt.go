package main

import (
	"crypto/aes"
	"crypto/cipher"
	"io"
	"os"

	"github.com/codegangsta/cli"
)

var decryptCommand = cli.Command{
	Name:  "decrypt",
	Usage: "decrypt a file",
	Flags: []cli.Flag{
		cli.StringFlag{Name: "key", Usage: "key to use for the encryption algo"},
	},
	Action: decryptAction,
}

func decryptFile(in, out string, key []byte) error {
	block, err := aes.NewCipher(hashKey(key))
	if err != nil {
		return err
	}
	f, err := os.Open(in)
	if err != nil {
		return err
	}
	defer f.Close()

	var iv [aes.BlockSize]byte
	r := &cipher.StreamReader{S: cipher.NewOFB(block, iv[:]), R: f}

	w, err := os.Create(out)
	if err != nil {
		return err
	}
	defer w.Close()
	if _, err := io.Copy(w, r); err != nil {
		return err
	}
	return nil
}

func decryptAction(context *cli.Context) {
	if len(context.Args()) != 2 {
		logger.Fatal("invalid number of arguments: <file in> <file out>")
	}
	key := getKey(context)
	if key == "" {
		logger.Fatal("no key provided")
	}
	if err := decryptFile(context.Args().Get(0), context.Args().Get(1), []byte(key)); err != nil {
		logger.Fatal(err)
	}
}
