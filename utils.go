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
	block, err := newBlock(key)
	if err != nil {
		return err
	}
	stream := cipher.NewOFB(block, iv[:])
	r, w := h(inf, outf, stream)
	bar, err := newProgressBar(inf)
	if err != nil {
		return err
	}
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
	return aes.NewCipher(key)
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
	return hashKey(key)
}

// hashKey hashes the provided key using md5 to ensure that it is
// 32 bytes long for used with the encryption algos
func hashKey(key string) []byte {
	h := md5.New()
	if _, err := fmt.Fprint(h, key); err != nil {
		panic(err)
	}
	return h.Sum(nil)
}
