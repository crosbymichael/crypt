package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"io"
	"os"

	"github.com/codegangsta/cli"
	"github.com/rakyll/pb"
)

type rwHandler func(in, out *os.File, stream cipher.Stream) (io.Reader, io.Writer)

func process(in, out string, key []byte, h rwHandler) error {
	outf, err := os.Create(out)
	if err != nil {
		return err
	}
	defer outf.Close()

	inf, err := os.Open(in)
	if err != nil {
		return err
	}
	defer inf.Close()

	bar, err := newProgressBar(inf)
	if err != nil {
		return err
	}

	var iv [aes.BlockSize]byte
	block, err := newBlock(key)
	if err != nil {
		return err
	}
	stream := cipher.NewOFB(block, iv[:])
	r, w := h(inf, outf, stream)
	mw := io.MultiWriter(w, bar)
	if _, err := io.Copy(mw, r); err != nil {
		return err
	}
	return nil
}

func newProgressBar(in *os.File) (io.Writer, error) {
	stat, err := in.Stat()
	if err != nil {
		return nil, err
	}
	bar := pb.New(int(stat.Size())).SetUnits(pb.U_BYTES)
	bar.ShowSpeed = true
	bar.ShowTimeLeft = false
	bar.Start()
	return bar, nil
}

func newBlock(key []byte) (cipher.Block, error) {
	return aes.NewCipher(hashKey(key))
}

func getKey(context *cli.Context) string {
	return context.GlobalString("key")
}

func hashKey(key []byte) []byte {
	h := md5.New()
	if _, err := h.Write(key); err != nil {
		panic(err)
	}
	return h.Sum(nil)
}

func handler(wrap rwHandler) func(*cli.Context) {
	return func(context *cli.Context) {
		if len(context.Args()) != 2 {
			logger.Fatal("invalid number of arguments: <file in> <file out>")
		}
		key := getKey(context)
		if key == "" {
			logger.Fatal("no key provided")
		}
		if err := process(context.Args().Get(0), context.Args().Get(1), []byte(key), wrap); err != nil {
			logger.Fatal(err)
		}
	}
}
