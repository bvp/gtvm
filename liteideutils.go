// liteideutils
package main

import (
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type liteIDEVer struct {
	ver        string
	updated_at string
}

type liteIDEfile struct {
	ver         string
	osType      string
	osArch      string
	qtType      string
	qtVer       string
	fileName    string
	fileNameURL string
	updated_at  string
}

func cacheLiteIDE() {
	liteIDECrawler(urlLiteIDE, "", false)
	fetchLiteIDEfileList()
}

func liteIDECrawler(url string, ver string, getFileList bool) {
	src := ""
	if !getFileList {
		src = url
		liteIDEVersions = nil
	} else {
		src = url + ver
		liteIDEfileList = nil
	}

	if doc, err = goquery.NewDocument(src); err != nil {
		log.Fatal(err)
	}
	doc.Find("#files").Each(func(i int, s *goquery.Selection) {
		s.Find("#files_list").Each(func(i int, s *goquery.Selection) {
			s.Find("tr").Each(func(i int, s *goquery.Selection) {
				//				sURL, _ := s.Find("th a").Attr("href")
				name := strings.TrimSpace(s.Find("th a").Text())
				sURL := urlLiteIDEDownload + ver + "/" + name
				updated_at := s.Find("td[headers=files_date_h] abbr").Text()
				if name != "" && name != "Parent folder" {
					//					fmt.Println(name)

					if getFileList {
						var osType, osArch, qtType, qtVer string
						platform := getPlatform(name, ver[1:])
						if len(platform) == 4 {
							osType = platform[0]
							osArch = platform[1]
							qtType = platform[2]
							qtVer = platform[3][2:]
						} else if len(platform) == 2 {
							osType = platform[0]
							osArch = platform[1]
						} else {
							osType = platform[0]
						}

						lf := liteIDEfile{ver: ver[1:], osType: osType, osArch: osArch, qtType: qtType, qtVer: qtVer, fileName: name, fileNameURL: sURL, updated_at: updated_at}
						//						fmt.Printf("%q\n", lf)
						liteIDEfileList = append(liteIDEfileList, lf)
					} else {
						lv := liteIDEVer{ver: name, updated_at: updated_at}
						liteIDEVersions = append(liteIDEVersions, lv)
					}
				}
			})
		})
	})
}

func fetchLiteIDEfileList() {
	tx, txerr := db.Begin()
	checkErr("In fetchLiteIDEfileList - Begin transaction", txerr)
	stmt, err = tx.Prepare("insert into liteideCache(ver, osType, osArch, qtType, qtVer, fileName, url, updated_at) values (?, ?, ?, ?, ?, ?, ?, ?)")
	checkErr("In fetchLiteIDEfileList - Prepare statement", err)
	defer stmt.Close()

	for _, v := range liteIDEVersions {
		liteIDECrawler(urlLiteIDE, v.ver, true)
		for _, f := range liteIDEfileList {
			_, err = stmt.Exec(f.ver, f.osType, f.osArch, f.qtType, f.qtVer, f.fileName, f.fileNameURL, f.updated_at)
			checkErr("In fetchLiteIDEfileList - Exec statement", err)
		}
	}
	tx.Commit()
}

func getPlatform(s, ver string) []string {
	platform := ""
	prefix := ""
	suffix := ""
	liteidex := "liteidex"

	prefixes := []string{liteidex + ver + ".", liteidex + ver + "-1" + "."}
	suffixes := []string{".tar.bz2", ".7z", ".zip"}

	for _, prefix = range prefixes {
		if strings.HasPrefix(s, prefix) {
			platform = s[len(prefix):]
		}
	}
	for _, suffix = range suffixes {
		if strings.HasSuffix(platform, suffix) {
			platform = platform[:len(platform)-len(suffix)]
			break
		}
	}
	splitted := strings.Split(platform, "-")
	//	fmt.Printf("Length %q is %d", splitted, len(splitted))
	return splitted
}
