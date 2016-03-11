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
	"time"

	"github.com/badgerodon/penv"
)

type latest struct {
	ver        string
	url        string
	fileName   string
	osPlatform string
	osArch     string
}

type installed struct {
	gtver goVer
	date  time.Time
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

// print versions from storage
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
	instVers := getInstalled(t)
	for _, v := range vers {
		if contains(instVers, v) {
			fmt.Printf("- %-10s[installed]\n", v)
		} else {
			fmt.Printf("- %s\n", v)
		}

	}
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
		create table if not exists golangCache (ver text, osPlatform text, osArch text, kind text, url text, fileName text, size text, hash text);
		create table if not exists liteideCache (ver text, osPlatform text, osArch text, qtType text, qtVer text, fileName text, url text, updated_at text);
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
		log.SetPrefix("ERROR")
		log.Printf("%s : %s\n", msg, err.Error())
		os.Exit(1)
	}
}

func contains(slice []installed, item string) bool {
	ok := false
	for _, e := range slice {
		if e.gtver.ver == item && e.gtver.osPlatform == runtime.GOOS && e.gtver.osArch == runtime.GOARCH {
			ok = true
		}
	}
	return ok
}

func getArchives() []string {
	var sArchives []string
	fmt.Println(strArchivesListing)
	files, _ := ioutil.ReadDir(archivesDir)
	for _, f := range files {
		sArchives = append(sArchives, f.Name())
		fmt.Println(f.Name())
	}
	return sArchives
}

func getInstalled(gt string) []installed {
	var (
		inDir      string
		sInstalled []installed
	)

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
			sInstalled = append(sInstalled, installed{goVer{ver: f.Name(), osPlatform: runtime.GOOS, osArch: runtime.GOARCH}, f.ModTime()})
		}
	}
	return sInstalled
}

func printInstalled(slice []installed, gt string) {
	fmt.Printf(strInstalledVersions, gt)
	for _, e := range slice {
		fmt.Printf("- %-10s (%s)\n", e.gtver.ver, e.date.Format("2006-01-02 15:04:05"))
	}
}

func refreshDb() {
	sqlStmt := `
		drop table if exists golangCache;
		drop table if exists liteideCache;
		create table if not exists golangCache (ver text, osPlatform text, osArch text, kind text, url text, fileName text, size text, hash text);
		create table if not exists liteideCache (ver text, osPlatform text, osArch text, qtType text, qtVer text, fileName text, url text, updated_at text);
		`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	fmt.Print(strUpdatingGoCache)
	cacheGoLang(urlGoLang)
	fmt.Println(ok)

	fmt.Print(strUpdatingLiteIDECache)
	cacheLiteIDE()
	fmt.Println(ok)
}

func firstStart() {
	printBanner()
	fmt.Printf(strFirstStart)
	fmt.Printf(strCreateWorkDirs)
	createWorkDirs()
	fmt.Printf(strCreateStorage)
	db = getDB()
	createTables()
	refreshDb()
}

func printBanner() {
	fmt.Println()
	fmt.Println("     //////  //////////  //      //  //      //   ")
	fmt.Println("  //            //      //      //  ////  ////    ")
	fmt.Println(" //  ////      //      //      //  //  //  //     ")
	fmt.Println("//    //      //        //  //    //      //      ")
	fmt.Println(" //////      //          //      //      //       ")
	fmt.Println()
	fmt.Println("Golang Tools Version Manager")
	fmt.Printf("Command line tool for manage Golang & LiteIDE versions\n\n")
}

func compareHash(ver string, hash []string) bool {
	re := regexp.MustCompile("[^0-9]")

	curOS := runtime.GOOS
	curArch := re.ReplaceAllString(runtime.GOARCH, "")

	var mhash string
	var result bool
	stmt, _ := db.Prepare("select hash from golangCache where ver = ? and osPlatform = ? and osArch = ?")
	defer stmt.Close()
	err = stmt.QueryRow(ver, curOS, curArch).Scan(&mhash)
	if mhash != "" {
		for _, h := range hash {
			fmt.Printf("* calculated hash - %s (%d), hash from db - %s (%d)\n", h, len(h), mhash, len(mhash))
			if strings.EqualFold(h, mhash) {
				result = true
                break
			} else {
				result = false
			}
		}
	} else {
		result = true
	}
	return result
}

