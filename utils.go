// utils
package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"runtime"
	"strings"
)

type latest struct {
	ver      string
	url      string
	fileName string
}

func removeDuplicates(xs *[]string) {
	found := make(map[string]bool)
	j := 0
	for i, x := range *xs {
		if !found[x] {
			found[x] = true
			(*xs)[j] = (*xs)[i]
			j++
		}
	}
	*xs = (*xs)[:j]
}

func printVersions(t string) {
	var (
		msg  string
		vers []string
		rows *sql.Rows
		err  error
	)

	if t == "go" {
		msg = "Golang versions:"
		rows, err = db.Query("select ver from golangCache")
		checkErr("In printVersions - Query", err)
	} else if t == "liteide" {
		msg = "LiteIDE versions:"
		rows, err = db.Query("select ver from liteideCache")
		checkErr("In printVersions - Query", err)
	}
	defer rows.Close()

	for rows.Next() {
		var ver string
		rows.Scan(&ver)
		vers = append(vers, ver)
	}
	removeDuplicates(&vers)
	//	fmt.Printf("%q\n", vers)
	rows.Close()

	fmt.Println(msg)
	for _, v := range vers {
		fmt.Printf("- %s\n", v)
	}

	//	for _, v := range liteIDEVersions {
	//		fmt.Printf("%q\n", v)
	//	}
}

// exists returns whether the given file or directory exists or not
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func createWorkDirs() {
	os.MkdirAll(archivesDir, 0777)
	os.MkdirAll(golangDir, 0777)
	os.MkdirAll(liteideDir, 0777)
}

func createTables() {
	sqlStmt := `
		--drop table if exists config;
		--drop table if exists golangCache;
		--drop table if exists liteideCache;
		create table if not exists config (id integer not null primary key, name text, gopath text, gobin text);
		create table if not exists golangCache (ver text, osType text, osArch text, kind text, url text, fileName text, size text, hash text);
		create table if not exists liteideCache (ver text, osType text, osArch text, qtType text, qtVer text, fileName text, url text, updated_at text);
		--delete from config;
		`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
}

func checkErr(msg string, err error) {
	if err != nil {
		log.Printf("%s : %s\n", msg, err.Error())
		os.Exit(1)
	}
}

func listArchives() {
	fmt.Println("Archives listing")
	files, _ := ioutil.ReadDir(archivesDir)
	for _, f := range files {
		fmt.Println(f.Name())
	}
}

func listInstalled(gt string) {
	fmt.Printf("*Installed %s versions\n", gt)
	var inDir string = ""
	if gt == "go" {
		inDir = golangDir
	} else if gt == "liteide" {
		inDir = liteideDir
	}
	files, _ := ioutil.ReadDir(inDir)
	for _, f := range files {
		fi, err := os.Stat(inDir + ps + f.Name())
		checkErr("Getting os.Stat", err)
		if fi.IsDir() {
			fmt.Printf("\t%s\n", f.Name())
		}
	}
}

func refreshDb() {
	sqlStmt := `
		drop table if exists golangCache;
		drop table if exists liteideCache;
		create table if not exists golangCache (ver text, osType text, osArch text, kind text, url text, fileName text, size text, hash text);
		create table if not exists liteideCache (ver text, osType text, osArch text, qtType text, qtVer text, fileName text, url text, updated_at text);
		`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	fmt.Print("Updating Golang cache...\t")
	cacheGoLang(urlGoLang)
	fmt.Println("OK")

	fmt.Print("Updating LiteIDE cache...\t")
	cacheLiteIDE()
	fmt.Println("OK")
}

func firstStart() {
	fmt.Printf("This is first start\n")
	fmt.Printf("* Creating work directory\n")
	createWorkDirs()
	fmt.Printf("* Creating storage\n")
	db = getDB()
	createTables()
	refreshDb()
}

func printBanner() {
	fmt.Println("")
	fmt.Println("     //////  //////////  //      //  //      //   ")
	fmt.Println("  //            //      //      //  ////  ////    ")
	fmt.Println(" //  ////      //      //      //  //  //  //     ")
	fmt.Println("//    //      //        //  //    //      //      ")
	fmt.Println(" //////      //          //      //      //       ")
	fmt.Println("")
	fmt.Println("Golang Tools Version Manager")
	fmt.Printf("Command line tool for manage Golang & LiteIDE versions\n\n")
}

func compareHash(ver, hash string) bool {
	re := regexp.MustCompile("[^0-9]")

	curOS := runtime.GOOS
	curArch := re.ReplaceAllString(runtime.GOARCH, "")

	var mhash string
	stmt, _ := db.Prepare("select hash from golangCache where osType = ? and osArch = ?")
	defer stmt.Close()
	err = stmt.QueryRow(curOS, curArch).Scan(&mhash)
	//	fmt.Println("Chechsum", mhash, strings.EqualFold(hash, mhash))
	if strings.EqualFold(hash, mhash) {
		//		fmt.Println("Checksum is OK")
		return true
	} else {
		fmt.Println("Checksum missmatch")
		return false
	}
}

func getLatest(t string) latest {
	var (
		verurl latest
	)
	re := regexp.MustCompile("[^0-9]")

	curOS := runtime.GOOS
	curArch := re.ReplaceAllString(runtime.GOARCH, "")

	if t == "go" {
		stmt, err = db.Prepare("select ver, url, fileName from golangCache where osType = ? and osArch = ?")
		if err != nil {
			log.Fatal(err)
		}
	} else if t == "liteide" {
		// stmt, err = db.Prepare("select ver, url, fileName from liteideCache where osType = ? and osArch = ?")
		stmt, err = db.Prepare("select ver, url, fileName from liteideCache where osType = ?")
		if err != nil {
			log.Fatal(err)
		}
	}

	defer stmt.Close()
	err = stmt.QueryRow(curOS, curArch).Scan(&verurl.ver, &verurl.url, &verurl.fileName)
	if err != nil {
		log.Fatal(err)
	}

	return verurl
}
