package main

import (
	"awesomeProject/src/filesystemPart"
	"fmt"
	"os"
	"strings"
)

func main() {
	filesystem, err := filesystemPart.InitFilesystem()
	if err != nil {
		panic("FileSystem not initialized in Main")
	}
	var ans string
	var signname string
	var signpwd string
	for true {
		fmt.Println("[S]ign in or [n]ew user")
		fmt.Scanln(&ans)
		ans = strings.TrimSpace(ans)
		ans = strings.ToLower(ans)
		switch ans {
		case "s":
			{
				fmt.Println("username")
				fmt.Scanln(&signname)
				signname = strings.TrimSpace(signname)
				signname = strings.ToLower(signname)
				fmt.Println("password")
				fmt.Scanln(&signpwd)
				signpwd = strings.TrimSpace(signpwd)
				signpwd = strings.ToLower(signpwd)
				user, err := filesystemPart.CheckUserAndGetUser(signpwd, signname, filesystem)
				if err != nil {
					return
				}
				_, err = user.DeleteImage("beginning")
				//_, err = user.AddImageToRepository("C:\\Users\\18645\\Pictures\\Saved Pictures\\goland.png")
				if err != nil {
					fmt.Println("not working for add images")
					return
				}
				return
			}
		case "n":
			var username string
			var firstname string
			var lastname string
			fmt.Println("username")
			fmt.Scanln(&username)
			fmt.Println("first name")
			fmt.Scanln(&firstname)
			fmt.Println("last name")
			fmt.Scanln(&lastname)
			fmt.Println("password")
			fmt.Scanln(&signpwd)
			_, err = filesystemPart.InitCreateUser(username, firstname, lastname, signpwd, filesystem)
			if err != nil {
				panic(err)
			}
		case "q":
			os.Exit(0)
		}

	}

}
