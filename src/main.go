package main

import (
	"fmt"
	"os"
	"shopify/src/filesystemPart"
)

func main() {
	filesystem, err := filesystemPart.InitFilesystem()
	if err != nil {
		panic("FileSystem not initialized in Main")
	}
	_, err = initSignIn(filesystem)
	if err != nil {
		fmt.Printf("Problem with signing in ===> %v\n", err)
		os.Exit(0)
	}
}
