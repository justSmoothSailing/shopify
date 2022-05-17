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

//Test if the users that were added were added properly with the proper values
func TestInitFilesystem2(t *testing.T) {
	fs, err := InitFilesystem()
	if err != nil {
		t.Fatalf("Failed to Initialize Repository %v", err)
	}
	//These Last checks work if you already have a populated
	//JSON with these two keys
	//If JSON is empty change to len(fs.userInfo) == 0
	juser, ok := fs.userInfo["jhaddad"]
	if ok != true {
		t.Fatalf("user not added to Repository")
	}
	if juser.Username != "jhaddad" {
		t.Fatalf("user name not parsed correctly")
	}
	if juser.FirstName != "james" {
		t.Fatalf("user first name not parsed correctly")
	}
	if juser.LastName != "haddad" {
		t.Fatalf("user last name not parsed correctly")
	}
	if juser.UserStorageUsage != 0 {
		t.Fatalf("user storage not parsed correctly")
	}
	wuser, ok := fs.userInfo["whaddad"]
	if ok != true {
		t.Fatalf("user not added to Repository")
	}
	if wuser.Username != "whaddad" {
		t.Fatalf("user name not parsed correctly")
	}
	if wuser.FirstName != "william" {
		t.Fatalf("user first name not parsed correctly")
	}
	if wuser.LastName != "haddad" {
		t.Fatalf("user last name not parsed correctly")
	}
	if wuser.UserStorageUsage != 0 {
		t.Fatalf("user storage not parsed correctly")
	}

}

// Test Should fail with a username already used
/*
func TestFilesystem_CreateUser(t *testing.T) {
	fs, err := InitFilesystem()
	if err != nil {
		t.Fatalf("Failed to Initialize Repository %v", err)
	}
	user, err := InitCreateUser("jimmyBoy", "Jimmy", "Boy", "765431", fs)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	if user.Username != "jimmyBoy" {
		t.Fatalf("user name not parsed correctly")
	}
	if user.FirstName != "Jimmy" {
		t.Fatalf("user first name not parsed correctly")
	}
	if user.LastName != "Boy" {
		t.Fatalf("user last name not parsed correctly")
	}
	if user.UserStorageUsage != 0 {
		t.Fatalf("user storage was not zero")
	}
	//Cannot check password or DirId since password is encrypted and
	//DirId is a random 32 bit integer
}

*/
