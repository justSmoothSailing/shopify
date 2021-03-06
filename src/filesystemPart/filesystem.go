package filesystemPart

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"time"
)

//Struct that holds all information about the repository
//Fields
//userInfo: type map          //maps a username to its user struct
//rootDir: type string        //directory path to root of directory
//userMetadata: type map      //maps a username to its metadata
//userPublicData: type map    //maps a username to an array of public images
//rootDirExists: type boolean //if root directory was created correctly
//userFileExists: type boolean//if the json file on users who created an account
//userFilePersist: type string//json file of users who created an account
//key: type []byte             //key in bytes
//usersSoFar: type struct Users//struct with an array of all users

type Filesystem struct {
	userInfo        map[string]*User
	rootDir         string
	userMetadata    map[string]*Metadata
	userPublicData  map[string][]string
	rootDirExists   bool
	userFileExists  bool
	userFilePersist string
	key             []byte
	usersSoFar      Users
}

// Metadata Struct that holds metadata about the user repository
//Field:
//User: type string          //username of a user in the repository
//StorageUsed: type int64    //length in bytes of a user folder
//AmountOfFiles: type int64  //amount of images in a user folder
type Metadata struct {
	User          string `json:"user"`
	StorageUsed   int64  `json:"storageUsed"`
	AmountOfFiles int64  `json:"amountOfFiles"`
}

const key = "super secret key no one knowssss"

// InitFilesystem Function that initializes a filesystem and populates the fields that are used
//from metadata files for information used in functions
//the metadata is used to verify users and use parts in other functions
//@param NONE
//return: *FileSystem, nil OR nil, error         //A pointer to a populated FileSystem(Or Repository)
func InitFilesystem() (*Filesystem, error) {
	var errorCaused error = nil
	filesystem := Filesystem{userInfo: make(map[string]*User),
		rootDir: "C:/Users/18645/Documents/temp/users", userMetadata: make(map[string]*Metadata), // <===== change the root dir to what you want
		userPublicData: make(map[string][]string), rootDirExists: false, userFileExists: false,
		userFilePersist: "C:/Users/18645/Documents/temp/users/users.json", usersSoFar: Users{}} // <===== also change the user persist file to root dir + "/user.json"
	filesystem.key = []byte(key)
	//Check if root directory exists, If not create it
	_, err := os.Stat(filesystem.rootDir)
	if os.IsNotExist(err) {
		errorCaused = os.MkdirAll(filesystem.rootDir, 0777)
		if errorCaused == nil {
			filesystem.rootDirExists = true
		}
		if errorCaused != nil {
			return nil, errorCaused
		}
	}
	// initialize and populate user data into filesystem
	fs, err := initUserData(&filesystem)
	if err != nil {
		return nil, err
	}
	return fs, nil
}

// initUserData initialize and populate filesystem with all user data
//@param filesystem         pointer to an instance of a filesystem
func initUserData(filesystem *Filesystem) (*Filesystem, error) {
	var errorCaused error = nil
	//Check if file of all current users exist, if not create it
	_, err := os.Stat(filesystem.userFilePersist)
	if os.IsNotExist(err) {
		_, errorCaused = os.Create(filesystem.userFilePersist)
		_, err := os.Stat(filesystem.userFilePersist)
		if !os.IsNotExist(err) {
			filesystem.userFileExists = true
		}
		if err != nil {
			return nil, err
		}
	}
	jsonFile, err := os.Open(filesystem.userFilePersist)
	if err != nil {
		errorCaused = err
	}
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {
			errorCaused = err
		}
	}(jsonFile)

	//Check if json file is empty, if so continue, else load all information
	//into filesystem struct field UsersSoFar
	fileInfo, _ := jsonFile.Stat()
	if fileInfo.Size() > 0 {
		byteValue, _ := ioutil.ReadFile(filesystem.userFilePersist)
		err = json.Unmarshal(byteValue, &filesystem.usersSoFar.AllUsers)
		if err == nil {
			for i := 0; i < len(filesystem.usersSoFar.AllUsers); i++ {
				var user = new(User)
				user.UserStorageUsage = filesystem.usersSoFar.AllUsers[i].UserStorageUsage
				user.DirId = filesystem.usersSoFar.AllUsers[i].DirId
				user.Password = filesystem.usersSoFar.AllUsers[i].Password
				user.FirstName = filesystem.usersSoFar.AllUsers[i].FirstName
				user.LastName = filesystem.usersSoFar.AllUsers[i].LastName
				user.Username = filesystem.usersSoFar.AllUsers[i].Username
				user.filesys = filesystem
				dirId := strconv.FormatInt(int64(user.DirId), 10)
				result := "/" + dirId
				user.metadataPath = filesystem.rootDir + result + ".json"
				user.ImgData = Images{}
				user.ImagesInRepo = make(map[string]*Image)
				user.ImgPath = filesystem.usersSoFar.AllUsers[i].ImgPath
				filesystem.userInfo[filesystem.usersSoFar.AllUsers[i].Username] = user
				_, err := filesystem.initMetadata(user)
				if err != nil {
					return nil, err
				}
				_, err = filesystem.initImgData(user)
				if err != nil {
					return nil, err
				}

			}
		} else {
			errorCaused = errors.New("json file could not be parsed")
		}
	}
	if errorCaused != nil {
		return nil, errorCaused
	}
	return filesystem, nil
}

