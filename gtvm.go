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
		os.Exit(1)
	}

	if _, err = os.Stat(gtvmDir + ps + gtvmStorage); os.IsNotExist(err) {
		firstStart()
	} else {
		//	os.Remove(gtvmDir + ps + gtvmStorage)
		db = getDB()
	}

}

func main() {
	parseCmdLine()

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
