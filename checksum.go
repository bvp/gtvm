package main

import (
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"io"
	"math"
	"os"
)

const filechunk = 8192 // we settle for 8KB

func checksum(f string) []string {
    var hashs []string
	file, err := os.Open(f)
	defer file.Close()

	if err != nil {
		panic(err.Error())
	}

	// calculate the file size
	info, _ := file.Stat()
	filesize := info.Size()
	blocks := uint64(math.Ceil(float64(filesize) / float64(filechunk)))
	hash1 := sha1.New()
	hash256 := sha256.New()

	for i := uint64(0); i < blocks; i++ {
		blocksize := int(math.Min(filechunk, float64(filesize-int64(i*filechunk))))
		buf := make([]byte, blocksize)

		file.Read(buf)
		io.WriteString(hash1, string(buf)) // append into the hash
		io.WriteString(hash256, string(buf)) // append into the hash
	}
    hashs = append(hashs, fmt.Sprintf("%x", hash1.Sum(nil)))
    hashs = append(hashs, fmt.Sprintf("%x", hash256.Sum(nil)))
	return hashs
}
