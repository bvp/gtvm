/*
GTVM is Golang Tools Version Manager
Manage Golang versions and LiteIDE versions install/uninstall


Basics

All configs, archives, installed tools stored in $HOME/.gtvm in Linux and %USERPROFILE%\.gtvm in Windows


*/
package main

import (
	"database/sql"
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/mattn/go-sqlite3"
)

const (
	gtvmDirName        = ".gtvm"
	golang             = "go"
	archives           = "archives"
	liteide            = "liteide"
	gtvmStorage        = "gtvmStorage.db"
	urlSF              = "http://sourceforge.net"
	urlLiteIDE         = "http://sourceforge.net/projects/liteide/files/"
	urlLiteIDEDownload = "http://downloads.sourceforge.net/liteide/"
	urlGoLang          = "https://golang.org/dl/"
)

var (
	ps              = string(filepath.Separator) // Separator
	goVersions      []goVer                      // Store Go versions info
	liteIDEVersions []liteIDEVer                 // Store LiteIDE versions info
	liteIDEfileList []liteIDEfile                // Store LiteIDE files info
	doc             *goquery.Document
	err             error
	db              *sql.DB
	stmt            *sql.Stmt

	gtvmDir     = ""
	archivesDir = ""
	liteideDir  = ""
	golangDir   = ""
)

func init() {
	usr, _ := user.Current()
	gtvmDir = usr.HomeDir + ps + gtvmDirName
	archivesDir = gtvmDir + ps + archives
	golangDir = gtvmDir + ps + golang
	liteideDir = gtvmDir + ps + liteide

	createWorkDirs()

	gvmwd, err := os.Stat(gtvmDir)

	if err != nil {
		fmt.Println(err)
	}

	if !gvmwd.IsDir() {
		fmt.Println("Go Tools Version Manager working destination is not a directory")
		//		os.Exit(1)
	}

	if _, err = os.Stat(gtvmDir + ps + gtvmStorage); os.IsNotExist(err) {
		// log.Println("firstStart()")
		firstStart()
	} else {
		//	os.Remove(gtvmDir + ps + gtvmStorage)
		db = getDB()
	}

}

func main() {
	parseCmdLine()

	//	download("http://downloads.sourceforge.net/liteide/X27.1/liteidex27.1.linux-64-system-qt4.8.tar.bz2", "liteidex27.1.linux-64-system-qt4.8.tar.bz2")
	//	download("https://storage.googleapis.com/golang/go1.4.2.linux-amd64.tar.gz", "go1.4.2.linux-amd64.tar.gz")
	//	download("http://downloads.sourceforge.net/liteide/X27.1/liteidex27.1.windows.zip", "liteidex27.1.windows.zip")
	//  checksum("go1.4.2.windows-amd64.zip")
	//	unzip(gvmDir+ps+archivesDir+ps+"liteidex27.1.windows.zip", gvmDir+ps+liteideDir, "27.1")
	//unGzipBzip2(gvmDir+ps+archivesDir+ps+"go1.4.2.linux-amd64.tar.gz", gvmDir+ps+goDir, "1.4.2")
	//unGzipBzip2(gvmDir+ps+archivesDir+ps+"liteidex27.1.linux-64-system-qt4.8.tar.bz2", gvmDir+ps+liteideDir, "27.1")

	//	fmt.Println("Work dir - " + gvmDir)
	//	listArchives()
	//	printVersions("go")
	//	printVersions("liteide")
	defer db.Close()
}

func getDB() *sql.DB {
	if db == nil {
		db, err = sql.Open("sqlite3", gtvmDir+ps+gtvmStorage)
		if err != nil {
			panic(err)
		}
	}

	return db
}
