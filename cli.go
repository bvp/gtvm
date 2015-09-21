// cli
package main

import (
	"fmt"
	"os"
	"strings"
)

func parseCmdLine() {
	if len(os.Args) > 1 {
		//		fmt.Println("Command-line params")
		args := os.Args[1:]
		//		fmt.Printf("%q - length %d\n", args, len(args))

		switch strings.ToLower(args[0]) {
		case "refresh":
			fmt.Println("Fetch remote version's list")
			refreshDb()
		case "ls":
			if len(args) == 2 {
				if args[1] == "go" {
					listInstalled("go")
				} else if args[1] == "liteide" {
					listInstalled("liteide")
				}
			} else if len(args) == 1 {
				listInstalled("go")
				listInstalled("liteide")
			}
		case "ls-remote":
			if len(args) == 2 {
				if args[1] == "go" {
					printVersions("go")
				} else if args[1] == "liteide" {
					printVersions("liteide")
				}
			} else if len(args) == 1 {
				printVersions("go")
				printVersions("liteide")
			}
		case "fetch":
			if len(args) == 2 {
				if args[1] == "go" {
					gtver := getLatest("go")
					//  fmt.Printf("Version: %s, URL: %s, Filename: %s\n", gtver.ver, gtver.url, gtver.fileName)
					fmt.Printf("Downloading version: %s\n", gtver.ver)
					download("golang", gtver.url, gtver.fileName)
					compareHash(gtver.ver, checksum(archivesDir+ps+gtver.fileName))
				} else if args[1] == "liteide" {
					gtver := getLatest("liteide")
					fmt.Printf("Downloading version: %s\n", gtver.ver)
					download("liteide", gtver.url, gtver.fileName)
				} else {
					usage()
				}
			} else if len(args) == 1 {
				fmt.Println("Please, set 'go' or 'liteide'")
			} else {
				usage()
			}
		case "install":
			if len(args) == 2 {
				if args[1] == "go" {
					gtver := getLatest("go")
					//  fmt.Printf("Version: %s, URL: %s, Filename: %s\n", gtver.ver, gtver.url, gtver.fileName)
					fmt.Printf("Downloading version: %s\n", gtver.ver)
					download("golang", gtver.url, gtver.fileName)
					compareHash(gtver.ver, checksum(archivesDir+ps+gtver.fileName))
					extract(gtver.fileName, gtver.ver)
				} else if args[1] == "liteide" {
					gtver := getLatest("liteide")
					fmt.Printf("Downloading version: %s\n", gtver.ver)
					download("liteide", gtver.url, gtver.fileName)
					extract(gtver.fileName, gtver.ver)
				} else {
					usage()
				}
			} else if len(args) == 1 {
				fmt.Println("Please, set 'go' or 'liteide'")
			} else {
				usage()
			}
		case "remove":
			if len(args) == 3 {
				if args[1] == "go" {
					fmt.Println("Remove Golang version", args[2])
				} else if args[1] == "liteide" {
					fmt.Println("Remove Liteide version", args[2])
				} else {
					usage()
				}
			} else if len(args) == 1 {
				fmt.Println("Please, set 'go' or 'liteide'")
			} else {
				usage()
			}
		case "archives":
			listArchives()
		case "config":
			if len(args) == 2 {
				//
			} else if len(args) == 1 {
				//
			} else {
				usage()
			}
		case "env":
			if len(args) == 2 {
				//
			} else if len(args) == 1 {
				//
			} else {
				usage()
			}
		case "help":
			usage()
		}
	} else {
		usage()
	}
}
