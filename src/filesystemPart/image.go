package filesystemPart

import (
	"errors"
	"os"
)

type Image struct {
	name        string
	size        int64
	path        string
	origPath    string
	nameExt     string
	owner       string
	permissions int
}

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
