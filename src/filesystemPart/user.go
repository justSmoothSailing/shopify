package filesystemPart

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

type Users struct {
	AllUsers []User `json:"users"`
}

type User struct {
	Username         string           `json:"username"`
	FirstName        string           `json:"firstName"`
	LastName         string           `json:"lastName"`
	Password         string           `json:"Password"`
	DirId            uint32           `json:"DirId"`
	UserStorageUsage int              `json:"UserStorageUsage"`
	ImagesInRepo     map[string]Image `json:"imagesInRepo"`
	filesys          *Filesystem
	metadataPath     string
}

//Initialize and Create User from
//@param uname : string      //is the supplied username of the user
//@param fname : string      //is the supplied lastname of the user
//@param pword : string      //is the supplied password of the user
//@param filesys: Filesystem //is a pointer to the filesystem

func InitCreateUser(uname string, fname string, lname string, pword string, filesys *Filesystem) (*User, error) {
	user := User{Username: uname, FirstName: fname, LastName: lname, Password: "", DirId: 0, UserStorageUsage: 0}
	_, ok := filesys.userInfo[uname]
	if ok {
		return nil, errors.New("username already exists")
	}
	adduser, exists := filesys.createUser(uname, pword, user)
	return adduser, exists
}

func (u *User) AddImageToRepository(path string) (bool, error) {
	//Check if Path is Valid
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, errors.New("image requested does not exist check filepath")
	}
	imageName := filepath.Base(path)
	ext := filepath.Ext(imageName)
	name := strings.TrimSuffix(imageName, ext)
	_, ok := u.ImagesInRepo[name]
	if ok {
		return false, errors.New("image with the same name exists")
	}
	img, err := initImage(path, name, imageName, u.Username)
	if err != nil {
		return false, errors.New("error adding image to repository")
	}
	bl, err := u.filesys.addImg(*img, u.Username)
	if err != nil {
		return false, err
	}
	return bl, err
}
