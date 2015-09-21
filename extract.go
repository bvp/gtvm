package main

import (
	"archive/tar"
	"archive/zip"
	"compress/bzip2"
	"compress/gzip"
	"fmt"
	"github.com/cheggaaa/pb"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func extract(filename, ver string) {
	ext := path.Ext(filename)
	//	fmt.Println(ext)
	if ext == ".zip" {
		//		fmt.Println("zip archive")
		unzip(archivesDir+ps+filename, gtvmDir, ver)
	} else if ext == (".bz2") || ext == (".gz") {
		//		fmt.Println("gz or bz2 archive")
		unGzipBzip2(archivesDir+ps+filename, gtvmDir, ver)
	} else if ext == ".7z" {
		fmt.Println("7z archive")
	}
}

func unzip(filename, dest, ver string) {
	var path string
	if filename == "" {
		fmt.Println("Can't unzip ", filename)
		os.Exit(1)
	}

	/* if filename[:2] == "go" {
		dest = dest + ps + strings.Replace(filename, "go", "go"+ps+ver, 1)
	} else if filename[:2] == "li" {
		dest = dest + ps + strings.Replace(filename, "liteide", "liteide"+ps+ver, 1)
	} */

	reader, err := zip.OpenReader(filename)
	checkErr("Extract error::OpenArchive", err)
	defer reader.Close()

	fl := int(len(reader.Reader.File))
	bar := pb.StartNew(fl)
	bar.ShowPercent = true
	bar.ShowCounters = false
	bar.ShowTimeLeft = false
	bar.Prefix("Extracting " + filename[strings.LastIndex(filename, ps)+1:] + " ")
	bar.Start()
	for _, f := range reader.Reader.File {
		zipped, err := f.Open()
		checkErr("Extract error::", err)
		defer zipped.Close()

		// path := filepath.Join(dest, ver, "./", f.Name)
		if f.Name[:2] == "go" {
			path = filepath.Join(dest, "./", strings.Replace(f.Name, "go", "go"+ps+ver, 1))
		} else if f.Name[:2] == "li" {
			path = filepath.Join(dest, "./", strings.Replace(f.Name, "liteide", "liteide"+ps+ver, 1))
		}
		fmt.Println(path)

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			writer, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, f.Mode())
			checkErr("Extract error::OpenFileFromArchive", err)
			defer writer.Close()

			if _, err = io.Copy(writer, zipped); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
		//  progress = (i / fl) * 100
		bar.Increment()
	}
	bar.Finish()
}

func tarFilesCount(sourcefile string) int {
	flreader, _ := os.Open(sourcefile)
	defer flreader.Close()
	var fltarReader *tar.Reader
	var flReader io.ReadCloser = flreader

	if strings.HasSuffix(sourcefile, ".gz") ||
		strings.HasSuffix(sourcefile, ".tgz") {
		flgzipReader, err := gzip.NewReader(flreader)
		checkErr("In tarFilesCounter - NewReader", err)
		fltarReader = tar.NewReader(flgzipReader)
		defer flReader.Close()
	} else if strings.HasSuffix(sourcefile, ".bz2") {
		flbz2Reader := bzip2.NewReader(flreader)
		fltarReader = tar.NewReader(flbz2Reader)
	} else {
		fltarReader = tar.NewReader(flreader)
	}

	trfl := fltarReader
	counter := 0
	for {
		_, err := trfl.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			checkErr("Extract error::ReadTarArchive", err)
		}
		counter++
	}
	fmt.Println("Files in archive -", counter)
	return counter
}

func unGzipBzip2(sourcefile, dest, ver string) {
	reader, err := os.Open(sourcefile)
	checkErr("In unGzipBzip2 - Open", err)
	defer reader.Close()

	var tarReader *tar.Reader
	var fileReader io.ReadCloser = reader

	if strings.HasSuffix(sourcefile, ".gz") ||
		strings.HasSuffix(sourcefile, ".tgz") {
		gzipReader, err := gzip.NewReader(reader)
		checkErr("In unGzipBzip2 - NewReader", err)
		tarReader = tar.NewReader(gzipReader)
		defer fileReader.Close()
	} else if strings.HasSuffix(sourcefile, ".bz2") {
		bz2Reader := bzip2.NewReader(reader)
		tarReader = tar.NewReader(bz2Reader)
	} else {
		tarReader = tar.NewReader(reader)
	}

	bar := pb.StartNew(tarFilesCount(sourcefile))
	bar.ShowPercent = true
	bar.ShowCounters = false
	bar.ShowTimeLeft = false
	bar.Prefix("Extracting " + sourcefile[strings.LastIndex(sourcefile, ps)+1:] + " ")
	bar.Start()
	for {
		header, err := tarReader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			checkErr("Extract error::ReadTarArchive", err)
		}

		filename := ""
		if header.Name[:2] == "go" {
			// filename = dest + ps + strings.Replace(header.Name, "go", "go"+ver, 1)
			filename = dest + ps + strings.Replace(header.Name, "go", "go"+ps+ver, 1)
		} else if header.Name[:2] == "li" {
			// filename = dest + ps + strings.Replace(header.Name, "liteide", "liteide"+ver, 1)
			filename = dest + ps + strings.Replace(header.Name, "liteide", "liteide"+ps+ver, 1)
		}
		switch header.Typeflag {
		case tar.TypeDir:
			err = os.MkdirAll(filename, os.FileMode(header.Mode))
			checkErr("In unGzipBzip2 - MkdirAll", err)
		case tar.TypeReg, tar.TypeRegA:
			writer, err := os.Create(filename)
			checkErr("In unGzipBzip2 - Create", err)
			io.Copy(writer, tarReader)
			err = os.Chmod(filename, os.FileMode(header.Mode))
			checkErr("In unGzipBzip2 - Chmod", err)
			writer.Close()
		case tar.TypeSymlink:
			err = os.Symlink(header.Linkname, filename)
		default:
			fmt.Printf("Unable to untar type: %c in file %s\n", header.Typeflag, filename)
		}
		bar.Increment()
	}
	bar.Finish()
}
