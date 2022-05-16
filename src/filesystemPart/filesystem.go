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
type Metadata struct {
	User          string `json:"user"`
	StorageUsed   int64  `json:"storageUsed"`
	AmountOfFiles int64  `json:"amountOfFiles"`
}

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
			var user User
			for i := 0; i < len(filesystem.usersSoFar.AllUsers); i++ {
				user.UserStorageUsage = filesystem.usersSoFar.AllUsers[i].UserStorageUsage
				user.DirId = filesystem.usersSoFar.AllUsers[i].DirId
				user.Password = filesystem.usersSoFar.AllUsers[i].Password
				user.FirstName = filesystem.usersSoFar.AllUsers[i].FirstName
				user.LastName = filesystem.usersSoFar.AllUsers[i].LastName
				user.Username = filesystem.usersSoFar.AllUsers[i].Username
				filesystem.userInfo[user.Username] = &user
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
	user.metadataPath = f.rootDir + result + user.Username + ".json"
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
	return &user, nil
}
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
