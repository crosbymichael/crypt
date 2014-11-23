package main

import (
	"crypto/aes"
	"crypto/cipher"
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
	r, w := p.io(p.stream(p.newIV()))
	if p.size > 0 {
		w = io.MultiWriter(w, p.newProgressBar())
	}
	_, err := io.Copy(w, r)
	return err
}

func (p *processor) newIV() []byte {
	var iv [aes.BlockSize]byte
	return iv[:]
}

func (p *processor) stream(iv []byte) cipher.Stream {
	return cipher.NewOFB(p.block, iv)
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
