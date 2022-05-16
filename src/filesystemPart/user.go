package filesystemPart

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

//Users struct
//Fields
//AllUsers: type []Users              //Holds an array of all users who have directories in Repository
type Users struct {
	AllUsers []User `json:"users"`
}

//User Struct
//Holds fields relative to each user.
//Fields:
//Username: type string               //username that was supplied by the user
//Firstname: type string			  //first name supplied by the user
//Lastname: type string				  //last name supplied by the user
//Password: type string				  //initially supplied by the user but later encrypted by the file system
//DirId: type uint32				  //directory id supplied by the filesystem
//UserStorage: type int				  //supplied by the filesystem amount of storage used in bytes
//ImagesIn Repo: type map			  //returns an Image when supplied the name of an Image
//filesys: type *FileSystem           //
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

// InitCreateUser initialize and create User
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

// AddImageToRepository add the image to the repository with the help of the filesystem
// @param path : the full path of the Image to be added
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
