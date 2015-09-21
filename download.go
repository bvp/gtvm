package main

import (
	"fmt"
	"github.com/cheggaaa/pb"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

func download(gt, url, outDir string) {
	sourceName, destName := url, archivesDir+ps+outDir
	if _, err := os.Stat(destName); err == nil {
		fmt.Println("Already downloaded...")
		if gt == "golang" {
			if compareHash("1.5", checksum(destName)) {
				return
			}
		}
		return
	}
	var source io.Reader
	var sourceSize int64
	resp, err := http.Get(sourceName)
	if err != nil {
		fmt.Printf("Can't get %s: %v\n", sourceName, err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Server return non-200 status: %v\n", resp.Status)
		return
	}
	i, _ := strconv.Atoi(resp.Header.Get("Content-Length"))
	sourceSize = int64(i)
	source = resp.Body

	dest, err := os.Create(destName)
	if err != nil {
		fmt.Printf("Can't create %s: %v\n", destName, err)
		return
	}
	defer dest.Close()

	bar := pb.New(int(sourceSize)).SetUnits(pb.U_BYTES).SetRefreshRate(time.Millisecond * 10)
	bar.ShowSpeed = true
	bar.Start()
	writer := io.MultiWriter(dest, bar)
	io.Copy(writer, source)
	bar.Finish()
}
