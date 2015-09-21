package main

import (
	"crypto/sha1"
	"fmt"
	"io"
	"math"
	"os"
)

const filechunk = 8192 // we settle for 8KB

func checksum(f string) string {

	file, err := os.Open(f)

	if err != nil {
		panic(err.Error())
	}

	defer file.Close()

	// calculate the file size
	info, _ := file.Stat()
	filesize := info.Size()
	blocks := uint64(math.Ceil(float64(filesize) / float64(filechunk)))
	hash := sha1.New()

	for i := uint64(0); i < blocks; i++ {
		blocksize := int(math.Min(filechunk, float64(filesize-int64(i*filechunk))))
		buf := make([]byte, blocksize)

		file.Read(buf)
		io.WriteString(hash, string(buf)) // append into the hash
	}

	//	fmt.Printf("%s checksum is %x - len: %s\n", file.Name(), hash.Sum(nil), fmt.Sprintf("%x", hash.Sum(nil)))
	return fmt.Sprintf("%x", hash.Sum(nil))
}
