package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mvanwaaijen/fileinfo"
)

func main() {
	filePath := flag.String("path", "", "please specify file path to .exe or .dll to show file info.")
	flag.Parse()
	if filePath == nil || len(*filePath) == 0 {
		flag.Usage()
		os.Exit(1)
	}
	fi, _ := fileinfo.New(*filePath)
	fmt.Printf("Description : %v\n", fi.GetFileDesc())
	fmt.Printf("File Version: %v\n", fi.GetFileVer())
	fmt.Printf("Product Name: %v\n", fi.GetProdName())
	fmt.Printf("Prod Version: %v\n", fi.GetProdVer())
	fmt.Printf("Org Filename: %v\n", fi.GetOrgName())
	fmt.Printf("File Hash   : %v\n", fi.GetHash())
}
