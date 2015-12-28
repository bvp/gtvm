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

	"github.com/badgerodon/penv"
)

type latest struct {
	ver        string
	url        string
	fileName   string
	osPlatform string
	osArch     string
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
	for _, v := range vers {
		fmt.Printf("- %s\n", v)
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

func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}

func listArchives() {
	fmt.Println(strArchivesListing)
	files, _ := ioutil.ReadDir(archivesDir)
	for _, f := range files {
		fmt.Println(f.Name())
	}
}

func listInstalled(gt string) []string {
	var (
		inDir      string
		sInstalled []string
	)

	fmt.Printf(strInstalledVersions, gt)
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
			fmt.Printf("\t%s (%s)\n", f.Name(), f.ModTime().Format("2006-01-02 15:04:05"))
			sInstalled = append(sInstalled, f.Name())
		}
	}
	return sInstalled
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

func compareHash(ver, hash string) bool {
	re := regexp.MustCompile("[^0-9]")

	curOS := runtime.GOOS
	curArch := re.ReplaceAllString(runtime.GOARCH, "")

	var mhash string
	var result bool
	stmt, _ := db.Prepare("select hash from golangCache where osPlatform = ? and osArch = ?")
	defer stmt.Close()
	err = stmt.QueryRow(curOS, curArch).Scan(&mhash)
	fmt.Println("* mhash - '" + mhash + "' and hash - '" + hash + "'")
	if mhash != "" {
		if strings.EqualFold(hash, mhash) {
			result = true
		} else {
			fmt.Println("\n" + strChecksumMismatch)
			result = false
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
	if contains(listInstalled("go"), goRoot) {
		err := penv.SetEnv("GOROOT", golangDir+ps+goRoot)
		checkErr("setGoRoot", err)
	} else {
		gtver := getLatest("go", goRoot, "")
		download("golang", gtver.url, gtver.fileName)
		// compareHash(gtver.ver, checksum(archivesDir+ps+gtver.fileName))
		fmt.Printf(strInstallGoVersion, goRoot)
		extract(gtver.fileName, gtver.ver)
		err := penv.SetEnv("GOROOT", golangDir+ps+goRoot)
		checkErr("setGoRoot", err)
	}
}
