package filesystemPart

import (
	"errors"
	"os"
)

//Image struct
//Fields
//name: type string           //name of image without extension
//size: type int64			  //length of image in bytes
//path: type string			  //path to the image in repository
//origPath: type string       //original path before in the repository
//owner: type string          //owner of the image (User.username)
//permissions: type int       //public or private depending on the user preferences
type Image struct {
	name        string
	size        int64
	path        string
	origPath    string
	nameExt     string
	owner       string
	permissions int
}

//Function that gets called from FileSystem to initialize all the values in the Image struct
//@param pathname: full path to the image
//@param pname: name of the image without extension
//@param imgname: name of the image with extension
//@param imgOwner: username of the user of the image
//@return: returns a pointer to an Image and nil OR nil and error
func initImage(pathname string, pname string, imgname string, imgOwner string) (*Image, error) {
	img := Image{}
	img.origPath = pathname
	img.nameExt = imgname
	img.name = pname
	info, err := os.Stat(pathname)
	if err != nil {
		return nil, errors.New("image has a size of zero")
	}
	img.size = info.Size()
	img.owner = imgOwner
	img.permissions = 1
	return &img, nil
}