func getLatest(t, ver, qt string) latest {
	var (
		verurl latest
	)
	re := regexp.MustCompile("[^0-9]")

	curOS := runtime.GOOS
	curArch := re.ReplaceAllString(runtime.GOARCH, "")

	if t == "go" {
		if ver != "" {
			stmt, err = db.Prepare("SELECT ver, url, fileName, osPlatform, osArch FROM golangCache WHERE osPlatform = ? AND osArch = ? AND ver = ? LIMIT 1")
			if err != nil {
				log.Fatal(err)
			}
			err = stmt.QueryRow(curOS, curArch, ver).Scan(&verurl.ver, &verurl.url, &verurl.fileName, &verurl.osPlatform, &verurl.osArch)
			if err != nil {
				if err.Error() == "sql: no rows in result set" {
					fmt.Println("Unknown version")
					os.Exit(0)
				} else {
					log.Fatal(err)
				}
			}
		} else {
			stmt, err = db.Prepare("SELECT ver, url, fileName, osPlatform, osArch FROM golangCache WHERE osPlatform = ? AND osArch = ?  LIMIT 1")
			if err != nil {
				log.Fatal(err)
			}
			err = stmt.QueryRow(curOS, curArch).Scan(&verurl.ver, &verurl.url, &verurl.fileName, &verurl.osPlatform, &verurl.osArch)
			if err != nil {
				if err.Error() == "sql: no rows in result set" {
					fmt.Println("Unknown version")
					os.Exit(0)
				} else {
					log.Fatal(err)
				}
			}
		}
	} else if t == "liteide" {
		if ver != "" {
			qt = "%"
			stmt, err = db.Prepare("select ver, url, fileName, osPlatform, osArch from liteideCache where osPlatform = ? and osArch = ? and ver = ? and qtVer=?")
			if err != nil {
				log.Fatal(err)
			}
		} else if qt != "" {
			stmt, err = db.Prepare("select ver, url, fileName, osPlatform, osArch from liteideCache where osPlatform = ? and osArch = ? and ver = ? and qtVer like ?")
			if err != nil {
				log.Fatal(err)
			}
		}
		stmt, err = db.Prepare("select ver, url, fileName, osPlatform, osArch from liteideCache where osPlatform = ? and osArch = ? and ver = ? and qtVer like ?")
		if err != nil {
			log.Fatal(err)
		}
		err = stmt.QueryRow(curOS, curArch, ver, qt).Scan(&verurl.ver, &verurl.url, &verurl.fileName, &verurl.osPlatform, &verurl.osArch)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				fmt.Println("Unknown version")
				os.Exit(0)
			} else {
				log.Fatal(err)
			}
		}
	}
	return verurl
}

func setGoRoot(goRoot string) {
	if contains(getInstalled("go"), goRoot) {
		err := penv.SetEnv("GOROOT", golangDir+ps+goRoot)
		checkErr("setGoRoot", err)
	} else {
		gtver := getLatest("go", goRoot, "")
		download("go", gtver.ver, gtver.url, gtver.fileName)
		if compareHash(gtver.ver, checksum(archivesDir+ps+gtver.fileName)) {
			fmt.Printf(strInstallGoVersion, goRoot)
			extract(gtver.fileName, gtver.ver)
			err := penv.SetEnv("GOROOT", golangDir+ps+goRoot)
			checkErr("setGoRoot", err)
		} else {
			fmt.Println(strChecksumMismatch)
		}
	}
}

func removeFile(f string) {
	err := os.Remove(f)
	checkErr("Remove file", err)
	fmt.Println(strFileRemoved)
}
