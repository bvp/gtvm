package main

import (
	"fmt"
	"os"
	"strings"
)

func parseCmdLine() {
	var gtver latest
	if len(os.Args) > 1 {
		args := os.Args[1:]
		//		fmt.Printf("%q - length %d\n", args, len(args))

		switch strings.ToLower(args[0]) {
		case "refresh":
			fmt.Println(strFetchRemote)
			refreshDb()
		case "installed":
			if len(args) > 1 {
				if args[1] == "go" {
					listInstalled("go")
				} else if args[1] == "liteide" {
					listInstalled("liteide")
				}
			} else {
				listInstalled("go")
				listInstalled("liteide")
			}
		case "ls":
			if len(args) > 1 {
				if args[1] == "go" {
					printVersions("go")
				} else if args[1] == "liteide" {
					printVersions("liteide")
				}
			} else {
				printVersions("go")
				printVersions("liteide")
			}
		case "fetch":
			if len(args) >= 2 {
				if args[1] == "go" {
					if len(args) >= 3 {
						gtver = getLatest("go", args[2], "")
						fmt.Printf(strDownloadingGo, gtver.ver)
					} else {
						gtver = getLatest("go", "", "")
						fmt.Printf(strDownloadingGo, gtver.ver)
					}
					download("golang", gtver.url, gtver.fileName)
					// compareHash(gtver.ver, checksum(archivesDir+ps+gtver.fileName))
				} else if args[1] == "liteide" {
					if len(args) >= 3 {
						fmt.Printf(strDownloadingLiteIDE, args[2])
						gtver = getLatest("liteide", args[2], "")
						if len(args) >= 4 {
							fmt.Printf(strDownloadingLiteIDEversion, args[2], args[3])
							gtver = getLatest("liteide", args[2], args[3])
						}
					} else {
						gtver = getLatest("liteide", "", "")
					}
					fmt.Printf(strDownloading, gtver.ver)
					download("liteide", gtver.url, gtver.fileName)
				} else {
					usage()
				}
			} else if len(args) == 1 {
				fmt.Println(strPleaseSetTool)
			} else {
				usage()
			}
		case "install", "i":
			if len(args) > 1 {
				if args[1] == "go" {
					if len(args) >= 3 {
						gtver = getLatest("go", args[2], "")
					} else {
						gtver = getLatest("go", "", "")
					}
					//  fmt.Printf("Version: %s, URL: %s, Filename: %s\n", gtver.ver, gtver.url, gtver.fileName)
					fmt.Printf(strDownloadingGo, gtver.ver)
					download("golang", gtver.url, gtver.fileName)
					// compareHash(gtver.ver, checksum(archivesDir+ps+gtver.fileName))
					extract(gtver.fileName, gtver.ver)
				} else if args[1] == "liteide" {
					if len(args) >= 3 {
						fmt.Printf(strDownloadingLiteIDE, args[2])
						gtver = getLatest("liteide", args[2], "")
						if len(args) >= 4 {
							fmt.Printf(strDownloadingLiteIDEversion, args[2], args[3])
							gtver = getLatest("liteide", args[2], args[3])
						}
					} else {
						gtver = getLatest("liteide", "", "")
					}
					fmt.Printf(strDownloadingInstalling, gtver.ver)
					download("liteide", gtver.url, gtver.fileName)
					extract(gtver.fileName, gtver.ver)
				} else {
					usage()
				}
			} else {
				fmt.Println(strPleaseSetTool)
				usage()
			}
		case "uninstall", "u":
			if len(args) == 3 {
				if args[1] == "go" {
					fmt.Println(strUninstallGo, args[2])
					errRemove := os.RemoveAll(gtvmDir + ps + "go" + ps + args[2])
					checkErr("Uninstall go", errRemove)
				} else if args[1] == "liteide" {
					fmt.Println(strUninstallLiteIDE, args[2])
					errRemove := os.RemoveAll(gtvmDir + ps + "liteide" + ps + args[2])
					checkErr("Uninstall liteide", errRemove)
				} else {
					usage()
				}
			} else if len(args) == 1 {
				fmt.Println(strPleaseSetTool)
			} else {
				usage()
			}
		case "use":
			if len(args) > 1 {
				if args[1] == "go" {
					if len(args) >= 3 {
						setGoRoot(args[2])
					}
				}
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
		case "help", "h":
			usage()
		}
	} else {
		usage()
	}
}
