package filesystemPart

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	ro "crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
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
//usersSoFar: type struct Users//struct with an array of all users

type Filesystem struct {
	userInfo        map[string]*User
	rootDir         string
	userMetadata    map[string]Metadata
	userPublicData  map[string][]string
	rootDirExists   bool
	userFileExists  bool
	userFilePersist string
	key             string
	usersSoFar      Users
}

//Struct that holds metadata about the user repository
//Field:
//User: type string          //username of a user in the repository
//StorageUsed: type int64    //length in bytes of a user folder
//AmountOfFiles: type int64  //amount of images in a user folder

type Metadata struct {
	User          string `json:"user"`
	StorageUsed   int64  `json:"storageUsed"`
	AmountOfFiles int64  `json:"amountOfFiles"`
}

//Function that initializes a filesystem and populates the fields that are used
//from metadata files for information used in functions
//the metadata is used to verify users and use parts in other functions
//@param NONE
//return: *FileSystem, nil OR nil, error         //A pointer to a populated FileSystem(Or Repository)

func InitFilesystem() (*Filesystem, error) {
	var errorCaused error = nil
	filesystem := Filesystem{userInfo: make(map[string]*User),
		rootDir: "C:/Users/18645/Documents/temp/users", userMetadata: make(map[string]Metadata),
		userPublicData: make(map[string][]string), rootDirExists: false, userFileExists: false,
		userFilePersist: "C:/Users/18645/Documents/temp/users/users.json", usersSoFar: Users{}}
	filesystem.key = randSeq(32)

	//Check if root directory exists, If not create it
	_, err := os.Stat(filesystem.rootDir)
	if os.IsNotExist(err) {
		errorCaused = os.MkdirAll(filesystem.rootDir, 0777)
		if errorCaused == nil {
			filesystem.rootDirExists = true
		}
	}

	//Check if file of all current users exist, if not create it
	_, err = os.Stat(filesystem.userFilePersist)
	if os.IsNotExist(err) {
		_, errorCaused = os.Create(filesystem.userFilePersist)
		_, err := os.Stat(filesystem.userFilePersist)
		if !os.IsNotExist(err) {
			filesystem.userFileExists = true
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
		byteValue, _ := ioutil.ReadAll(jsonFile)
		err = json.Unmarshal(byteValue, &filesystem.usersSoFar.AllUsers)
		if err == nil {
			for i := 0; i < len(filesystem.usersSoFar.AllUsers); i++ {
				var user *User = new(User)
				user.UserStorageUsage = filesystem.usersSoFar.AllUsers[i].UserStorageUsage
				user.DirId = filesystem.usersSoFar.AllUsers[i].DirId
				user.Password = filesystem.usersSoFar.AllUsers[i].Password
				user.FirstName = filesystem.usersSoFar.AllUsers[i].FirstName
				user.LastName = filesystem.usersSoFar.AllUsers[i].LastName
				user.Username = filesystem.usersSoFar.AllUsers[i].Username
				filesystem.userInfo[filesystem.usersSoFar.AllUsers[i].Username] = user
			}
		} else {
			errorCaused = errors.New("json file could not be parsed")
		}
	}
	if errorCaused != nil {
		return nil, errorCaused
	}
	return &filesystem, errorCaused
}

//Creates User and adds them to the filesystem, creates directory, and initializes metadata
//@param uname              //username of the user
//@param pword              //password supplied by the user
//@param user               //original user of the call
//@return *User, error      //if successful, returns pointer to User and nil, otherwise nil, error
func (f *Filesystem) createUser(uname string, pword string, user User) (*User, error) {
	var keys []byte
	hpword := encrypt(keys, pword)
	newPassword := hex.EncodeToString(hpword)
	user.Password = newPassword
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	//used name of folder for designated user
	user.DirId = r.Uint32()
	dirId := strconv.FormatInt(int64(user.DirId), 10)
	result := "/" + dirId
	user.metadataPath = f.rootDir + result + ".json"
	err := os.MkdirAll(f.rootDir+result, os.ModeDir)
	if err != nil {
		return nil, err
	}
	//Add user to the filesystem map
	f.userInfo[uname] = &user
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
	//Initialize all metadata values to 0 except for the username
	metadata := Metadata{user.Username, 0, 0}
	jsonFile, err := os.OpenFile(user.metadataPath, os.O_APPEND|os.O_WRONLY, 0777)
	if err != nil {
		return nil, err
	}
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {
		}
	}(jsonFile)
	//Append user to the array of users and update the json file
	content, err := json.MarshalIndent(metadata, " ", " ")
	if err != nil {
		return nil, err
	}
	err = ioutil.WriteFile(user.metadataPath, content, 0644)
	if err != nil {
		return nil, err
	}

	return &user, nil
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

//Adds a new user to the filesystem
//@param User               //Newly populated user
//@return err               //nil if user was added properly, err if not
func (f *Filesystem) addUser(user User) error {
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
	f.usersSoFar.AllUsers = append(f.usersSoFar.AllUsers, user)
	content, err := json.MarshalIndent(f.usersSoFar.AllUsers, " ", " ")
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
func (f *Filesystem) addImg(img Image, uname string) (bool, error) {
	user, ok := f.userInfo[uname]
	if !ok {
		return false, errors.New("could not add image to repository: user not found")
	}
	newPathname := f.rootDir + strconv.FormatInt(int64(user.DirId), 10) + img.nameExt
	err := os.Rename(img.origPath, newPathname)
	if err != nil {
		return false, err
	}
	_, err = os.Stat(newPathname)
	if err != nil {
		return false, errors.New("failed to move image to repository")
	}
	img.path = newPathname
	userdata, ok := f.userMetadata[user.Username]
	if !ok {
		return false, errors.New("user metadata not found")
	}
	userdata.StorageUsed += img.size
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
	//Append user to the array of users and update the json file
	content, err := json.MarshalIndent(userdata, " ", " ")
	if err != nil {
		return false, err
	}
	err = ioutil.WriteFile(f.userFilePersist, content, 0644)
	if err != nil {
		return false, err
	}
	return true, nil
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}
func encrypt(data []byte, passphrase string) []byte {
	block, _ := aes.NewCipher([]byte(createHash(passphrase)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(ro.Reader, nonce); err != nil {
		panic(err.Error())
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext
}
