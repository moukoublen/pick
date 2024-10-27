//go:build install_shellcheck

package main

import (
	"fmt"
	"math/rand/v2"
	"os"
	"time"
)

func randString(ln int) string {
	const randStringVocab = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	b := make([]byte, ln)
	for i := range ln {
		b[i] = randStringVocab[rand.IntN(ln)]
	}
	return string(b)
}

func tmpDir(name string) (string, func(), error) {
	os.TempDir()
	fullPath := fmt.Sprintf("%s/%s_%s", os.TempDir(), name, randString(6))

	err := os.Mkdir(fullPath, 0o755) // drwxr-xr-x
	if err != nil {
		return "", func() {}, err
	}

	return fullPath, func() {
		_ = os.RemoveAll(fullPath) // rm -rf
	}, nil
}

func main() {
	tmpFullPath, cleanUp, err := tmpDir("shellcheck")
	defer cleanUp()

	time.Sleep(4 * time.Second)
	fmt.Printf(">>> %s %#v\n\n", tmpFullPath, err)
}
