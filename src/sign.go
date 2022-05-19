package main

import (
	"awesomeProject/src/filesystemPart"
	"fmt"
	"golang.org/x/term"
	_ "golang.org/x/term"
	"os"
	"strings"
	"syscall"
)

func initSignIn(f *filesystemPart.Filesystem) (bool, error) {
	var ans string
	fmt.Println("[s]ign in or [n]ew user or [q]uit")
	_, err := fmt.Scanln(&ans)
	if err != nil {
		return false, err
	}
	ans = strings.TrimSpace(ans)
	ans = strings.ToLower(ans)
	switch ans {
	case "s":
		_, err := signIn(f)
		if err != nil {
			return false, err
		}
	case "n":
		_, err := signUp(f)
		if err != nil {
			return false, err
		}
		_, err = signIn(f)
		if err != nil {
			return false, err
		}
	case "q":
		fmt.Println("Exiting Repository....")
		os.Exit(0)
	}
	return true, nil
}

func signIn(f *filesystemPart.Filesystem) (bool, error) {
	var signname string
	fmt.Println("username")
	_, err := fmt.Scanln(&signname)
	if err != nil {
		return false, err
	}
	signname = strings.TrimSpace(signname)
	fmt.Println("password (will not echo)")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	user, err := filesystemPart.CheckUserAndGetUser(string(bytePassword), signname, f)
	if err != nil {
		return false, err
	}
	_, err = runFileSystem(user)
	if err != nil {
		return false, err
	}
	return true, nil
}

func signUp(f *filesystemPart.Filesystem) (bool, error) {
	var username string
	var firstname string
	var lastname string
	fmt.Println("username")
	_, err := fmt.Scanln(&username)
	if err != nil {
		return false, err
	}
	fmt.Println("first name")
	_, err = fmt.Scanln(&firstname)
	if err != nil {
		return false, err
	}
	fmt.Println("last name")
	_, err = fmt.Scanln(&lastname)
	if err != nil {
		return false, err
	}
	fmt.Println("password (will not echo)")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return false, err
	}
	_, err = filesystemPart.InitCreateUser(username, firstname, lastname, string(bytePassword), f)
	if err != nil {
		return false, err
	}
	return true, nil
}

func runFileSystem(user *filesystemPart.User) (bool, error) {
	var answer string
	for {
		fmt.Println("Type")
		fmt.Println("a ==>To add image")
		fmt.Println("d ==>To delete image")
		fmt.Println("q ==>To quit")
		fmt.Println("delete ==>To Remove yourself as a user")

		_, err := fmt.Scanln(&answer)
		if err != nil {
			return false, err
		}
		switch answer {
		case "a":
			_, err := addImage(user)
			if err != nil {
				fmt.Printf("error occured adding image. %v\n", err)
				continue
			}
			fmt.Println("Image successfully added")
			continue
		case "d":
			_, err := deleteImage(user)
			if err != nil {
				fmt.Printf("error occured deleting image. %v\n", err)
				continue
			}
			fmt.Println("Image deleted successfully")
			continue
		case "delete":
			_, err := user.RemoveMyUserAccount()
			if err != nil {
				fmt.Println("Sorry to see you go")
				return true, nil
			}
		case "l":
			_, err := user.ListImages()
			if err != nil {
				return false, err
			}
			continue
		case "q":
			fmt.Println("Exiting.....")
			os.Exit(0)
		}

	}
}

func deleteImage(user *filesystemPart.User) (bool, error) {
	fmt.Println("type the name of the image (with no extension)")
	var name string
	_, err := fmt.Scanln(&name)
	if err != nil {
		return false, err
	}
	_, err = user.DeleteImage(name)
	if err != nil {
		return false, err
	}
	return true, nil
}

func addImage(user *filesystemPart.User) (bool, error) {
	fmt.Println("type the absolute path to the image")
	var path string
	_, err := fmt.Scanln(&path)
	if err != nil {
		return false, err
	}
	_, err = user.AddImageToRepository(path)
	if err != nil {
		return false, err
	}
	return true, nil
}
