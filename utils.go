package main

import "crypto/md5"

func hashKey(key []byte) []byte {
	h := md5.New()
	if _, err := h.Write(key); err != nil {
		panic(err)
	}
	return h.Sum(nil)
}
