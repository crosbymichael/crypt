package main

import (
	"crypto/aes"
	"crypto/cipher"
	"io"
	"os"

	"github.com/codegangsta/cli"
	"github.com/rakyll/pb"
)

var encryptCommand = cli.Command{
	Name:  "encrypt",
	Usage: "encript a file",
	Flags: []cli.Flag{
		cli.StringFlag{Name: "key", Usage: "key to use for the encryption algo"},
	},
	Action: encryptAction,
}

func getKey(context *cli.Context) string {
	return context.String("key")
}

func encryptFile(in, out string, key []byte) error {
	block, err := aes.NewCipher(hashKey(key))
	if err != nil {
		return err
	}
	of, err := os.Create(out)
	if err != nil {
		return err
	}
	defer of.Close()
	var iv [aes.BlockSize]byte
	w := &cipher.StreamWriter{S: cipher.NewOFB(block, iv[:]), W: of}
	f, err := os.Open(in)
	if err != nil {
		return err
	}
	defer f.Close()
	stat, err := f.Stat()
	if err != nil {
		return err
	}
	bar := pb.New(int(stat.Size())).SetUnits(pb.U_BYTES)
	bar.ShowSpeed = true
	bar.ShowTimeLeft = false
	bar.Start()
	mw := io.MultiWriter(w, bar)

	if _, err := io.Copy(mw, f); err != nil {
		return err
	}
	return nil
}

func encryptAction(context *cli.Context) {
	if len(context.Args()) != 2 {
		logger.Fatal("invalid number of arguments: <file in> <file out>")
	}
	key := getKey(context)
	if key == "" {
		logger.Fatal("no key provided")
	}
	if err := encryptFile(context.Args().Get(0), context.Args().Get(1), []byte(key)); err != nil {
		logger.Fatal(err)
	}
}
