package main

import (
	"crypto/aes"
	"crypto/cipher"
	"io"
	"os"

	"github.com/codegangsta/cli"
	"github.com/rakyll/pb"
)

var decryptCommand = cli.Command{
	Name:  "decrypt",
	Usage: "decrypt a file",
	Flags: []cli.Flag{
		cli.StringFlag{Name: "key", Usage: "key to use for the encryption algo"},
	},
	Action: handler(decrypt),
}

func decrypt(in, out string, key []byte) error {
	block, err := aes.NewCipher(hashKey(key))
	if err != nil {
		return err
	}
	f, err := os.Open(in)
	if err != nil {
		return err
	}
	defer f.Close()
	stat, err := f.Stat()
	if err != nil {
		return err
	}

	var iv [aes.BlockSize]byte
	r := &cipher.StreamReader{S: cipher.NewOFB(block, iv[:]), R: f}

	w, err := os.Create(out)
	if err != nil {
		return err
	}
	defer w.Close()
	bar := pb.New(int(stat.Size())).SetUnits(pb.U_BYTES)
	bar.ShowSpeed = true
	bar.ShowTimeLeft = false
	bar.Start()
	mw := io.MultiWriter(w, bar)

	if _, err := io.Copy(mw, r); err != nil {
		return err
	}
	return nil
}
