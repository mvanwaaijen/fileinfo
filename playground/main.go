package main

import (
	"flag"
	"fmt"
	"os"

	w32testing "github.com/mvanwaaijen/fileinfo"
)

func main() {
	filePath := flag.String("path", "", "specify file path to .exe or .dll to show file info.")
	flag.Parse()
	if filePath == nil || len(*filePath) == 0 {
		flag.Usage()
		os.Exit(1)
	}
	fi := w32testing.New(*filePath)
	fmt.Printf("Description : %v\n", fi.GetFileDesc())
	fmt.Printf("File Version: %v\n", fi.GetFileVer())
	fmt.Printf("Product Name: %v\n", fi.GetProdName())
	fmt.Printf("Prod Version: %v\n", fi.GetProdVer())
	fmt.Printf("Org Filename: %v\n", fi.GetOrgName())
}
