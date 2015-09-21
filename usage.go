// usage
package main

import (
	"fmt"
)

func usage() {
	printBanner()
	fmt.Println("USAGE:")
	fmt.Printf("  %-35s - %s\n", "refresh", "Fetch remote version's list")
	fmt.Printf("  %-35s - %s\n", "ls [go|liteide]", "Show installed version's list")
	fmt.Printf("  %-35s - %s\n", "ls-remote [go|liteide]", "Show remote version's list")
	fmt.Printf("  %-35s - %s\n", "fetch [go|liteide [version]]", "Download only tool's archive")
	fmt.Printf("  %-35s - %s\n", "install [go|liteide [version]]", "Install tool")
	fmt.Printf("  %-35s - %s\n", "remove <go|liteide [version]>", "Remove tool")
	fmt.Printf("  %-35s - %s\n", "archives", "List local archives")
	fmt.Println("")
}
