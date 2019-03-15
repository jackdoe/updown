package main

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/minio/sio"
	"github.com/miquella/ask"
)

func exit(err error) {
	fmt.Fprintf(os.Stderr, "%s\n", err.Error())
	os.Exit(1)
}

func main() {
	var pkeyFrom = flag.String("k", "-", "encryption key, filename or '-' to ask; key is sha256(value)")
	var pdecrypt = flag.Bool("d", false, "decrypt")
	flag.Parse()

	h := sha256.New()
	if *pkeyFrom == "-" {
		key, err := ask.HiddenAsk("Key: ")
		if err != nil {
			exit(err)
		}
		h.Write([]byte(key))
	} else {
		f, err := os.Open(*pkeyFrom)
		if err != nil {
			exit(err)
		}
		io.Copy(h, f)
	}

	shakey := h.Sum(nil)

	if *pdecrypt {
		reader, err := sio.DecryptReader(os.Stdin, sio.Config{Key: shakey})
		if err != nil {
			exit(err)
		}
		io.Copy(os.Stdout, reader)
	} else {
		reader, err := sio.EncryptReader(os.Stdin, sio.Config{Key: shakey})
		if err != nil {
			exit(err)
		}
		io.Copy(os.Stdout, reader)
	}
}
