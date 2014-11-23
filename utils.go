package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"fmt"
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
	var iv [aes.BlockSize]byte
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	r, w := h(inf, outf, cipher.NewOFB(block, iv[:]))
	bar, err := newProgressBar(inf)
	if err != nil {
		return err
	}
	_, err = io.Copy(io.MultiWriter(w, bar), r)
	return err
}

func newProgressBar(in *os.File) (io.Writer, error) {
	stat, err := in.Stat()
	if err != nil {
		return nil, err
	}
	bar := pb.New(int(stat.Size())).SetUnits(pb.U_BYTES)
	bar.ShowSpeed, bar.ShowTimeLeft = true, false
	bar.Start()
	return bar, nil
}

func getKey(context *cli.Context) []byte {
	key := context.GlobalString("key")
	if key == "" {
		fmt.Fprint(os.Stdout, "please enter your key:\n> ")
		s := bufio.NewScanner(os.Stdin)
		s.Scan()
		key = s.Text()
		if key == "" {
			return nil
		}
	}
	h := md5.New()
	if _, err := fmt.Fprint(h, key); err != nil {
		panic(err)
	}
	return h.Sum(nil)
}
