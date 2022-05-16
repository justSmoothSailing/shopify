package main

import (
	"awesomeProject/src/filesystemPart"
	"fmt"
)

func main() {
	filesystem, err := filesystemPart.InitFilesystem()
	if err != nil {
		panic("FileSystem not initialized in Main")
	}
	var username string
	var firstname string
	var lastname string
	var password string
	fmt.Println("username")
	fmt.Scanln(&username)
	fmt.Println("first name")
	fmt.Scanln(&firstname)
	fmt.Println("last name")
	fmt.Scanln(&lastname)
	fmt.Println("password")
	fmt.Scanln(&password)
	_, err = filesystemPart.InitCreateUser(username, firstname, lastname, password, filesystem)
	if err != nil {
		panic(err)
	}

}
