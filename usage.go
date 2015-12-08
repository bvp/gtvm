package main

import (
	"fmt"
)

func usage() {
	printBanner()
	fmt.Println("USAGE:")
	fmt.Printf("  %-40s - %s\n", "refresh", "Fetch remote version's list")
	fmt.Printf("  %-40s - %s\n", "ls [go|liteide]", "Show remote version's list")
	fmt.Printf("  %-40s - %s\n", "fetch [go|liteide [version]]", "Download only tool's archive. Without version - fetch latest")
	fmt.Printf("  %-40s - %s\n", "install|i [go|liteide [version]]", "Install tool. Without version - install latest")
	fmt.Printf("  %-40s - %s\n", "uninstall|u <go|liteide [version]>", "Remove tool")
	fmt.Printf("  %-40s - %s\n", "use <go|liteide [version]> [--default]", "Use tool. Without version - use latest")
	fmt.Printf("  %-40s - %s\n", "installed [go|liteide]", "Show installed version's list")
	fmt.Printf("  %-40s - %s\n", "archives", "List local archives")
	fmt.Println()
}