//Creates User and adds them to the filesystem, creates directory, and initializes metadata
//@param uname              //username of the user
//@param pword              //password supplied by the user
//@param user               //original user of the call
//@return *User, error      //if successful, returns pointer to User and nil, otherwise nil, error
func (f *Filesystem) createUser(uname string, pword string, user *User) (*User, error) {
	user.Password = pword
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	//used name of folder for designated user
	user.DirId = r.Uint32()
	dirId := strconv.FormatInt(int64(user.DirId), 10)
	result := "/" + dirId
	user.metadataPath = f.rootDir + result + ".json"
	user.ImgPath = f.rootDir + result + "/" + "img.json"
	err := os.MkdirAll(f.rootDir+result, os.ModeDir)
	if err != nil {
		return nil, err
	}
	//Add user to the filesystem map
	f.userInfo[uname] = user
	_, exists := f.userInfo[uname]
	if !exists {
		return nil, errors.New("user was not added to Repository")
	}
	//Add user to the repository
	err = f.addUser(user)
	if err != nil {
		return nil, err
	}
	//Create Metadata for user
	bl, err := f.createUserMetadata(user.metadataPath)
	if !bl {
		return nil, err
	}
	_, err = f.createUserImgData(user.ImgPath)
	if err != nil {
		return nil, err
	}
	//Initialize all metadata values to 0 except for the username
	metadata := Metadata{user.Username, 0, 0}
	f.userMetadata[user.Username] = &metadata
	jsonFile, err := os.OpenFile(user.metadataPath, os.O_APPEND|os.O_WRONLY, 0777)
	if err != nil {
		return nil, err
	}
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {
		}
	}(jsonFile)
	//update the metadata json file
	content, err := json.MarshalIndent(metadata, "", "")
	if err != nil {
		return nil, err
	}
	err = ioutil.WriteFile(user.metadataPath, content, 0644)
	if err != nil {
		return nil, err
	}

	return user, nil
}

//Creates a file that holds metadata for a particular User
//@param: path                //The path the metadata file will have
//@return: bool, err         // returns true if file was created successfully or false and error if not
func (f *Filesystem) createUserMetadata(path string) (bool, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		_, err = os.Create(path)
		_, err := os.Stat(path)
		if os.IsNotExist(err) {
			return false, errors.New("error creating metadata file in repository")
		}
	}

	return true, nil
}

//Creates a file that holds metadata about the images in the user's repository
//@param: path     path to metadata file
//@return: bool, err
func (f *Filesystem) createUserImgData(path string) (bool, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		_, err = os.Create(path)
		_, err := os.Stat(path)
		if os.IsNotExist(err) {
			return false, errors.New("error creating img data file in repository")
		}
	}
	return true, nil
}

// initMetaData initialize and populate filesystem with all user data
//@param user         pointer to an instance of a user
func (f *Filesystem) initMetadata(user *User) (bool, error) {
	jsonFile, err := os.Open(f.userFilePersist)
	if err != nil {
		return false, err
	}
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {
		}
	}(jsonFile)
	var metadata Metadata
	byteValue, _ := ioutil.ReadFile(user.metadataPath)
	err = json.Unmarshal(byteValue, &metadata)
	if err != nil {
		return false, errors.New("json unmarshal went wrong")
	}
	f.userMetadata[user.Username] = &metadata
	return true, nil
}

