package main

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"regexp"
	"strings"
)

type goVer struct {
	ver        string
	osPlatform string
	osArch     string
	kind       string
	url        string
	fileName   string
	size       string
	hash       string
}

func cacheGoLang(url string) {
	re := regexp.MustCompile("[^0-9]")
	goVersions = nil
	//	checkErr(err)
	tx, txerr := db.Begin()
	checkErr("In cacheGoLang - Begin transaction", txerr)
	stmt, err := tx.Prepare("insert into golangCache(ver, osPlatform, osArch, kind, url, fileName, size, hash) values (?, ?, ?, ?, ?, ?, ?, ?)")
	checkErr("In cacheGoLang - Prepare statement", txerr)
	defer stmt.Close()

	if doc, err = goquery.NewDocument(url); err != nil {
		log.Fatal(err)
	}
	doc.Find("#page .container div[class*='toggle']").Each(func(i int, s *goquery.Selection) {
		ver := strings.TrimRight(s.Find(".collapsed h2").Text(), " ▹")

		s.Find(".codetable").Each(func(i int, s *goquery.Selection) {
			s.Find("tr").Each(func(i int, s *goquery.Selection) {
				url, _ := s.Find("td.filename a").Attr("href")
				fileName := s.Find("td.filename").Text()
				kind := s.Find("td:nth-of-type(2)").Text()
				osPlatform := strings.ToLower(s.Find("td:nth-of-type(3)").Text())
				osArch := re.ReplaceAllString(s.Find("td:nth-of-type(4)").Text(), "")
				size := s.Find("td:nth-of-type(5)").Text()
				hash := s.Find("td:nth-of-type(6)").Text()
				if url != "" && kind != "Source" {
					gv := goVer{ver: ver[2:], osPlatform: osPlatform, osArch: osArch, kind: kind, url: url, fileName: fileName, size: size, hash: hash}
					goVersions = append(goVersions, gv)
					_, err = stmt.Exec(ver[2:], osPlatform, osArch, kind, url, fileName, size, hash)
					checkErr("In cacheGoLang - Exec statement", txerr)
				}
			})
		})
	})
	tx.Commit()
}
