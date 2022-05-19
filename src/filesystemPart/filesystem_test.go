package filesystemPart

import (
	"os"
	"testing"
)

//Test if the initialization of the filesystem works properly
func TestInitFilesystem(t *testing.T) {
	fs, err := InitFilesystem()
	if err != nil {
		t.Fatalf("Failed to Initialize Repository %v", err)
	}
	if fs.userFilePersist != "C:/Users/18645/Documents/temp/users/users.json" {
		t.Fatalf("users info file path is wrong %v", fs.userFileExists)
	}
	if fs.rootDir != "C:/Users/18645/Documents/temp/users" {
		t.Fatalf("root directory set for repository is wrong")
	}
	if _, err := os.Stat(fs.rootDir); os.IsNotExist(err) {
		t.Fatalf("Root Directory was not created")
	}
	if _, err := os.Stat(fs.userFilePersist); os.IsNotExist(err) {
		t.Fatalf("Json file was not created")
	}
}

//Test if the new user was added properly with the proper values
//Should change the values of the InitCreateUser after every test
//or else it fails
func TestInitFilesystem2(t *testing.T) {
	fs, err := InitFilesystem()
	if err != nil {
		t.Fatalf("Failed to Initialize Repository %v", err)
	}
	juser, err := InitCreateUser("doubleTrouble", "bucket", "head", "head123", fs)
	if err != nil {
		t.Fatalf("initializing and creating user failed")
	}
	//These Last checks work if you already have a populated
	//JSON with these two keys
	//If JSON is empty change to len(fs.userInfo) == 0
	juser, ok := fs.userInfo["doubleTrouble"]
	if ok != true {
		t.Fatalf("user not added to Repository")
	}
	if juser.Username != "doubleTrouble" {
		t.Fatalf("user name not parsed correctly")
	}
	if juser.FirstName != "bucket" {
		t.Fatalf("user first name not parsed correctly")
	}
	if juser.LastName != "head" {
		t.Fatalf("user last name not parsed correctly")
	}
	if juser.UserStorageUsage != 0 {
		t.Fatalf("user storage not parsed correctly")
	}

}

// Keep getting this error but I tried sprintf and string(var) and formatInt
//still getting the same error
// This test should find goland.png and delete it from the user
// repository of user: ad
// TODO: Check if this works after build (WORKED!)
func TestUser_AddImageToRepository(t *testing.T) {
	fs, err := InitFilesystem()
	if err != nil {
		t.Fatalf("Failed to Initialize Repository %v", err)
	}
	user, ok := fs.userInfo["wil"]
	if !ok {
		t.Fatalf("User should have existed but does not")
	}
	_, ok = user.ImagesInRepo["goland"]
	if !ok {
		t.Fatal("image not found in users map")
	}

}