// initImgData initialize and populate filesystem with all Image data
//@param filesystem         pointer to an instance of a filesystem
func (f *Filesystem) initImgData(user *User) (bool, error) {
	_, err := os.Stat(user.ImgPath)
	if os.IsNotExist(err) {
		_, err = os.Create(user.ImgPath)
		if err != nil {
			return false, err
		}
		_, err := os.Stat(user.ImgPath)
		if !os.IsNotExist(err) {
			return false, err
		}
	}
	jsonFile, err := os.Open(user.ImgPath)
	if err != nil {
		return false, err
	}
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {
		}
	}(jsonFile)
	fileInfo, _ := jsonFile.Stat()
	if fileInfo.Size() > 0 {
		byteValue, _ := ioutil.ReadFile(user.ImgPath)
		if err != nil {
			return false, err
		}
		err = json.Unmarshal(byteValue, &user.ImgData.AllImages)
		if err == nil {
			user.ImagesInRepo = make(map[string]*Image)
			for i := 0; i < len(user.ImgData.AllImages); i++ {
				var img = new(Image)
				img.Name = user.ImgData.AllImages[i].Name
				img.NameExt = user.ImgData.AllImages[i].NameExt
				img.Size = user.ImgData.AllImages[i].Size
				img.Path = user.ImgData.AllImages[i].Path
				img.Permissions = user.ImgData.AllImages[i].Permissions
				img.OrigPath = user.ImgData.AllImages[i].OrigPath
				img.Owner = user.ImgData.AllImages[i].Owner
				user.ImagesInRepo[img.Name] = img
			}
		}
	}
	return true, nil
}

//removeIndex             remove an element from a particular index and return the slice without it
//@param img              original slice
//@param index            index to remove
func removeImageIndex(img []Image, index int) []Image {
	return append(img[:index], img[index+1:]...)
}

//removeIndex             remove an element from a particular index and return the slice without it
//@param user              original slice
//@param index            index to remove
func removeUserIndex(user []User, index int) []User {
	return append(user[:index], user[index+1:]...)
}

//deleteImgData              delete an Image from a user repository
//@param user:               pointer to the User of the repository
//@param img:                pointer to the Image to be deleted
func (f *Filesystem) deleteImgData(user *User, img *Image) (bool, error) {
	metadata := f.userMetadata[user.Username]
	metadata.StorageUsed -= img.Size
	metadata.AmountOfFiles -= 1
	jsonFile, err := os.OpenFile(user.metadataPath, os.O_APPEND|os.O_WRONLY, 0777)
	if err != nil {
		return false, err
	}
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {
		}
	}(jsonFile)
	//Append user to the array of users and update the json file
	content, err := json.MarshalIndent(*metadata, "", "")
	if err != nil {
		return false, err
	}
	err = ioutil.WriteFile(user.metadataPath, content, 0644)
	if err != nil {
		return false, err
	}

	name := img.Name
	delete(user.ImagesInRepo, name)
	var index int
	for i, image := range user.ImgData.AllImages {
		if image.Name == name {
			index = i
		}
	}
	user.ImgData.AllImages = removeImageIndex(user.ImgData.AllImages, index)
	_, err = os.Create(user.ImgPath)
	if err != nil {
		return false, err
	}
	jsonFile, err = os.OpenFile(user.ImgPath, os.O_APPEND|os.O_WRONLY, 0777)
	if err != nil {
		return false, err
	}
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {
		}
	}(jsonFile)
	//Append image to the array of images and update the json file
	content, err = json.MarshalIndent(user.ImgData.AllImages, "", "")
	if err != nil {
		return false, err
	}
	err = ioutil.WriteFile(user.ImgPath, content, 0644)
	if err != nil {
		return false, err
	}
	err = os.Remove(user.filesys.rootDir + "/" + strconv.FormatInt(int64(user.DirId), 10) + "/" + img.NameExt)
	if err != nil {
		return false, errors.New("error deleting image")
	}
	return true, nil
}

//Adds a new user to the filesystem
//@param User               //Newly populated user
//@return err               //nil if user was added properly, err if not
func (f *Filesystem) addUser(user *User) error {
	var errorCause error
	jsonFile, err := os.OpenFile(f.userFilePersist, os.O_APPEND|os.O_WRONLY, 0777)
	if err != nil {
		errorCause = err
	}
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {
			errorCause = err
		}
	}(jsonFile)
	//Append user to the array of users and update the json file
	f.usersSoFar.AllUsers = append(f.usersSoFar.AllUsers, *user)
	content, err := json.MarshalIndent(f.usersSoFar.AllUsers, "", "")
	if err != nil {
		errorCause = err
	}
	err = ioutil.WriteFile(f.userFilePersist, content, 0644)
	if err != nil {
		errorCause = err
	}
	return errorCause
}

