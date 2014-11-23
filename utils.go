package main

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"os"

	"github.com/codegangsta/cli"
)

// getKey returns a key that is 32 bytes in size from the
// cli context or from STDIN
func getKey(context *cli.Context) []byte {
	k := context.GlobalString("key")
	if k == "" {
		fmt.Fprint(os.Stdout, "please enter your key:\n> ")
		s := bufio.NewScanner(os.Stdin)
		s.Scan()
		k = s.Text()
		if k == "" {
			return nil
		}
	}
	h := md5.New()
	if _, err := fmt.Fprint(h, k); err != nil {
		panic(err)
	}
	return h.Sum(nil)
}
