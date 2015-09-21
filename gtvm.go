package main

import (
	"database/sql"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	//	"runtime"
	"github.com/PuerkitoBio/goquery"
	_ "github.com/mattn/go-sqlite3"
)

var (
	ps                 = string(filepath.Separator)
	gtvmDirName        = ".gtvm"
	golang             = "go"
	archives           = "archives"
	liteide            = "liteide"
	gtvmStorage        = "gtvmStorage.db"
	goVersions         []goVer
	liteIDEVersions    []liteIDEVer
	liteIDEfileList    []liteIDEfile
	urlSF              = "http://sourceforge.net"
	urlLiteIDE         = "http://sourceforge.net/projects/liteide/files/"
	urlLiteIDEDownload = "http://downloads.sourceforge.net/liteide/"
	urlGoLang          = "https://golang.org/dl/"
	doc                *goquery.Document
	err                error
	db                 *sql.DB
	stmt               *sql.Stmt

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
		log.Println(err)
	}

	if !gvmwd.IsDir() {
		log.Println("Go Tools Version Manager working destination is not a directory")
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
	runtime.GOMAXPROCS(runtime.NumCPU())
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