//Add image to a user directory in the repository
//@param img Image                 //The Image the user wants to add to their directory
//@param uname string              //username of the user who is adding the Image
//@return bool, error              //returns true, nil if image was added successfully OR false, error
func (f *Filesystem) addImg(img *Image, uname string) (bool, error) {
	user, ok := f.userInfo[uname]
	if !ok {
		return false, errors.New("could not add image to repository: user not found")
	}
	newPathname := f.rootDir + "/" + strconv.FormatInt(int64(user.DirId), 10) + "/" + img.NameExt
	imageFile, err := os.Stat(img.OrigPath)
	if err != nil {
		return false, err
	}

	if !imageFile.Mode().IsRegular() {
		return false, fmt.Errorf("%s is not a regular file", img.OrigPath)
	}
	source, err := os.Open(img.OrigPath)
	if err != nil {
		return false, err
	}
	defer func(source *os.File) {
		err := source.Close()
		if err != nil {
		}
	}(source)
	destination, err := os.Create(newPathname)
	if err != nil {
		return false, err
	}
	defer func(destination *os.File) {
		err := destination.Close()
		if err != nil {

		}
	}(destination)
	_, err = io.Copy(destination, source)
	if err != nil {
		return false, err
	}
	img.Path = newPathname
	userdata, ok := f.userMetadata[user.Username]
	if !ok {
		return false, errors.New("user metadata not found")
	}
	userdata.StorageUsed += img.Size
	userdata.AmountOfFiles += 1
	jsonFile, err := os.OpenFile(user.metadataPath, os.O_APPEND|os.O_WRONLY, 0777)
	if err != nil {
		return false, err
	}
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {
		}
	}(jsonFile)
	// Update the metadata of user after adding the image to user folder
	content, err := json.MarshalIndent(userdata, " ", " ")
	if err != nil {
		return false, err
	}
	// Write the new metadata of user to user metadata json
	err = ioutil.WriteFile(user.metadataPath, content, 0644)
	if err != nil {
		return false, err
	}
	user.ImagesInRepo[img.Name] = img
	user.ImgData.AllImages = append(user.ImgData.AllImages, *img)
	jsonFile, err = os.OpenFile(user.ImgPath, os.O_APPEND|os.O_WRONLY, 0777)
	if err != nil {
		return false, err
	}
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {
		}
	}(jsonFile)
	// Update the metadata of user after adding the image to user folder
	content, err = json.MarshalIndent(user.ImgData.AllImages, " ", " ")
	if err != nil {
		return false, err
	}
	// Write the new metadata of user to user metadata json
	err = ioutil.WriteFile(user.ImgPath, content, 0644)
	if err != nil {
		return false, err
	}
	return true, nil
}

//removeUser removes user from the filesystem and all metadata
//@param u: type *User        pointer to the user to be deleted
func (f *Filesystem) removeUser(u *User) (bool, error) {
	_, ok := f.userInfo[u.Username]
	if !ok {
		return false, errors.New("user not found")
	} else {
		delete(f.userInfo, u.Username)
	}
	var index int
	for i, check := range f.usersSoFar.AllUsers {
		if check.Username == u.Username {
			index = i
			break
		}
	}
	f.usersSoFar.AllUsers = removeUserIndex(f.usersSoFar.AllUsers, index)
	jsonFile, err := os.OpenFile(f.userFilePersist, os.O_APPEND|os.O_WRONLY, 0777)
	if err != nil {
		return false, err
	}

	//get  the array of users and update the json file
	content, err := json.MarshalIndent(f.usersSoFar.AllUsers, "", "")
	if err != nil {
		return false, err
	}
	err = ioutil.WriteFile(f.userFilePersist, content, 0644)
	if err != nil {
		return false, err
	}
	err = jsonFile.Close()
	if err != nil {
		return false, err
	}
	err = os.Remove(u.filesys.rootDir + "/" + strconv.FormatInt(int64(u.DirId), 10) + ".json")
	if err != nil {
		return false, errors.New("error deleting user metadata")
	}
	err = os.RemoveAll(u.filesys.rootDir + "/" + strconv.FormatInt(int64(u.DirId), 10))
	if err != nil {
		return false, errors.New("error deleting user directory")
	}
	return true, nil
}
