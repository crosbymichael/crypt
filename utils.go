package main

import (
	"crypto/md5"

	"github.com/codegangsta/cli"
)

func getKey(context *cli.Context) string {
	return context.String("key")
}

func hashKey(key []byte) []byte {
	h := md5.New()
	if _, err := h.Write(key); err != nil {
		panic(err)
	}
	return h.Sum(nil)
}

func handler(wrap func(string, string, []byte) error) func(*cli.Context) {
	return func(context *cli.Context) {
		if len(context.Args()) != 2 {
			logger.Fatal("invalid number of arguments: <file in> <file out>")
		}
		key := getKey(context)
		if key == "" {
			logger.Fatal("no key provided")
		}
		if err := wrap(context.Args().Get(0), context.Args().Get(1), []byte(key)); err != nil {
			logger.Fatal(err)
		}
	}
}
