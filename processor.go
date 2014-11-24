package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"os"

	"github.com/rakyll/pb"
)

type Action int

const (
	Encrypt Action = iota
	Decrypt
)

func newProcessor(in, out *os.File, key []byte, a Action) (*processor, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	stat, err := in.Stat()
	if err != nil {
		return nil, err
	}
	return &processor{
		In:     in,
		Out:    out,
		action: a,
		size:   int(stat.Size()),
		block:  block,
	}, nil
}

type processor struct {
	action Action
	In     *os.File
	Out    *os.File
	size   int
	block  cipher.Block
}

func (p *processor) Run() error {
	iv, err := p.getIV()
	if err != nil {
		return err
	}
	r, w := p.io(p.stream(iv))
	if p.size > 0 {
		w = io.MultiWriter(w, p.newProgressBar())
	}
	_, err = io.Copy(w, r)
	return err
}

func (p *processor) randomIV(iv []byte) error {
	_, err := rand.Read(iv)
	return err
}

func (p *processor) getIV() ([]byte, error) {
	iv := make([]byte, aes.BlockSize)
	switch p.action {
	case Encrypt:
		if err := p.randomIV(iv); err != nil {
			return nil, err
		}
		// write the iv as the first bytes of the output file
		wrote, err := p.Out.Write(iv)
		if err != nil {
			return nil, err
		}
		if wrote != aes.BlockSize {
			return nil, fmt.Errorf("unable to write correct iv length %d != %d", wrote, aes.BlockSize)
		}
	case Decrypt:
		// read the previous iv as the first bytes of the input file
		read, err := p.In.Read(iv)
		if err != nil {
			return nil, err
		}
		if read != aes.BlockSize {
			return nil, fmt.Errorf("did not read correct iv amount %d != %d", read, aes.BlockSize)
		}
	}
	return iv, nil
}

func (p *processor) stream(iv []byte) cipher.Stream {
	switch p.action {
	case Encrypt:
		return cipher.NewCFBEncrypter(p.block, iv)
	case Decrypt:
		return cipher.NewCFBDecrypter(p.block, iv)
	}
	return nil
}

func (p *processor) io(s cipher.Stream) (r io.Reader, w io.Writer) {
	switch p.action {
	case Encrypt:
		r = p.In
		w = &cipher.StreamWriter{S: s, W: p.Out}
	case Decrypt:
		w = p.Out
		r = &cipher.StreamReader{S: s, R: p.In}
	}
	return
}

func (p *processor) newProgressBar() io.Writer {
	bar := pb.New(p.size).SetUnits(pb.U_BYTES)
	bar.ShowSpeed, bar.ShowTimeLeft = true, false
	bar.Start()
	return bar
}
