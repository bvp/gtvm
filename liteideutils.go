package main

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"regexp"
	"strings"
)

type liteIDEVer struct {
	ver       string
	updatedAt string
}

type liteIDEfile struct {
	verMajor    string
	verMinor    string
	osPlatform  string
	osArch      string
	qtType      string
	qtVer       string
	fileName    string
	fileNameURL string
	updatedAt   string
}

func cacheLiteIDE() {
	liteIDECrawler(urlLiteIDE, "", false)
	fetchLiteIDEfileList()
}

func liteIDECrawler(url string, ver string, getFileList bool) {
	var liteideRegex = regexp.MustCompile(`(?P<Prefix>liteidex)(?P<Major>\d+)?(?:\.|-)+(?P<Minor>(?:\d+(?:(\.|-)\d+)?))?(?:\.)?(?P<Platform>windows|linux|macosx)?(?:-)?(?P<Arch>(amd)?\d+)?(?:-)?(?P<Variant>system)?(?:-)?(?P<QtVer>qt\d)?(?:-)?(?P<Variant2>system)?(?P<Ext>7z|zip|tar)?`)
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
				updatedAt := s.Find("td[headers=files_date_h] abbr").Text()
				if name != "" && name != "Parent folder" {
					if getFileList {
						var verMajor, verMinor, osPlatform, osArch, qtType, qtVer string
						match := liteideRegex.FindStringSubmatch(name)
						result := make(map[string]string)
						for i, name := range liteideRegex.SubexpNames() {
							result[name] = match[i]
						}
						verMajor = result["Major"]
						verMinor = result["Minor"]
						osPlatform = result["Platform"]
						if result["Arch"] == "amd64" {
							osArch = "64"
						} else if result["Arch"] == "386" {
							osArch = "32"
						} else {
							osArch = result["Arch"]
						}
						qtVer = result["QtVer"]
						if qtType = result["Variant"]; result["Variant2"] != "" {
							qtType = result["Variant2"]
						}

						lf := liteIDEfile{verMajor: verMajor, verMinor: verMinor, osPlatform: osPlatform, osArch: osArch, qtType: qtType, qtVer: qtVer, fileName: name, fileNameURL: sURL, updatedAt: updatedAt}
						liteIDEfileList = append(liteIDEfileList, lf)
					} else {
						lv := liteIDEVer{ver: name, updatedAt: updatedAt}
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
	stmt, err = tx.Prepare("insert into liteideCache(ver, osPlatform, osArch, qtType, qtVer, fileName, url, updated_at) values (?, ?, ?, ?, ?, ?, ?, ?)")
	checkErr("In fetchLiteIDEfileList - Prepare statement", err)
	defer stmt.Close()

	for _, v := range liteIDEVersions {
		liteIDECrawler(urlLiteIDE, v.ver, true)
		for _, f := range liteIDEfileList {
			var ver string
			if ver = f.verMajor; f.verMinor != "" {
				ver = f.verMajor + "." + f.verMinor
			}
			_, err = stmt.Exec(ver, f.osPlatform, f.osArch, f.qtType, f.qtVer, f.fileName, f.fileNameURL, f.updatedAt)
			checkErr("In fetchLiteIDEfileList - Exec statement", err)
		}
	}
	tx.Commit()
}
